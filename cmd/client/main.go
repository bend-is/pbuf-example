package main

import (
	"context"
	"fmt"
	pb "github.com/bend-is/pbuf-example/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"sync"
	"time"
)

var texts = [...]string{
	"Hello there. Thanks for the follow. Did you notice, that I am an egg? A talking egg? Damn!",
	"Thanks mate! Feel way better now",
	"Yeah that is crazy",
	"Hi",
	"Thanks",
	"Okay",
}

var logins = [...]string{
	"Super_Cool_Guy",
	"HotPepper228",
	"DemonSoul",
	"Ratchet",
	"Clank",
	"Kratos",
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("[FATAL] fail to dial: %v", err)
	}
	defer conn.Close()

	userClient := pb.NewUsersAPIClient(conn)
	postClient := pb.NewPostsAPIClient(conn)

	user := register(context.Background(), userClient, logins[rand.Intn(len(logins))])
	displayCurrentUsers(context.Background(), userClient)

	leaveAPosts(context.Background(), postClient, user)
	displayCurrentPosts(context.Background(), postClient)
}

func register(ctx context.Context, client pb.UsersAPIClient, login string) *pb.User {
	fmt.Println("Trying to register...")

	user, err := client.CreateUser(context.Background(), &pb.CreateUserRequest{Login: login})
	if err != nil {
		log.Fatalf("[FATAL] faild to register user: %v", err)
	}

	fmt.Printf("\nSuccessfully register user [%s] with id [%s]\n", user.Login, user.Id)

	return user
}

func displayCurrentUsers(ctx context.Context, client pb.UsersAPIClient) {
	userList, err := client.GetUsers(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatalf("[FATAL] faild to get a list if users: %v", err)
	}

	fmt.Printf("\nList of users:\n")

	for i, u := range userList.Users {
		fmt.Printf("%d) %s\n", i+1, u.Login)
	}
}

func leaveAPosts(ctx context.Context, client pb.PostsAPIClient, user *pb.User) {
	fmt.Println("\nTrying to leave a posts...")

	var wg sync.WaitGroup
	for i := 0; i < rand.Intn(len(texts)); i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			_, err := client.CreatePost(ctx, &pb.CreatePostRequest{
				ReplayId: "",
				Text:     texts[i],
				UserId:   user.Id,
			})
			if err != nil {
				log.Printf("[ERROR] error while leaving a post: %v", err)
			}
		}(i)
	}
	wg.Wait()
}

func displayCurrentPosts(ctx context.Context, client pb.PostsAPIClient) {
	postList, err := client.GetPosts(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatalf("[FATAL] faild to get a list if posts: %v", err)
	}

	fmt.Printf("\nList of posts:\n")
	for i, p := range postList.Posts {
		fmt.Printf(
			"%d) %s (Published at %s by %s)\n",
			i+1,
			p.Text,
			time.Unix(p.PublishedAt.Seconds, 0).Format(time.ANSIC),
			p.UserId,
		)
	}
}
