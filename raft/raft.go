package raft

import (
	"errors"

	raft "github.com/adylanrff/raft-algorithm/proto"
	"github.com/adylanrff/raft-algorithm/raft/model"
)

type Entry struct {
}

type RaftServerHandler interface {
	RequestVote(req *model.RequestVoteRequest) (*model.RequestVoteResponse, error)
	AppendEntries(req *model.AppendEntriesRequest) (*model.AppendEntriesResponse, error)
}

type defaultRaftServerHandler struct {
}

func NewRaftServerHandler() *defaultRaftServerHandler {
	return &defaultRaftServerHandler{}
}

func (h *defaultRaftServerHandler) Handle(req *model.RaftMessageDTO) (resp *model.RaftMessageDTO, err error) {
	switch req.GetMethod() {
	case RaftMethodName_AppendEntries:
		return h.handleAppendEntries(req)
	case RaftMethodName_RequestVotes:
		return h.handleRequestVote(req)
	}

	return nil, errors.New("unrecognized method")
}

func (h *defaultRaftServerHandler) handleAppendEntries(req *model.RaftMessageDTO) (resp *model.RaftMessageDTO, err error) {
	raftReq := req.GetRaftRequest()
	appendEntriesReq := &model.AppendEntriesRequest{
		AppendEntriesRequest: raftReq.GetAppendEntriesRequest(),
	}

	appendEntriesResp, appendEntriesErr := h.AppendEntries(appendEntriesReq)
	resp = &model.RaftMessageDTO{
		RaftMessage: &raft.RaftMessage{
			Method:          req.GetMethod(),
			RaftMessageType: raft.RaftMessageType_RaftMessageType_RaftResponse,
			Payload: &raft.RaftMessage_RaftResponse{
				RaftResponse: &raft.RaftResponse{
					Response: &raft.RaftResponse_AppendEntriesResponse{
						AppendEntriesResponse: appendEntriesResp.AppendEntriesResponse,
					},
				},
			},
		},
	}
	err = appendEntriesErr

	return
}

func (h *defaultRaftServerHandler) handleRequestVote(req *model.RaftMessageDTO) (resp *model.RaftMessageDTO, err error) {
	raftReq := req.GetRaftRequest()
	requestVoteReq := &model.RequestVoteRequest{
		RequestVoteRequest: raftReq.GetRequestVoteRequest(),
	}

	requestVoteResp, requestVoteErr := h.RequestVote(requestVoteReq)
	resp = &model.RaftMessageDTO{
		RaftMessage: &raft.RaftMessage{
			Method:          req.GetMethod(),
			RaftMessageType: raft.RaftMessageType_RaftMessageType_RaftResponse,
			Payload: &raft.RaftMessage_RaftResponse{
				RaftResponse: &raft.RaftResponse{
					Response: &raft.RaftResponse_RequestVoteResponse{
						RequestVoteResponse: requestVoteResp.RequestVoteResponse,
					},
				},
			},
		},
	}
	err = requestVoteErr

	return
}

// AppendEntries implements RaftServerHandler
func (*defaultRaftServerHandler) AppendEntries(req *model.AppendEntriesRequest) (*model.AppendEntriesResponse, error) {
	panic("unimplemented")
}

// RequestVote implements RaftServerHandler
func (*defaultRaftServerHandler) RequestVote(req *model.RequestVoteRequest) (*model.RequestVoteResponse, error) {
	panic("unimplemented")
}
