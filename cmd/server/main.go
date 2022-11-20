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

	raftServer := server.NewServer(port)
	raftServer.AddHandler(raft.RaftMethodName_RequestVotes, func(req *server.ServerMessageDTO) (resp *server.ServerMessageDTO, err error) {})
	raftServer.Run()
}
