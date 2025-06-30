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

	// Prepare chat request
	chatRequest := &gofabric.ChatRequest{
		Prompts: []gofabric.PromptRequest{
			{
				UserInput:   "Write a Golang function that calculates the factorial of a number.",
				Vendor:      "Gemini",
				Model:       "gemini-2.0-flash",
				ContextName: "",
				PatternName: "coding_master",
			},
		},
		Language:    "en",
		ChatOptions: gofabric.ChatOptions{},
	}

	// Start streaming chat
	responses, err := client.Chat(ctx, chatRequest)
	if err != nil {
		log.Fatalf("Error sending chat request: %v", err)
	}

	for response := range responses {
		switch response.Type {
		case string(gofabric.StreamResponseTypeContent):
			log.Printf("Content (%s): %s\n", response.Format, response.Content)
		case string(gofabric.StreamResponseTypeError):
			log.Printf("Error: %s\n", response.Content)
		case string(gofabric.StreamResponseTypeComplete):
			log.Printf("Chat completed")
		}
	}
}
