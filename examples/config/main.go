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

	// Fetch the configuration
	config, err := client.GetConfig(ctx)
	if err != nil {
		log.Fatalf("Error fetching configuration: %v\n", err)
	}
	log.Printf("Configuration: %v\n", config)
}
