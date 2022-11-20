package raft

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
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

	// timer
	electionTimer   *time.Timer
	heartbeatTicker *time.Ticker

	raftClient RaftClient

	appendEntriesSignal chan struct{}
}

func NewRaft(id string, address string, cfg RaftConfig) *defaultRaft {
	electionTimeout := rand.Intn(int(cfg.ElectionTimeout*2)-int(cfg.ElectionTimeout)) + int(cfg.ElectionTimeout)
	cfg.ElectionTimeout = time.Duration(electionTimeout)

	fmt.Println(electionTimeout)
	return &defaultRaft{
		address:       address,
		state:         model.NewRaftState(id),
		config:        cfg,
		electionTimer: time.NewTimer(cfg.ElectionTimeout),
		raftClient:    NewRaftClient(rpc.NewRPCClient()),

		appendEntriesSignal: make(chan struct{}),
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
	r.appendEntriesSignal <- struct{}{}
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
	log.WithFields(log.Fields{}).Infof("do candidate action startzzz")
	electionChan := make(chan struct{}, 1)
	electionChan <- struct{}{}

	log.WithFields(log.Fields{}).Infof("do candidate action send election")

	for {
		select {
		case <-r.electionTimer.C:
			log.WithFields(log.Fields{}).Infof("reset election timer")
			// reset election timer
			r.resetElectionTimer()
			// do another election
			electionChan <- struct{}{}
		case <-r.appendEntriesSignal:
			log.WithFields(log.Fields{}).Infof("heartbeat received")
			r.state.ChangeRole(model.RaftRoleFollower)
			return
		case <-electionChan:
			log.WithFields(log.Fields{}).Infof("do election")
			shouldReturn := r.doElection()
			if shouldReturn {
				return
			}
		}
	}
}

func (r *defaultRaft) doElection() bool {
	log.WithFields(log.Fields{}).Infof("do election start")

	r.state.Lock()
	defer r.state.Unlock()

	// 1. increment current term
	r.state.CurrentTerm++
	// 2. vote for self
	r.state.VotedFor = r.state.ID
	// 3. reset election timer
	r.resetElectionTimer()
	// calculate majority:
	majority := math.Ceil(float64(len(r.config.ClusterMemberAddreses)) / float64(2))
	// 1 because we voted for ourselves
	var votes int32 = 1

	// Do voting
	log.WithFields(log.Fields{}).Infof("do voting")
	requestVoteReq := model.NewRequestVoteRequestDTO(
		r.state.CurrentTerm,
		r.state.ID,
		int32(len(r.state.Logs)),
		0)

	voteRespChan := make(chan *model.RequestVoteResponseDTO, len(r.config.ClusterMemberAddreses))

	var wg sync.WaitGroup
	for _, memberAddress := range r.config.ClusterMemberAddreses {
		if memberAddress != r.address {
			wg.Add(1)
			go func(memberAddress string) {
				defer wg.Done()
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

	wg.Add(1)
	go func() {
		defer wg.Done()

		voteSeen := 0
		for voteResp := range voteRespChan {
			voteSeen++
			if voteResp.VoteGranted {
				atomic.AddInt32(&votes, 1)
			}
			if voteSeen == len(r.config.ClusterMemberAddreses)-1 {
				return
			}
		}
	}()

	wg.Wait()

	log.WithFields(log.Fields{
		"id":             r.state.ID,
		"voting_result":  votes,
		"majority":       majority,
		"will_be_leader": votes >= int32(majority),
		"term":           r.state.CurrentTerm,
	}).Infof("voting result")

	if votes >= int32(majority) {
		r.state.ChangeRole(model.RaftRoleLeader)
		return true
	}

	return false
}

func (r *defaultRaft) doLeaderAction() {
	log.WithFields(log.Fields{}).Infof("do leader action start")
	r.heartbeatTicker = time.NewTicker(r.config.HeartbeatInterval)

	for range r.heartbeatTicker.C {
		for _, memberAddress := range r.config.ClusterMemberAddreses {
			if memberAddress != r.address {
				go r.raftClient.AppendEntries(memberAddress, model.NewAppendEntriesRequestDTO())
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

func (r *defaultRaft) resetElectionTimer() {
	log.WithFields(log.Fields{}).Infof("reset election timer")

	r.electionTimer.Reset(r.config.ElectionTimeout)
}
