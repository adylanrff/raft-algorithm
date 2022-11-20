package model

import raftPb "github.com/adylanrff/raft-algorithm/proto/raft"

type RequestVoteRequestDTO struct {
	*raftPb.RequestVoteRequest
}

type RequestVoteResponseDTO struct {
	*raftPb.RequestVoteResponse
}
