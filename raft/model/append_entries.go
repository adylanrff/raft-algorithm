package model

import raft "github.com/adylanrff/raft-algorithm/proto"

type AppendEntriesRequest struct {
	*raft.AppendEntriesRequest
}

type AppendEntriesResponse struct {
	*raft.AppendEntriesResponse
}
