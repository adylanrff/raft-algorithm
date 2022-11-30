package raft

import (
	serverPb "github.com/adylanrff/raft-algorithm/proto/server"
	"github.com/adylanrff/raft-algorithm/raft/model"
	"github.com/adylanrff/raft-algorithm/rpc"
)

type RaftClient interface {
	// TODO: add entries implementation
	AppendEntries(address string, req *model.AppendEntriesRequestDTO) (resp *model.AppendEntriesResponseDTO, err error)
	RequestVote(address string, req *model.RequestVoteRequestDTO) (resp *model.RequestVoteResponseDTO, err error)
}

type defaultRaftClient struct {
	rpcClient rpc.RPCClient
}

// AppendEntries implements RaftClient
// TODO: implement entries
func (c *defaultRaftClient) AppendEntries(address string, req *model.AppendEntriesRequestDTO) (resp *model.AppendEntriesResponseDTO, err error) {
	msg := &rpc.ServerMessageDTO{
		ServerMessage: &serverPb.ServerMessage{
			Method: RaftMethodName_AppendEntries,
			Payload: &serverPb.ServerMessage_ServerRequest{
				ServerRequest: &serverPb.ServerRequest{
					Request: &serverPb.ServerRequest_AppendEntriesRequest{
						AppendEntriesRequest: req.AppendEntriesRequest,
					},
				},
			},
		},
	}

	rsp, err := c.rpcClient.Call(address, msg)
	if err != nil {
		return nil, err
	}

	appendEntriesResp := rsp.GetServerResponse().GetAppendEntriesResponse()
	return &model.AppendEntriesResponseDTO{
		AppendEntriesResponse: appendEntriesResp,
	}, nil
}

// RequestVote implements RaftClient
func (c *defaultRaftClient) RequestVote(address string, req *model.RequestVoteRequestDTO) (resp *model.RequestVoteResponseDTO, err error) {
	msg := &rpc.ServerMessageDTO{
		ServerMessage: &serverPb.ServerMessage{
			Method: RaftMethodName_RequestVotes,
			Payload: &serverPb.ServerMessage_ServerRequest{
				ServerRequest: &serverPb.ServerRequest{
					Request: &serverPb.ServerRequest_RequestVoteRequest{
						RequestVoteRequest: req.RequestVoteRequest,
					},
				},
			},
		},
	}

	rsp, err := c.rpcClient.Call(address, msg)
	if err != nil {
		return nil, err
	}

	requestVoteRsp := rsp.GetServerResponse().GetRequestVoteResponse()
	return &model.RequestVoteResponseDTO{
		RequestVoteResponse: requestVoteRsp,
	}, nil
}

func NewRaftClient(rpcClient rpc.RPCClient) RaftClient {
	return &defaultRaftClient{
		rpcClient: rpcClient,
	}
}
