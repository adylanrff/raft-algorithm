.PHONY: all

all: clean raft-server client

update-proto:
	protoc --proto_path=./proto --go_out=proto/raft --go_opt=paths=source_relative raft.proto
	protoc --proto_path=./proto --go_out=proto/server --go_opt=paths=source_relative server.proto

raft-server:
	go build -o bin/raft cmd/server/main.go

client:
	go build -o bin/client cmd/client/main.go

clean: 
	rm -rf bin