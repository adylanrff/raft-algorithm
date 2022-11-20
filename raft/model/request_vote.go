package model

import raftPb "github.com/adylanrff/raft-algorithm/proto/raft"

type RequestVoteRequestDTO struct {
	*raftPb.RequestVoteRequest
}

func NewRequestVoteRequestDTO(
	term int32,
	candidateID string,
	lastLogIndex int32,
	lastLogTerm int32) *RequestVoteRequestDTO {
	return &RequestVoteRequestDTO{
		RequestVoteRequest: &raftPb.RequestVoteRequest{
			Term:         term,
			CandidateId:  candidateID,
			LastLogIndex: lastLogIndex,
			LastLogTerm:  int32(lastLogTerm),
		},
	}
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
