package rpc

import (
	"errors"
	"io"

	serverPb "github.com/adylanrff/raft-algorithm/proto/server"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func ParseServerMessage(reader io.Reader) (*ServerMessageDTO, error) {
	byteSize, err := readByteSize(reader)
	if err != nil {
		return nil, err
	}

	return readServerMessage(reader, byteSize)
}

func readByteSize(reader io.Reader) (int, error) {
	byteSizeBuf := make([]byte, 1)
	n, err := reader.Read(byteSizeBuf)
	if err != nil {
		if err != io.EOF {
			log.WithFields(log.Fields{
				"n": n,
			}).Error("read error")
		}

		return 0, err
	}

	byteSize := int(byteSizeBuf[0])
	return byteSize, nil
}

func readServerMessage(reader io.Reader, byteSize int) (*ServerMessageDTO, error) {
	messageBuf := make([]byte, byteSize)
	n, err := reader.Read(messageBuf)
	if err != nil {
		if err != io.EOF {
			log.WithFields(log.Fields{
				"n": n,
			}).Error("read error")
		}

		return nil, err
	}

	if n != byteSize {
		log.WithFields(log.Fields{
			"n":                n,
			"header_byte_size": byteSize,
		}).Error("byte size differs from header size")
		return nil, errors.New("unexpected byte size")
	}

	var serverMessage serverPb.ServerMessage

	err = proto.Unmarshal(messageBuf[:n], &serverMessage)
	if err != nil {
		return nil, err
	}

	return &ServerMessageDTO{
		ServerMessage: &serverMessage,
	}, nil
}
