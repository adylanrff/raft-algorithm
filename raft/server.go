package raft

import (
	"fmt"
	"io"
	"net"

	"github.com/adylanrff/raft-algorithm/raft/model"
	log "github.com/sirupsen/logrus"
)

// Server - implements the raft server
type Server struct {
	port int

	Parse   func(reader io.Reader) (req *model.RaftMessageDTO, err error)
	Handler func(req *model.RaftMessageDTO) (resp *model.RaftMessageDTO, err error)
}

func NewServer(port int) *Server {
	return &Server{
		port:    port,
		Parse:   ParseRaftMessage,
		Handler: NewRaftServerHandler().Handle,
	}
}

func (s *Server) Run() {
	srv, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", s.port))
	if err != nil {
		panic(err)
	}
	defer srv.Close()

	log.WithFields(log.Fields{
		"port": s.port,
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

func (s *Server) handleConn(conn net.Conn) error {
	var (
		err error
		msg *model.RaftMessageDTO
	)

	msg, err = s.Parse(conn)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("parse message error")
		return err
	}

	resp, err := s.Handler(msg)
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
