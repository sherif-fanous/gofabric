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

	sessionName := "session_example"

	// Check if the session exists
	if ok, err := client.SessionExists(ctx, sessionName); err != nil {
		log.Fatalf("Error checking session existence: %v\n", err)
	} else if ok {
		log.Printf("Session '%s' already exists\n", sessionName)
	} else {
		log.Printf("Session '%s' does not exist, proceeding to create it\n", sessionName)
	}

	// Create a new session
	err := client.CreateSession(ctx, sessionName, strings.NewReader("[]"))
	if err != nil {
		log.Fatalf("Error creating session: %v\n", err)
	}
	log.Printf("Session '%s' created successfully\n", sessionName)

	// Check again if the session exists
	if ok, err := client.SessionExists(ctx, sessionName); err != nil {
		log.Printf("Error checking session existence: %v\n", err)

		return
	} else if ok {
		log.Printf("Session '%s' already exists\n", sessionName)
	} else {
		log.Printf("Session '%s' does not exist, proceeding to create it\n", sessionName)
	}

	// List all sessions
	sessions, err := client.ListSessions(ctx)
	if err != nil {
		log.Fatalf("Error listing sessions: %v\n", err)
	}
	log.Printf("Available sessions: %v\n", sessions)

	// Fetch the session metadata
	sessionMetadata, err := client.GetSessionMetadata(ctx, sessionName)
	if err != nil {
		log.Fatalf("Error fetching session metadata: %v\n", err)
	}
	log.Printf("Fetched session: %+v\n", sessionMetadata)

	// Rename the session
	if err := client.RenameSession(ctx, sessionName, "renamed_"+sessionName); err != nil {
		log.Fatalf("Error renaming session: %v\n", err)
	}
	log.Printf("Session '%s' renamed to 'renamed_%s'\n", sessionName, sessionName)

	// Delete the session
	if err := client.DeleteSession(ctx, "renamed_"+sessionName); err != nil {
		log.Fatalf("Error deleting session: %v\n", err)
	}
	log.Printf("Session 'renamed_%s' deleted successfully\n", sessionName)
}
