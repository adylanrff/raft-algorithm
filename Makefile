all: server client

update-proto:
	protoc --proto_path=./proto --go_out=proto --go_opt=paths=source_relative raft/raft.proto
	protoc --proto_path=./proto --go_out=proto --go_opt=paths=source_relative server/server.proto


server:
	go build -o bin/raft cmd/server/main.go

client:
	go build -o bin/client cmd/client/main.go
