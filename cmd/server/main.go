package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/adylanrff/raft-algorithm/raft"
	"github.com/adylanrff/raft-algorithm/server"
	"github.com/adylanrff/raft-algorithm/util"
	log "github.com/sirupsen/logrus"
)

var port int
var logPath string

func init() {
	flag.IntVar(&port, "port", 8000, "server port")
	flag.StringVar(&logPath, "log_path", "server.log", "log path")
}

func main() {
	util.InitLogger(logPath)

	address := fmt.Sprintf("127.0.0.1:%d", port)
	serverID := raft.GetRaftServerIDFromAddress(address)

	raftHandler := raft.NewRaft(serverID, raft.RaftConfig{
		ElectionTimeout:   500 * time.Millisecond,
		HeartbeatInterval: 200 * time.Millisecond,
		ClusterMemberAddreses: []string{
			"127.0.0.1:8000",
			"127.0.0.1:8001",
			"127.0.0.1:8002",
		},
	})

	raftServerHandler := raft.NewRaftServerhandler(raftHandler)

	server := server.NewServer(address)
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
