package raft

import (
	raftModel "github.com/adylanrff/raft-algorithm/raft/model"

	serverPb "github.com/adylanrff/raft-algorithm/proto/server"
	"github.com/adylanrff/raft-algorithm/server"
)

// TODO: This sucks... do a better way to handle various request/response
// especially the response generation
// perhaps codegen would be a good way (just like what grpc do)

type RaftServerHandler struct {
	raftHandler Raft
}

func NewRaftServerhandler(raftHandler Raft) *RaftServerHandler {
	return &RaftServerHandler{
		raftHandler: raftHandler,
	}
}

func (h *RaftServerHandler) AppendEntriesHandler(req *server.ServerMessageDTO) (resp *server.ServerMessageDTO, err error) {
	appendEntriesDTO := &raftModel.AppendEntriesRequestDTO{
		AppendEntriesRequest: req.GetServerRequest().GetAppendEntriesRequest(),
	}

	appendEntriesResp, appendEntriesErr := h.raftHandler.AppendEntries(appendEntriesDTO)

	resp = &server.ServerMessageDTO{
		ServerMessage: &serverPb.ServerMessage{
			Method: req.GetMethod(),
			Payload: &serverPb.ServerMessage_ServerResponse{
				ServerResponse: &serverPb.ServerResponse{
					Response: &serverPb.ServerResponse_AppendEntriesResponse{
						AppendEntriesResponse: appendEntriesResp.AppendEntriesResponse,
					},
				},
			},
		},
	}

	err = appendEntriesErr

	return
}

func (h *RaftServerHandler) RequestVoteHandler(req *server.ServerMessageDTO) (resp *server.ServerMessageDTO, err error) {
	requestVoteDTO := &raftModel.RequestVoteRequestDTO{
		RequestVoteRequest: req.GetServerRequest().GetRequestVoteRequest(),
	}

	requestVoteResp, requestVoteErr := h.raftHandler.RequestVote(requestVoteDTO)
	resp = &server.ServerMessageDTO{
		ServerMessage: &serverPb.ServerMessage{
			Method: req.GetMethod(),
			Payload: &serverPb.ServerMessage_ServerResponse{
				ServerResponse: &serverPb.ServerResponse{
					Response: &serverPb.ServerResponse_RequestVoteResponse{
						RequestVoteResponse: requestVoteResp.RequestVoteResponse,
					},
				},
			},
		},
	}
	err = requestVoteErr

	return
}
