package server

import (
	serverPb "github.com/adylanrff/raft-algorithm/proto/server"
	"google.golang.org/protobuf/proto"
)

type ServerMessageDTO struct {
	*serverPb.ServerMessage
}

func (dto *ServerMessageDTO) ToBytes() ([]byte, error) {
	marshalledMessage, err := proto.Marshal(dto.ServerMessage)
	if err != nil {
		return nil, err
	}

	return append([]byte{byte(len(marshalledMessage))}, marshalledMessage...), nil
}
