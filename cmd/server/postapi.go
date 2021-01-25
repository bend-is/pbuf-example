package main

import (
	"context"
	"errors"
	pb "github.com/bend-is/pbuf-example/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
)

type PostAPIServer struct {
	mu          sync.Mutex
	subscribers []chan *pb.Post
	pb.UnimplementedPostsAPIServer
}

func NewPostAPIServer() *PostAPIServer {
	return &PostAPIServer{}
}

func (p *PostAPIServer) GetPosts(ctx context.Context, emp *empty.Empty) (*pb.ListOfPosts, error) {
	posts := &pb.ListOfPosts{
		Posts: make([]*pb.Post, len(serverState.Posts)),
	}

	copy(posts.Posts, serverState.Posts)

	return posts, nil
}

func (p *PostAPIServer) Subscribe(emp *empty.Empty, stream pb.PostsAPI_SubscribeServer) error {
	ch := make(chan *pb.Post)
	defer p.unsubscribe(ch)

	p.mu.Lock()
	p.subscribers = append(p.subscribers, ch)
	p.mu.Unlock()

	for {
		select {
		case post := <-ch:
			if err := stream.Send(post); err != nil {
				return err
			}
		}
	}
}

func (p *PostAPIServer) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error) {
	if req.Text == "" {
		return nil, errors.New("empty text was specified")
	}
	if req.UserId == "" {
		return nil, errors.New("no user was specified")
	}

	serverState.Lock()
	defer serverState.Unlock()

	if len(serverState.Posts) == 0 {
		serverState.Posts = make([]*pb.Post, 0, 1)
	}

	id, err := generateNewID()
	if err != nil {
		return nil, err
	}

	newPost := &pb.Post{
		Id:          id,
		ReplayId:    req.ReplayId,
		Text:        req.Text,
		UserId:      req.UserId,
		PublishedAt: timestamppb.Now(),
	}
	serverState.Posts = append(serverState.Posts, newPost)

	p.mu.Lock()
	for _, ch := range p.subscribers {
		go func(ch chan *pb.Post) { ch <- newPost }(ch)
	}
	p.mu.Unlock()

	return newPost, nil
}

func (p *PostAPIServer) unsubscribe(ch chan *pb.Post) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i, v := range p.subscribers {
		if v == ch {
			p.subscribers = append(p.subscribers[:i], p.subscribers[i+1:]...)
			break
		}
	}

	close(ch)
}
