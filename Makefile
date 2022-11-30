.PHONY: all

all: raft-server client

update-proto:
	protoc --proto_path=./proto --go_out=proto/raft --go_opt=paths=source_relative raft.proto
	protoc --proto_path=./proto --go_out=proto/server --go_opt=paths=source_relative server.proto

raft-server:
	go build -o bin/raft cmd/server/main.go

client:
	go build -o bin/client cmd/client/main.go

run-3-server:
	echo "running on port 8000" && ./bin/raft -port=8000 -log_path=server1.log &
	echo "running on port 8001" && ./bin/raft -port=8001 -log_path=server2.log &
	echo "running on port 8002" && ./bin/raft -port=8002 -log_path=server3.log &

run-5-server:
	echo "running on port 8000" && ./bin/raft -port=8000 -log_path=server1.log &
	echo "running on port 8001" && ./bin/raft -port=8001 -log_path=server2.log &
	echo "running on port 8002" && ./bin/raft -port=8002 -log_path=server3.log &
	echo "running on port 8003" && ./bin/raft -port=8003 -log_path=server4.log &
	echo "running on port 8004" && ./bin/raft -port=8004 -log_path=server5.log &

clean: 
	rm -rf bin
	rm *.log
