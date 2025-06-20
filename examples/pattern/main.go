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

	patternName := "pattern_example"

	// Check if the pattern exists
	if ok, err := client.PatternExists(ctx, patternName); err != nil {
		log.Fatalf("Error checking pattern existence: %v\n", err)
	} else if ok {
		log.Printf("Pattern '%s' already exists\n", patternName)
	} else {
		log.Printf("Pattern '%s' does not exist, proceeding to create it\n", patternName)
	}

	// Create a new pattern
	err := client.CreatePattern(ctx, patternName, strings.NewReader("Pattern content"))
	if err != nil {
		log.Fatalf("Error creating pattern: %v\n", err)
	}
	log.Printf("Pattern '%s' created successfully\n", patternName)

	// Check again if the pattern exists
	if ok, err := client.PatternExists(ctx, patternName); err != nil {
		log.Printf("Error checking pattern existence: %v\n", err)

		return
	} else if ok {
		log.Printf("Pattern '%s' already exists\n", patternName)
	} else {
		log.Printf("Pattern '%s' does not exist, proceeding to create it\n", patternName)
	}

	// List all patterns
	patterns, err := client.ListPatterns(ctx)
	if err != nil {
		log.Fatalf("Error listing patterns: %v\n", err)
	}
	log.Printf("Available patterns: %v\n", patterns)

	// Fetch the pattern metadata
	patternMetadata, err := client.GetPatternMetadata(ctx, patternName)
	if err != nil {
		log.Fatalf("Error fetching pattern metadata: %v\n", err)
	}
	log.Printf("Fetched pattern: %+v\n", patternMetadata)

	// Rename the pattern
	if err := client.RenamePattern(ctx, patternName, "renamed_"+patternName); err != nil {
		log.Fatalf("Error renaming pattern: %v\n", err)
	}
	log.Printf("Pattern '%s' renamed to 'renamed_%s'\n", patternName, patternName)

	// Delete the pattern
	if err := client.DeletePattern(ctx, "renamed_"+patternName); err != nil {
		log.Fatalf("Error deleting pattern: %v\n", err)
	}
	log.Printf("Pattern 'renamed_%s' deleted successfully\n", patternName)
}
