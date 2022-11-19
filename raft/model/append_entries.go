package model

import raft "github.com/adylanrff/raft-algorithm/proto"

type AppendEntriesRequestDTO struct {
	*raft.AppendEntriesRequest
}

type AppendEntriesResponseDTO struct {
	*raft.AppendEntriesResponse
}
