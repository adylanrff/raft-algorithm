package main

import (
	"flag"
	"os"

	"github.com/adylanrff/raft-algorithm/raft"
	log "github.com/sirupsen/logrus"
)

var port int
var logPath string

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)

	flag.IntVar(&port, "port", 8000, "server port")
	flag.StringVar(&logPath, "log_path", "server.log", "log path")
}

func main() {
	f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)

	raftServer := raft.NewServer(port)
	raftServer.Run()
}
