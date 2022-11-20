package model

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

// RaftState
type RaftState struct {
	sync.Mutex

	// persistent state
	ID          string
	CurrentTerm int32
	VotedFor    string // CandidateID for the current term election
	Log         []RaftLog
	Role        RaftRole

	// volatile state
	CommitIndex int32
	LastApplied int32

	// leaders only
	NextIndex  map[string]int
	MatchIndex map[string]int
}

func (s *RaftState) Vote(candidateID string) {
	log.WithFields(log.Fields{
		"id":       s.ID,
		"vote_for": candidateID,
	}).Info("voting done")

	s.VotedFor = candidateID
}

func (s *RaftState) ChangeRole(raftRole RaftRole) {
	log.WithFields(log.Fields{
		"role_before": s.Role.String(),
		"role_after":  raftRole.String(),
		"id":          s.ID,
	}).Info("raft role changed")

	s.Role = raftRole
}

func NewRaftState(id string) RaftState {
	return RaftState{
		ID:          id,
		CurrentTerm: 0,
		VotedFor:    "",
		Log:         make([]RaftLog, 0),
		Role:        RaftRoleFollower,

		CommitIndex: 0,
		LastApplied: 0,

		// for leaders only, set to nil first
		NextIndex:  nil,
		MatchIndex: nil,
	}
}

// RaftLog
type RaftLog struct{}

// RaftRole
type RaftRole int32

const (
	RaftRoleFollower RaftRole = iota
	RaftRoleCandidate
	RaftRoleLeader
)

func (r RaftRole) String() string {
	switch r {
	case RaftRoleCandidate:
		return "candidate"
	case RaftRoleFollower:
		return "follower"
	case RaftRoleLeader:
		return "leader"
	}

	return ""
}
