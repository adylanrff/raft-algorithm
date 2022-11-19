package raft

import (
	"errors"

	raft "github.com/adylanrff/raft-algorithm/proto"
	"github.com/adylanrff/raft-algorithm/raft/model"
	log "github.com/sirupsen/logrus"
)

type RaftServerHandler interface {
	RequestVote(req *model.RequestVoteRequestDTO) (*model.RequestVoteResponseDTO, error)
	AppendEntries(req *model.AppendEntriesRequestDTO) (*model.AppendEntriesResponseDTO, error)
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

// AppendEntries implements RaftServerHandler
func (*defaultRaftServerHandler) AppendEntries(req *model.AppendEntriesRequestDTO) (*model.AppendEntriesResponseDTO, error) {
	log.WithFields(log.Fields{
		"req":    req.AppendEntriesRequest,
		"method": "AppendEntries",
	}).Debug("append_entries_request_start")

	return &model.AppendEntriesResponseDTO{}, nil
}

// RequestVote implements RaftServerHandler
func (*defaultRaftServerHandler) RequestVote(req *model.RequestVoteRequestDTO) (*model.RequestVoteResponseDTO, error) {
	log.WithFields(log.Fields{
		"req":    req.RequestVoteRequest,
		"method": "RequestVote",
	}).Debug("request_vote_request_start")

	return &model.RequestVoteResponseDTO{}, nil
}

func (h *defaultRaftServerHandler) handleAppendEntries(req *model.RaftMessageDTO) (resp *model.RaftMessageDTO, err error) {
	raftReq := req.GetRaftRequest()
	appendEntriesReq := &model.AppendEntriesRequestDTO{
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
	requestVoteReq := &model.RequestVoteRequestDTO{
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
