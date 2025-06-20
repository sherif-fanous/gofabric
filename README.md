# gofabric

`gofabric` is a Go client library for interacting with the [Fabric](https://github.com/danielmiessler/fabric) API server. It provides a convenient way to manage contexts, patterns, and sessions, initiate chat sessions, and retrieve configuration and available models/strategies from the Fabric API.

## Features

- **Client Initialization**: Create a client with API key and custom HTTP client options.
- **Chat Functionality**: Initiate chat sessions with streaming responses using Server-Sent Events (SSE).
- **Entity Management**: Create, delete, retrieve, list, and rename `contexts`, `patterns`, and `sessions`.
- **Configuration Management**: Get and update the Fabric API server configuration.
- **Model and Strategy Listing**: Retrieve lists of available models and strategies.

## Installation

To install `gofabric`, use `go get`:

```bash
go get github.com/sherif-fanous/gofabric
```

## Usage

### Basic Client Initialization

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/sherif-fanous/gofabric"
)

func main() {
    // Replace with your Fabric API server host
    host := "http://localhost:8000"
    // Replace with your API key if required by the server
    apiKey := "your-api-key"

    client := gofabric.NewClient(
        host,
        gofabric.WithAPIKey(apiKey),
    )

    // Example: Get server configuration
    config, err := client.GetConfig(context.Background())
    if err != nil {
        log.Fatalf("Failed to get config: %v", err)
    }
    fmt.Printf("Fabric API Config: %+v\n", config)
}
```

### Chatting with the API

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/sherif-fanous/gofabric"
)

func main() {
    host := "http://localhost:8000"
    apiKey := "your-api-key"

    client := gofabric.NewClient(
        host,
        gofabric.WithAPIKey(apiKey),
    )

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
    responses, err := client.Chat(context.Background(), chatRequest)
    if err != nil {
        log.Fatalf("Error sending chat request: %v", err)
    }

    for response := range responses {
        switch response.Type {
        case "content":
            log.Printf("Content (%s): %s\n", response.Format, response.Content)
        case "error":
            log.Printf("Error: %s\n", response.Content)
        case "complete":
            log.Printf("Chat completed")
        }
    }
}
```

### Managing Entities (Contexts, Patterns, Sessions)

The client provides methods for `Context`, `Pattern`, and `Session` management. Here's an example for `Context`:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "strings"

    "github.com/sherif-fanous/gofabric"
)

func main() {
    host := "http://localhost:8000"
    apiKey := "your-api-key"

    client := gofabric.NewClient(
        host,
        gofabric.WithAPIKey(apiKey),
    )

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
```

Similar methods are available for `Patterns` and `Sessions`:

- `CreatePattern`, `DeletePattern`, `PatternExists`, `GetPatternMetadata`, `ListPatterns`, `RenamePattern`
- `CreateSession`, `DeleteSession`, `SessionExists`, `GetSessionMetadata`, `ListSessions`, `RenameSession`

For more detailed examples on how to use the API, refer to the [`examples/`](examples/) directory.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

This project is licensed under the MIT License.
