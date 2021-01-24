package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	state "github.com/korjavin/pbuf-example/state"
	"log"
	"os"
	"time"

	messages "github.com/korjavin/pbuf-example/messages"
	"google.golang.org/protobuf/proto"
)

func main() {
	st := newState()
	displayPosts(st)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var m messages.Direct
		hexMsg := scanner.Bytes()
		msg := make([]byte, hex.DecodedLen(len(hexMsg)))
		_, err := hex.Decode(msg, hexMsg)
		if err != nil {
			log.Fatalf("[FATAL]  %v", err)
		}

		err = proto.Unmarshal(msg, &m)
		if err != nil {
			log.Fatalf("[FATAL]  %v", err)
		}
		addPostToState(st, &m)
		displayPosts(st)
	}

	if scanner.Err() != nil {
		log.Fatalf("[FATAL]  %v", scanner.Err())
	}
}

func newState() *state.State {
	return &state.State{
		Posts:             make(map[string]*state.ListOfPosts),
		ScheduledMessages: make([]*messages.Text, 0),
		DirectMessages:    make(map[string]*state.ListOfDirects),
	}
}

func addPostToState(st *state.State, m *messages.Direct) {
	post := &state.Post{
		Timestamp: time.Now().Unix(),
		Text:      m.Text,
	}

	if posts, ok := st.Posts[m.Account]; ok {
		st.Posts[m.Account].Posts = append(posts.Posts, post)

		return
	}

	st.Posts[m.Account] = &state.ListOfPosts{Posts: []*state.Post{post}}
}

func displayPosts(st *state.State) {
	if len(st.Posts) == 0 {
		fmt.Printf("No posts yet. But you can be first!\n")
		return
	}

	fmt.Printf("\nList of posts:\n\n")
	for k, v := range st.Posts {
		fmt.Printf("%s's posts:\n", k)
		for i, post := range v.Posts {
			fmt.Printf(
				"%d) %s (Published at %s)\n",
				i+1,
				post.Text,
				time.Unix(post.Timestamp, 0).Format(time.RFC1123),
			)
		}
		fmt.Printf("--End--\n\n")
	}
}
