package main

import (
	"context"
	"fmt"
	pb "github.com/bend-is/pbuf-example/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("[FATAL] fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewPostsAPIClient(conn)

	stream, err := client.Subscribe(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatalf("[FATAL] fail to listen stream: %v", err)
	}

	for {
		post, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("[FATAL] fail to listen stream: %v", err)
		}

		fmt.Printf("New post from user with id %s:\n\n", post.UserId)
		fmt.Printf("%s\n", post.Text)
		fmt.Printf("---%s---\n\n", time.Unix(post.PublishedAt.Seconds, 0).Format(time.ANSIC))
	}
}
