package model

import raftPb "github.com/adylanrff/raft-algorithm/proto/raft"

type RequestVoteRequestDTO struct {
	*raftPb.RequestVoteRequest
}

type RequestVoteResponseDTO struct {
	*raftPb.RequestVoteResponse
}

func NewRequestVoteResponseDTO(term int32, voteGranted bool) *RequestVoteResponseDTO {
	return &RequestVoteResponseDTO{
		RequestVoteResponse: &raftPb.RequestVoteResponse{
			Term:        term,
			VoteGranted: voteGranted,
		},
	}
}
