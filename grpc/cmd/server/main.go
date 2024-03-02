package main

import (
	"context"
	"flag"
	"github.com/igorakimy/bigtech_microservices/internal/app"
	"log"
)

func main() {
	flag.Parse()

	application, err := app.NewApp(context.Background())
	if err != nil {
		log.Fatalf("failed to init app: %v", err.Error())
	}

	if err = application.Run(); err != nil {
		log.Fatalf("failed to run app: %v", err.Error())
	}
}
