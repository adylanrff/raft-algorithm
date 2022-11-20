package rpc

import (
	"net"
)

type RPCClient interface {
	Call(address string, payload *ServerMessageDTO) (resp *ServerMessageDTO, err error)
}

type defaultRpcClient struct {
}

// Call implements RPCClient
func (*defaultRpcClient) Call(address string, payload *ServerMessageDTO) (resp *ServerMessageDTO, err error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	byteSrvMsg, err := payload.ToBytes()
	if err != nil {
		return nil, err
	}

	_, err = conn.Write(byteSrvMsg)
	if err != nil {
		return nil, err
	}

	resp, err = ParseServerMessage(conn)
	if err != nil {
		return nil, err
	}

	return
}

func NewRPCClient() RPCClient {
	return &defaultRpcClient{}
}
