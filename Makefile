all: build

update-proto:
	protoc --proto_path=./proto --go_out=proto --go_opt=paths=source_relative raft.proto

build:
	go build -o bin/raft cmd/main.go
