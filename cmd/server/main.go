package main

import (
	"flag"

	"github.com/adylanrff/raft-algorithm/raft"
	"github.com/adylanrff/raft-algorithm/server"
	"github.com/adylanrff/raft-algorithm/util"
)

var port int
var logPath string

func init() {
	flag.IntVar(&port, "port", 8000, "server port")
	flag.StringVar(&logPath, "log_path", "server.log", "log path")
}

func main() {
	util.InitLogger(logPath)

	server := server.NewServer(port)

	raftHandler := raft.NewRaftHandler()
	raftServerHandler := raft.NewRaftServerhandler(raftHandler)

	server.AddHandler(raft.RaftMethodName_RequestVotes, raftServerHandler.RequestVoteHandler)
	server.AddHandler(raft.RaftMethodName_AppendEntries, raftServerHandler.AppendEntriesHandler)
	server.Run()
}
