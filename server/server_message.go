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

type ErrorServerMessageDTO struct {
	*ServerMessageDTO
}

func NewErrorMessageDTO(errorCode int32, errorMsg string) *ErrorServerMessageDTO {
	errorResponse := &serverPb.ErrorResponse{
		ErrorCode: errorCode,
		ErrorMsg:  errorMsg,
	}

	return &ErrorServerMessageDTO{
		ServerMessageDTO: &ServerMessageDTO{
			ServerMessage: &serverPb.ServerMessage{
				Payload: &serverPb.ServerMessage_ServerResponse{
					ServerResponse: &serverPb.ServerResponse{
						Response: &serverPb.ServerResponse_ErrorResponse{
							ErrorResponse: errorResponse,
						},
					},
				},
			},
		},
	}
}
