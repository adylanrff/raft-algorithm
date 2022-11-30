package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/adylanrff/raft-algorithm/raft"
	"github.com/adylanrff/raft-algorithm/rpc"
	"github.com/adylanrff/raft-algorithm/util"
	log "github.com/sirupsen/logrus"
)

var port int
var logPath string

func init() {
	flag.IntVar(&port, "port", 8000, "server port")
	flag.StringVar(&logPath, "log_path", "server.log", "log path")

	flag.Parse()
	rand.Seed(time.Now().UnixNano())
}

func main() {
	util.InitLogger(logPath)

	address := fmt.Sprintf("127.0.0.1:%d", port)
	serverID := raft.GetRaftServerIDFromAddress(address)

	raftHandler := raft.NewRaft(serverID, address, raft.RaftConfig{
		ElectionTimeout: 1000 * time.Millisecond,
		IdleTimeout:     500 * time.Millisecond,
		ClusterMemberAddreses: []string{
			"127.0.0.1:8000",
			"127.0.0.1:8001",
			"127.0.0.1:8002",
			"127.0.0.1:8003",
			"127.0.0.1:8004",
		},
	})

	raftServerHandler := raft.NewRaftServerhandler(raftHandler)

	server := rpc.NewServer(address)
	server.AddHandler(raft.RaftMethodName_RequestVotes, raftServerHandler.RequestVoteHandler)
	server.AddHandler(raft.RaftMethodName_AppendEntries, raftServerHandler.AppendEntriesHandler)

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		server.Run()
	}()

	go func() {
		defer wg.Done()
		raftHandler.Run()
	}()

	wg.Wait()

	log.Debug("server terminated")
}
