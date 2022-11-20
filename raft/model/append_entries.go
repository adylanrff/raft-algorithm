package model

import raftPb "github.com/adylanrff/raft-algorithm/proto/raft"

type AppendEntriesRequestDTO struct {
	*raftPb.AppendEntriesRequest
}

func NewAppendEntriesRequestDTO() *AppendEntriesRequestDTO {
	return &AppendEntriesRequestDTO{
		AppendEntriesRequest: &raftPb.AppendEntriesRequest{},
	}
}

type AppendEntriesResponseDTO struct {
	*raftPb.AppendEntriesResponse
}
