package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/sherif-fanous/gofabric"
)

func main() {
	fabricServerURL, ok := os.LookupEnv("FABRIC_SERVER_URL")
	if !ok {
		log.Fatal("FABRIC_SERVER_URL environment variable is not set")
	}

	client := gofabric.NewClient(fabricServerURL)
	ctx := context.Background()

	contextName := "context_example.md"

	// Check if the context exists
	if ok, err := client.ContextExists(ctx, contextName); err != nil {
		log.Fatalf("Error checking context existence: %v\n", err)
	} else if ok {
		log.Printf("Context '%s' already exists\n", contextName)
	} else {
		log.Printf("Context '%s' does not exist, proceeding to create it\n", contextName)
	}

	// Create a new context
	err := client.CreateContext(ctx, contextName, strings.NewReader("Context content"))
	if err != nil {
		log.Fatalf("Error creating context: %v\n", err)
	}
	log.Printf("Context '%s' created successfully\n", contextName)

	// Check again if the context exists
	if ok, err := client.ContextExists(ctx, contextName); err != nil {
		log.Printf("Error checking context existence: %v\n", err)

		return
	} else if ok {
		log.Printf("Context '%s' already exists\n", contextName)
	} else {
		log.Printf("Context '%s' does not exist, proceeding to create it\n", contextName)
	}

	// List all contexts
	contexts, err := client.ListContexts(ctx)
	if err != nil {
		log.Fatalf("Error listing contexts: %v\n", err)
	}
	log.Printf("Available contexts: %v\n", contexts)

	// Fetch the context metadata
	contextMetadata, err := client.GetContextMetadata(ctx, contextName)
	if err != nil {
		log.Fatalf("Error fetching context metadata: %v\n", err)
	}
	log.Printf("Fetched context: %+v\n", contextMetadata)

	// Rename the context
	if err := client.RenameContext(ctx, contextName, "renamed_"+contextName); err != nil {
		log.Fatalf("Error renaming context: %v\n", err)
	}
	log.Printf("Context '%s' renamed to 'renamed_%s'\n", contextName, contextName)

	// Delete the context
	if err := client.DeleteContext(ctx, "renamed_"+contextName); err != nil {
		log.Fatalf("Error deleting context: %v\n", err)
	}
	log.Printf("Context 'renamed_%s' deleted successfully\n", contextName)
}
