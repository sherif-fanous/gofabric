package main

import (
	"context"
	"log"
	"os"

	"github.com/sherif-fanous/gofabric"
)

func main() {
	fabricServerURL, ok := os.LookupEnv("FABRIC_SERVER_URL")
	if !ok {
		log.Fatal("FABRIC_SERVER_URL environment variable is not set")
	}

	client := gofabric.NewClient(fabricServerURL)
	ctx := context.Background()

	// List available strategies
	strategys, err := client.ListStrategies(ctx)
	if err != nil {
		log.Fatalf("Error listing strategies: %v\n", err)
	}
	log.Printf("Available strategies: %v\n", strategys)
}
