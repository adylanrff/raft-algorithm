package model

import raft "github.com/adylanrff/raft-algorithm/proto"

type RequestVoteRequestDTO struct {
	*raft.RequestVoteRequest
}

type RequestVoteResponseDTO struct {
	*raft.RequestVoteResponse
}
