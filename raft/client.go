package raft

type RaftClient interface {
	// TODO: add entries implementation
	AppendEntries(address string, entries []interface{}) error
	RequestVote(address string) error
}

type defaultRaftClient struct {
}

// AppendEntries implements RaftClient
func (*defaultRaftClient) AppendEntries(address string, entries []interface{}) error {
	panic("unimplemented")
}

// RequestVote implements RaftClient
func (*defaultRaftClient) RequestVote(address string) error {
	panic("unimplemented")
}

func NewRaftClient() RaftClient {
	return &defaultRaftClient{}
}
