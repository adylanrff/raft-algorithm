package raft

import (
	"time"

	"github.com/adylanrff/raft-algorithm/raft/model"
	log "github.com/sirupsen/logrus"
)

type Raft interface {
	RequestVote(req *model.RequestVoteRequestDTO) (*model.RequestVoteResponseDTO, error)
	AppendEntries(req *model.AppendEntriesRequestDTO) (*model.AppendEntriesResponseDTO, error)
}

type defaultRaft struct {
	state  model.RaftState
	config RaftConfig

	// timer
	electionTimer time.Timer

	appendEntriesSignal chan struct{}
}

func NewRaft(id string, cfg RaftConfig) *defaultRaft {
	return &defaultRaft{
		state:         model.NewRaftState(id),
		config:        cfg,
		electionTimer: *time.NewTimer(cfg.ElectionTimeout),
	}
}

// RequestVote implements RaftServerHandler
func (r *defaultRaft) RequestVote(req *model.RequestVoteRequestDTO) (*model.RequestVoteResponseDTO, error) {
	log.WithFields(log.Fields{
		"req":    req.RequestVoteRequest,
		"method": "RequestVote",
	}).Debug("request_vote_request_start")

	r.state.Lock()
	defer r.state.Unlock()

	candidateTerm := req.GetTerm()

	// reply false if the term is less than our current term
	if candidateTerm < r.state.CurrentTerm {
		return model.NewRequestVoteResponseDTO(r.state.CurrentTerm, false), nil
	}

	// Reply false also if we already voted and it is not the candidate ID
	if r.state.VotedFor != req.GetCandidateId() && r.state.VotedFor != "" {
		return model.NewRequestVoteResponseDTO(r.state.CurrentTerm, false), nil
	}

	// Grant vote!
	r.state.Vote(req.GetCandidateId())

	return model.NewRequestVoteResponseDTO(r.state.CurrentTerm, true), nil
}

// AppendEntries implements RaftServerHandler
func (r *defaultRaft) AppendEntries(req *model.AppendEntriesRequestDTO) (*model.AppendEntriesResponseDTO, error) {
	log.WithFields(log.Fields{
		"req":    req.AppendEntriesRequest,
		"method": "AppendEntries",
	}).Debug("append_entries_request_start")

	// received a heartbeat
	r.appendEntriesSignal <- struct{}{}
	return &model.AppendEntriesResponseDTO{}, nil
}

func (r *defaultRaft) Run() error {
	for {
		switch r.state.Role {
		case model.RaftRoleFollower:
			r.doFollowerAction()
		case model.RaftRoleCandidate:
			r.doCandidateAction()
		case model.RaftRoleLeader:
			r.doLeaderAction()
		}
	}
}

func (r *defaultRaft) doFollowerAction() {
	for {
		select {
		case <-r.electionTimer.C:
			// be a candidate
			// TODO: fix race condition possibilities,
			// should not happen because the only thread that modifies the state.role is this thread
			r.state.ChangeRole(model.RaftRoleCandidate)
			return
		case <-r.appendEntriesSignal:
			log.Debug("heartbeat received")
			// reset election timeout
			r.resetElectionTimer()
			continue
		}
	}
}

func (r *defaultRaft) doCandidateAction() {
	electionChan := make(chan struct{})

	for {
		select {
		case <-r.electionTimer.C:
			// reset election timer
			r.resetElectionTimer()
			// do another election
			electionChan <- struct{}{}
		case <-r.appendEntriesSignal:
			r.state.ChangeRole(model.RaftRoleFollower)
			return
		case <-electionChan:
			shouldReturn := r.doElection()
			if shouldReturn {
				return
			}
		}
	}
}

func (r *defaultRaft) doElection() bool {
	r.state.Lock()
	defer r.state.Lock()

	// 1. increment current term
	r.state.CurrentTerm++
	// 2. vote for self
	r.state.Vote(r.state.ID)
	// 3. reset election timer
	r.resetElectionTimer()
	// calculate majority:
	majority := len(r.config.ClusterMemberAddreses)/2 + 1
	// 1 because we voted for ourselves
	votes := 1

	// TODO: call requestvote for each addresses
	for _, address := range r.config.ClusterMemberAddreses {
		requestVoteResult := false
		if address != "" {
			requestVoteResult = true
		}
		if requestVoteResult {
			votes++
		}
	}

	if votes >= majority {
		r.state.ChangeRole(model.RaftRoleLeader)
		return true
	}

	return false
}

func (r *defaultRaft) doLeaderAction() {
	ticker := time.NewTicker(r.config.HeartbeatInterval)

	for range ticker.C {
		// for each addresses, call empty appendEntries
		continue
	}
}

func (r *defaultRaft) resetElectionTimer() {
	if !r.electionTimer.Stop() {
		<-r.electionTimer.C
	}
	r.electionTimer.Reset(r.config.ElectionTimeout)
}
