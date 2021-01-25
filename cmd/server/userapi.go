package main

import (
	"context"
	"errors"
	pb "github.com/bend-is/pbuf-example/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

type UserAPIServer struct {
	pb.UnimplementedUsersAPIServer
}

func NewUserAPIServer() *UserAPIServer {
	return &UserAPIServer{}
}

func (u *UserAPIServer) GetUsers(ctx context.Context, emp *empty.Empty) (*pb.ListOfUsers, error) {
	users := &pb.ListOfUsers{
		Users: make([]*pb.User, len(serverState.Users)),
	}

	copy(users.Users, serverState.Users)

	return users, nil
}

func (u *UserAPIServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	if req.Login == "" {
		return nil, errors.New("empty login was specified")
	}

	serverState.Lock()
	defer serverState.Unlock()

	if len(serverState.Users) == 0 {
		serverState.Users = make([]*pb.User, 0, 1)
	}

	id, err := generateNewID()
	if err != nil {
		return nil, err
	}

	newUser := &pb.User{Id: id, Login: req.Login}
	serverState.Users = append(serverState.Users, newUser)

	return newUser, nil
}
