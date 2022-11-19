package main

import (
	"flag"
	"fmt"
	"net"

	pb "github.com/adylanrff/raft-algorithm/proto"
	"github.com/adylanrff/raft-algorithm/raft"
	"github.com/adylanrff/raft-algorithm/raft/model"
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

	raftMsg := model.RaftMessageDTO{
		RaftMessage: &pb.RaftMessage{
			Method:          raft.RaftMethodName_RequestVotes,
			RaftMessageType: pb.RaftMessageType_RaftMessageType_RaftRequest,
			Payload: &pb.RaftMessage_RaftRequest{
				RaftRequest: &pb.RaftRequest{
					Request: &pb.RaftRequest_RequestVoteRequest{
						RequestVoteRequest: &pb.RequestVoteRequest{},
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

	raftMsgDTO, err := raft.ParseRaftMessage(conn)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", raftMsgDTO)
}
