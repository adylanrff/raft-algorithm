package raft

import (
	"math/rand"
	"time"

	"github.com/adylanrff/raft-algorithm/raft/model"
	"github.com/adylanrff/raft-algorithm/rpc"
	log "github.com/sirupsen/logrus"
)

type Raft interface {
	RequestVote(req *model.RequestVoteRequestDTO) (*model.RequestVoteResponseDTO, error)
	AppendEntries(req *model.AppendEntriesRequestDTO) (*model.AppendEntriesResponseDTO, error)
}

type defaultRaft struct {
	address string
	state   model.RaftState
	config  RaftConfig

	raftClient RaftClient

	receiveAppendEntriesSignal chan struct{}
	sendAppendEntriesSignal    chan struct{}
}

func NewRaft(id string, address string, cfg RaftConfig) *defaultRaft {
	electionTimeout := rand.Intn(int(cfg.ElectionTimeout*2)-int(cfg.ElectionTimeout)) + int(cfg.ElectionTimeout)
	cfg.ElectionTimeout = time.Duration(electionTimeout)

	return &defaultRaft{
		address:    address,
		state:      model.NewRaftState(id),
		config:     cfg,
		raftClient: NewRaftClient(rpc.NewRPCClient()),

		receiveAppendEntriesSignal: make(chan struct{}),
		sendAppendEntriesSignal:    make(chan struct{}),
	}
}

// RequestVote implements RaftServerHandler
func (r *defaultRaft) RequestVote(req *model.RequestVoteRequestDTO) (*model.RequestVoteResponseDTO, error) {
	r.state.Lock()
	defer r.state.Unlock()

	log.WithFields(log.Fields{
		"id":       r.state.ID,
		"votedFor": r.state.VotedFor,
		"req":      req.RequestVoteRequest,
		"method":   "RequestVote",
	}).Debug("request_vote_request_start")

	// reply false if the term is less than our current term
	if req.GetTerm() < r.state.CurrentTerm {
		log.WithFields(log.Fields{
			"candidate_term": req.GetTerm(),
			"current_term":   r.state.CurrentTerm,
		}).Error("outdated term")
		return model.NewRequestVoteResponseDTO(r.state.CurrentTerm, false), nil
	}

	// Reply false also if we already voted and it is not the candidate ID
	if r.state.VotedFor != req.GetCandidateId() && r.state.VotedFor != "" {
		log.WithFields(log.Fields{
			"voted_for":    r.state.VotedFor,
			"candidate_id": req.GetCandidateId(),
		}).Error("already voted before")
		return model.NewRequestVoteResponseDTO(r.state.CurrentTerm, false), nil
	}

	// Grant vote!
	r.state.VotedFor = req.GetCandidateId()
	resp := model.NewRequestVoteResponseDTO(r.state.CurrentTerm, true)
	log.WithFields(log.Fields{
		"id":       r.state.ID,
		"votedFor": r.state.VotedFor,
		"resp":     resp,
		"method":   "RequestVote",
	}).Debug("request_vote_request_end")

	return resp, nil
}

// AppendEntries implements RaftServerHandler
func (r *defaultRaft) AppendEntries(req *model.AppendEntriesRequestDTO) (*model.AppendEntriesResponseDTO, error) {
	log.WithFields(log.Fields{
		"req":    req.AppendEntriesRequest,
		"method": "AppendEntries",
	}).Debug("append_entries_request_start")

	// received a heartbeat
	r.receiveAppendEntriesSignal <- struct{}{}
	return &model.AppendEntriesResponseDTO{}, nil
}

func (r *defaultRaft) Run() error {
	go r.watcher()

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
	log.WithFields(log.Fields{}).Infof("do follower action start")

	for {
		select {
		case <-time.After(r.config.ElectionTimeout):
			// be a candidate
			// TODO: fix race condition possibilities,
			// should not happen because the only thread that modifies the state.role is this thread
			r.state.ChangeRole(model.RaftRoleCandidate)
			return
		case <-r.receiveAppendEntriesSignal:
			log.Debug("append entries received")
		}
	}
}

func (r *defaultRaft) doCandidateAction() {
	log.WithFields(log.Fields{}).Infof("do candidate action startzzz")

	votes := 1
	majority := len(r.config.ClusterMemberAddreses) / 2

	log.WithFields(log.Fields{}).Infof("do candidate action send election")
	voteRespChan := r.doElection()

	for {
		select {
		case resp, ok := <-voteRespChan:
			if ok && resp.VoteGranted {
				votes++
			}
			if votes > majority {
				r.state.ChangeRole(model.RaftRoleLeader)
				return
			}
		case <-time.After(r.config.ElectionTimeout):
			log.WithFields(log.Fields{}).Infof("election timeout")
			// do another election
			return
		case <-r.receiveAppendEntriesSignal:
			log.WithFields(log.Fields{}).Infof("heartbeat received")
			r.state.Lock()
			r.state.ChangeRole(model.RaftRoleFollower)
			r.state.VotedFor = ""
			r.state.Unlock()
			return
		}
	}
}

func (r *defaultRaft) doElection() chan *model.RequestVoteResponseDTO {
	log.WithFields(log.Fields{}).Infof("do election start")

	r.state.Lock()
	defer r.state.Unlock()

	// 1. increment current term
	r.state.CurrentTerm++
	// 2. vote for self
	r.state.VotedFor = r.state.ID

	// Do voting
	log.WithFields(log.Fields{}).Infof("do voting")
	requestVoteReq := model.NewRequestVoteRequestDTO(
		r.state.CurrentTerm,
		r.state.ID,
		int32(len(r.state.Logs)),
		0)

	voteRespChan := make(chan *model.RequestVoteResponseDTO, len(r.config.ClusterMemberAddreses))

	for _, memberAddress := range r.config.ClusterMemberAddreses {
		if memberAddress != r.address {
			go func(memberAddress string) {
				resp, err := r.raftClient.RequestVote(memberAddress, requestVoteReq)
				if err != nil {
					log.WithFields(log.Fields{
						"err": err,
					}).Error("error requesting vote")
					return
				}

				log.WithFields(log.Fields{
					"id":      r.state.ID,
					"address": memberAddress,
					"resp":    resp,
				}).Infof("do voting done")
				voteRespChan <- resp
			}(memberAddress)
		}
	}

	return voteRespChan
}

func (r *defaultRaft) doLeaderAction() {
	log.WithFields(log.Fields{}).Infof("do leader action start")

	for {
		select {
		case <-r.sendAppendEntriesSignal:
			log.Debug("append entries sent")
		case <-time.After(r.config.IdleTimeout):
			log.Debug("idle timeout")
			for _, memberAddress := range r.config.ClusterMemberAddreses {
				if memberAddress != r.address {
					go r.raftClient.AppendEntries(memberAddress, model.NewAppendEntriesRequestDTO())
				}
			}
		}
	}
}

// for debugging purposes
func (r *defaultRaft) watcher() {
	ticker := time.NewTicker(time.Millisecond * 2000)

	for {
		select {
		case <-ticker.C:
			log.WithFields(log.Fields{
				"id":        r.state.ID,
				"role":      r.state.Role,
				"term":      r.state.CurrentTerm,
				"voted_for": r.state.VotedFor,
			}).Debug("raft node status")
		}
	}
}
