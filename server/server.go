package server

import (
	"errors"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
)

type Handler func(req *ServerMessageDTO) (resp *ServerMessageDTO, err error)

// Server - implements the raft server
type Server struct {
	address  string
	Handlers map[string]Handler

	Parse func(reader io.Reader) (req *ServerMessageDTO, err error)
}

func NewServer(address string) *Server {
	return &Server{
		address:  address,
		Parse:    ParseServerMessage,
		Handlers: make(map[string]Handler),
	}
}

func (s *Server) Run() {
	srv, err := net.Listen("tcp", s.address)
	if err != nil {
		panic(err)
	}
	defer srv.Close()

	log.WithFields(log.Fields{
		"address": s.address,
	}).Info("running server")

	for {
		conn, err := srv.Accept()
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("connection eror")
			continue
		}

		// TODO: add workercount limiter
		go s.handleConn(conn)
	}
}

func (s *Server) AddHandler(method string, handler Handler) {
	s.Handlers[method] = handler
}

func (s *Server) handleConn(conn net.Conn) error {
	var (
		err error
		msg *ServerMessageDTO
	)

	defer func() {
		// TODO: use better way, create ErrorDTO
		if err != nil {
			errorServerMessageDTO := NewErrorMessageDTO(-1, err.Error())
			errorServerMessageBytes, convErr := errorServerMessageDTO.ToBytes()
			if convErr != nil {
				log.WithFields(
					log.Fields{
						"errorDTO": errorServerMessageDTO,
						"err":      convErr,
					}).Error("convert error")

				errBytes := []byte(convErr.Error())
				msg := append([]byte{byte(len(errBytes))}, errBytes...)
				_, err = conn.Write(msg)
				if err != nil {
					log.WithFields(
						log.Fields{
							"errorDTO": errorServerMessageDTO,
							"err":      convErr,
						}).Error("write error")
				}
				return
			}
			conn.Write(errorServerMessageBytes)
		}
	}()

	msg, err = s.Parse(conn)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("parse message error")
		return err
	}

	resp, err := s.handleMsg(msg)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("handle error")
		return err
	}

	respByte, err := resp.ToBytes()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("handle error")
		return err
	}

	var (
		n int
	)
	n, err = conn.Write(respByte)
	log.WithFields(log.Fields{
		"err":           err,
		"bytes_written": n,
		"resp":          resp,
	}).Debug("write response")

	return err
}

func (s *Server) handleMsg(msg *ServerMessageDTO) (resp *ServerMessageDTO, err error) {
	method := msg.GetMethod()
	handler, ok := s.Handlers[method]
	if !ok {
		return nil, errors.New("unrecognized methods")
	}

	return handler(msg)
}
