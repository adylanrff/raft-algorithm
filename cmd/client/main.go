package main

import (
	"flag"
	"fmt"
	"net"

	raftPb "github.com/adylanrff/raft-algorithm/proto/raft"
	serverPb "github.com/adylanrff/raft-algorithm/proto/server"
	"github.com/adylanrff/raft-algorithm/raft"
	"github.com/adylanrff/raft-algorithm/server"
	"github.com/adylanrff/raft-algorithm/util"
)

var port int
var logPath string

func init() {
	flag.IntVar(&port, "port", 8000, "server target port")
	flag.StringVar(&logPath, "log_path", "server.log", "log path")
}

func main() {
	util.InitLogger(logPath)

	raftMsg := server.ServerMessageDTO{
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

	conn, err := net.Dial("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	byteRaftMsg, err := raftMsg.ToBytes()
	if err != nil {
		panic(err)
	}

	_, err = conn.Write(byteRaftMsg)
	if err != nil {
		panic(err)
	}

	raftMsgDTO, err := server.ParseServerMessage(conn)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", raftMsgDTO)
}
