package model

import raft "github.com/adylanrff/raft-algorithm/proto"

type RequestVoteRequest struct {
	*raft.RequestVoteRequest
}

type RequestVoteResponse struct {
	*raft.RequestVoteResponse
}
