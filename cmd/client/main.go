package main

import (
	"flag"
	"fmt"

	raftPb "github.com/adylanrff/raft-algorithm/proto/raft"
	serverPb "github.com/adylanrff/raft-algorithm/proto/server"
	"github.com/adylanrff/raft-algorithm/raft"
	"github.com/adylanrff/raft-algorithm/rpc"
	"github.com/adylanrff/raft-algorithm/util"
)

var port int
var logPath string

func init() {
	flag.IntVar(&port, "port", 8000, "server target port")
	flag.StringVar(&logPath, "log_path", "server.log", "log path")

	flag.Parse()
}

func main() {
	util.InitLogger(logPath)

	raftMsg := rpc.ServerMessageDTO{
		ServerMessage: &serverPb.ServerMessage{
			Method: raft.RaftMethodName_RequestVotes,
			Payload: &serverPb.ServerMessage_ServerRequest{
				ServerRequest: &serverPb.ServerRequest{
					Request: &serverPb.ServerRequest_RequestVoteRequest{
						RequestVoteRequest: &raftPb.RequestVoteRequest{},
					},
				},
			},
		},
	}

	rpcClient := rpc.NewRPCClient()
	resp, err := rpcClient.Call(fmt.Sprintf("0.0.0.0:%d", port), &raftMsg)

	fmt.Printf("resp=%+v,err=%v", resp, err)
}
