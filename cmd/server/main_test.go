package main

import (
	"encoding/json"
	"github.com/korjavin/pbuf-example/messages"
	"github.com/korjavin/pbuf-example/state"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

func BenchmarkOrderProtoMarshal(b *testing.B) {
	obj := getTestData()

	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(obj)
		if err != nil {
			b.Fatal("Marshaling error:", err)
		}
	}
}

func BenchmarkOrderJSONMarshal(b *testing.B) {
	obj := getTestData()

	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(obj)
		if err != nil {
			b.Fatal("Marshaling error:", err)
		}
	}
}

func BenchmarkOrderProtoUnmarshal(b *testing.B) {
	obj := getTestData()

	data, err := proto.Marshal(obj)
	if err != nil {
		b.Fatal("Marshaling error:", err)
	}
	for i := 0; i < b.N; i++ {
		var res state.State
		err := proto.Unmarshal(data, &res)
		if err != nil {
			b.Fatal("Unmarshalling error:", err)
		}
	}
}

func BenchmarkOrderJSONUnmarshal(b *testing.B) {
	obj := getTestData()

	data, err := json.Marshal(obj)
	if err != nil {
		b.Fatal("Marshaling error:", err)
	}
	for i := 0; i < b.N; i++ {
		var res state.State
		err := json.Unmarshal(data, &res)
		if err != nil {
			b.Fatal("Unmarshalling error:", err)
		}
	}
}

func getTestData() *state.State {
	tmpst := time.Now().Unix()

	return &state.State{
		Posts: map[string]*state.ListOfPosts{
			"key1": {Posts: []*state.Post{
				{Timestamp: tmpst, Text: "Some text number one"},
				{Timestamp: tmpst, Text: "Some text number two"},
				{Timestamp: tmpst, Text: "Some text number three"},
			}},
			"key2": {Posts: []*state.Post{
				{Timestamp: tmpst, Text: "Some text number one"},
				{Timestamp: tmpst, Text: "Some text number two"},
				{Timestamp: tmpst, Text: "Some text number three"},
			}},
			"key3": {Posts: []*state.Post{
				{Timestamp: tmpst, Text: "Some text number one"},
				{Timestamp: tmpst, Text: "Some text number two"},
				{Timestamp: tmpst, Text: "Some text number three"},
			}},
		},
		ScheduledMessages: []*messages.Text{
			{Text: "Message text number one", PublishAt: tmpst},
			{Text: "Message text number two", PublishAt: tmpst},
			{Text: "Message text number three", PublishAt: tmpst},
			{Text: "Message text number five", PublishAt: tmpst},
		},
		DirectMessages: map[string]*state.ListOfDirects{
			"key1": {Directs: []*state.Direct{
				{From: "Account_1", Text: "Some text number one"},
				{From: "Account_2", Text: "Some text number two"},
				{From: "Account_3", Text: "Some text number three"},
			}},
			"key2": {Directs: []*state.Direct{
				{From: "Account_1", Text: "Some text number one"},
				{From: "Account_2", Text: "Some text number two"},
				{From: "Account_3", Text: "Some text number three"},
			}},
			"key3": {Directs: []*state.Direct{
				{From: "Account_1", Text: "Some text number one"},
				{From: "Account_2", Text: "Some text number two"},
				{From: "Account_3", Text: "Some text number three"},
			}},
		},
	}
}
