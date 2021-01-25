package main

import (
	pb "github.com/bend-is/pbuf-example/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

var serverState *ServerState

type ServerState struct {
	sync.RWMutex
	*pb.State
}

func main() {
	serverState = newServer()
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("[FATAL] failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterPostsAPIServer(grpcServer, NewPostAPIServer())
	pb.RegisterUsersAPIServer(grpcServer, NewUserAPIServer())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[FATAL] grpc server error: %v", err)
	}
}

func newServer() *ServerState {
	return &ServerState{
		State: &pb.State{
			Posts: make([]*pb.Post, 0),
			Users: make([]*pb.User, 0),
		},
	}
}
