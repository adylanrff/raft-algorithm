package model

import (
	pb "github.com/adylanrff/raft-algorithm/proto"
	"google.golang.org/protobuf/proto"
)

type RaftMessageDTO struct {
	*pb.RaftMessage
}

func (dto *RaftMessageDTO) ToBytes() ([]byte, error) {
	marshalledMessage, err := proto.Marshal(dto.RaftMessage)
	if err != nil {
		return nil, err
	}

	return append([]byte{byte(len(marshalledMessage))}, marshalledMessage...), nil
}
