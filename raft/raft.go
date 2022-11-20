package raft

import (
	"github.com/adylanrff/raft-algorithm/raft/model"
	log "github.com/sirupsen/logrus"
)

type Raft interface {
	RequestVote(req *model.RequestVoteRequestDTO) (*model.RequestVoteResponseDTO, error)
	AppendEntries(req *model.AppendEntriesRequestDTO) (*model.AppendEntriesResponseDTO, error)
}

type defaultRaftHandler struct {
}

func NewRaftHandler() *defaultRaftHandler {
	return &defaultRaftHandler{}
}

// AppendEntries implements RaftServerHandler
func (*defaultRaftHandler) AppendEntries(req *model.AppendEntriesRequestDTO) (*model.AppendEntriesResponseDTO, error) {
	log.WithFields(log.Fields{
		"req":    req.AppendEntriesRequest,
		"method": "AppendEntries",
	}).Debug("append_entries_request_start")

	return &model.AppendEntriesResponseDTO{}, nil
}

// RequestVote implements RaftServerHandler
func (*defaultRaftHandler) RequestVote(req *model.RequestVoteRequestDTO) (*model.RequestVoteResponseDTO, error) {
	log.WithFields(log.Fields{
		"req":    req.RequestVoteRequest,
		"method": "RequestVote",
	}).Debug("request_vote_request_start")

	return &model.RequestVoteResponseDTO{}, nil
}
