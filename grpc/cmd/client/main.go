package main

import (
	"context"
	desc "github.com/igorakimy/bigtech_microservices/pkg/note/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const (
	address = "localhost:50051"
	noteID  = 12
)

func main() {
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer func() { _ = conn.Close() }()

	client := desc.NewNoteV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.Get(ctx, &desc.GetRequest{
		Id: noteID,
	})
	if err != nil {
		log.Fatalf("failed to get note by id: %v", err)
	}

	log.Printf("Note info:%v\n", r.GetNote())
}
