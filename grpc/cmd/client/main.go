package main

import (
	"context"
	desc "github.com/igorakimy/bigtech_microservices/pkg/note/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"time"
)

const (
	address = "localhost:50051"
	noteID  = 42
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("service.pem", "")
	if err != nil {
		log.Fatalf("could not process the credentials: %v", err)
	}

	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(creds),
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
