package raft

import (
	"errors"
	"io"

	raft "github.com/adylanrff/raft-algorithm/proto"
	"github.com/adylanrff/raft-algorithm/raft/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func ParseRaftMessage(reader io.Reader) (*model.RaftMessageDTO, error) {
	byteSize, err := readByteSize(reader)
	if err != nil {
		return nil, err
	}

	return readRaftMessage(reader, byteSize)
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

func readRaftMessage(reader io.Reader, byteSize int) (*model.RaftMessageDTO, error) {
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

	var raftMessage raft.RaftMessage

	err = proto.Unmarshal(messageBuf[:n], &raftMessage)
	if err != nil {
		return nil, err
	}

	return &model.RaftMessageDTO{
		RaftMessage: &raftMessage,
	}, nil
}
