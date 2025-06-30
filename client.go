package gofabric

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tmaxmax/go-sse"
)

const (
	apiKeyHeaderName         = "X-API-Key"
	defaultHTTPClientTimeout = 60 * time.Second
)

// Client represents a client for the Fabric API server.
type Client struct {
	// The URL of the Fabric API server
	host string
	// The API key for authentication
	apiKey string
	// The HTTP client for making requests
	httpClient *http.Client
}

// Option represents a function that configures the Client using the functional options pattern.
type Option func(*Client)

// NewClient creates a new Client instance with the specified host and options.
// The Client will use the default HTTP client with a timeout of 60 seconds.
//
// To set the API key, use the WithAPIKey option.
// To customize the HTTP client, use the WithHTTPClient option.
func NewClient(host string, opts ...Option) *Client {
	client := &Client{
		host: host,
		httpClient: &http.Client{
			Timeout: defaultHTTPClientTimeout,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// WithAPIKey sets the API key for the client.
func WithAPIKey(apiKey string) Option {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}

// WithHTTPClient sets the HTTP client for the client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

func (c *Client) doRequest(
	ctx context.Context,
	method string,
	path string,
	body io.Reader,
) (*http.Response, error) {
	parsedURL, err := url.Parse(c.host)
	if err != nil {
		return nil, fmt.Errorf("invalid host %q: %w", c.host, err)
	}

	url := parsedURL.JoinPath(path).String()

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s %s: %w", method, url, err)
	}

	if c.apiKey != "" {
		req.Header.Set(apiKeyHeaderName, c.apiKey)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %s %s: %w", method, url, err)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		var body *string
		if err == nil && len(bodyBytes) > 0 {
			bodyStr := string(bodyBytes)
			body = &bodyStr
		}

		return nil, &HTTPError{
			URL:        url,
			StatusCode: resp.StatusCode,
			Body:       body,
		}
	}

	return resp, nil
}

func createEntity(
	client *Client,
	ctx context.Context,
	entityType EntityType,
	entityName string,
	body io.Reader,
) error {
	resp, err := client.doRequest(ctx, http.MethodPost, "/"+string(entityType)+"s/"+entityName, body)
	if err != nil {
		return fmt.Errorf("failed to create %s `%s`: %w", entityType, entityName, err)
	}
	defer func() { _ = resp.Body.Close() }()

	return nil
}

func deleteEntity(client *Client, ctx context.Context, entityType EntityType, entityName string) error {
	resp, err := client.doRequest(ctx, http.MethodDelete, "/"+string(entityType)+"s/"+entityName, nil)
	if err != nil {
		return fmt.Errorf("failed to delete %s `%s`: %w", entityType, entityName, err)
	}
	defer func() { _ = resp.Body.Close() }()

	return nil
}

func entityExists(
	client *Client,
	ctx context.Context,
	entityType EntityType,
	entityName string,
) (bool, error) {
	resp, err := client.doRequest(ctx, http.MethodGet, "/"+string(entityType)+"s/exists/"+entityName, nil)
	if err != nil {
		return false, fmt.Errorf(
			"failed to check if %s `%s` exists: %w",
			entityType,
			entityName,
			err,
		)
	}
	defer func() { _ = resp.Body.Close() }()

	var exists bool
	if err := json.NewDecoder(resp.Body).Decode(&exists); err != nil {
		return false, fmt.Errorf(
			"failed to decode %s existence check for `%s`: %w",
			entityType,
			entityName,
			err,
		)
	}

	return exists, nil
}

func getEntity[T Entity](
	client *Client,
	ctx context.Context,
	entityType EntityType,
	entityName string,
) (*T, error) {
	resp, err := client.doRequest(ctx, http.MethodGet, "/"+string(entityType)+"s/"+entityName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s `%s`: %w", entityType, entityName, err)
	}
	defer func() { _ = resp.Body.Close() }()

	var entity T
	if err := json.NewDecoder(resp.Body).Decode(&entity); err != nil {
		return nil, fmt.Errorf("failed to decode %s `%s`: %w", entityType, entityName, err)
	}

	return &entity, nil
}

func listEntity(client *Client, ctx context.Context, entityType string) ([]string, error) {
	resp, err := client.doRequest(ctx, http.MethodGet, "/"+entityType+"/names", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list %s: %w", entityType, err)
	}
	defer func() { _ = resp.Body.Close() }()

	var entitys []string
	if err := json.NewDecoder(resp.Body).Decode(&entitys); err != nil {
		return nil, fmt.Errorf("failed to decode %s: %w", entityType, err)
	}

	return entitys, nil
}

func renameEntity(
	client *Client,
	ctx context.Context,
	entityType EntityType,
	oldEntityName string,
	newEntityName string,
) error {
	resp, err := client.doRequest(
		ctx,
		http.MethodPut,
		"/"+string(entityType)+"s/rename/"+oldEntityName+"/"+newEntityName,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to rename %s `%s` to `%s`: %w",
			entityType,
			oldEntityName,
			newEntityName,
			err,
		)
	}
	defer func() { _ = resp.Body.Close() }()

	return nil
}

// Chat initiates a chat session with the specified chat request.
func (c *Client) Chat(ctx context.Context, chatRequest *ChatRequest) (<-chan StreamResponse, error) {
	data, err := json.Marshal(chatRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to encode chat request: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/chat", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to initiate chat: %w", err)
	}

	streamResponseChannel := make(chan StreamResponse)

	// Start goroutine to read Server-sent events
	go func() {
		defer func() { _ = resp.Body.Close() }()
		defer close(streamResponseChannel)

		for event, err := range sse.Read(resp.Body, nil) {
			if err != nil {
				streamResponseChannel <- StreamResponse{
					Type:    string(StreamResponseTypeError),
					Format:  "plain",
					Content: fmt.Errorf("failed to read SSE response: %w", err).Error(),
				}

				break // Exit the goroutine on reading error
			}

			var streamResponse StreamResponse
			if err := json.Unmarshal([]byte(event.Data), &streamResponse); err != nil {
				streamResponseChannel <- StreamResponse{
					Type:    string(StreamResponseTypeError),
					Format:  "plain",
					Content: fmt.Errorf("failed to parse SSE response: %w", err).Error(),
				}

				return // Exit the goroutine on parsing error
			}

			streamResponseChannel <- streamResponse

			if streamResponse.Type == string(StreamResponseTypeComplete) {
				return // Exit the goroutine when the response is complete
			}
		}
	}()

	return streamResponseChannel, nil
}

// CreateContext creates a new context.
func (c *Client) CreateContext(ctx context.Context, name string, body io.Reader) error {
	return createEntity(c, ctx, EntityTypeContext, name, body)
}

// CreatePattern creates a new pattern.
func (c *Client) CreatePattern(ctx context.Context, name string, body io.Reader) error {
	return createEntity(c, ctx, EntityTypePattern, name, body)
}

// CreateSession creates a new session.
func (c *Client) CreateSession(ctx context.Context, name string, body io.Reader) error {
	return createEntity(c, ctx, EntityTypeSession, name, body)
}

// DeleteContext deletes a context.
func (c *Client) DeleteContext(ctx context.Context, name string) error {
	return deleteEntity(c, ctx, EntityTypeContext, name)
}

// DeletePattern deletes a pattern.
func (c *Client) DeletePattern(ctx context.Context, name string) error {
	return deleteEntity(c, ctx, EntityTypePattern, name)
}

// DeleteSession deletes a session.
func (c *Client) DeleteSession(ctx context.Context, name string) error {
	return deleteEntity(c, ctx, EntityTypeSession, name)
}

// GetConfig retrieves the configuration of fabric.
func (c *Client) GetConfig(ctx context.Context) (*Config, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/config", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var config Config
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return &config, nil
}

// GetContextMetadata retrieves the metadata of a context.
func (c *Client) GetContextMetadata(ctx context.Context, name string) (*Context, error) {
	return getEntity[Context](c, ctx, EntityTypeContext, name)
}

// GetPatternMetadata retrieves the metadata of a pattern.
func (c *Client) GetPatternMetadata(ctx context.Context, name string) (*Pattern, error) {
	return getEntity[Pattern](c, ctx, EntityTypePattern, name)
}

// GetSessionMetadata retrieves the metadata of a session.
func (c *Client) GetSessionMetadata(ctx context.Context, name string) (*Session, error) {
	return getEntity[Session](c, ctx, EntityTypeSession, name)
}

// ListContexts retrieves the list of contexts.
func (c *Client) ListContexts(ctx context.Context) ([]string, error) {
	return listEntity(c, ctx, "contexts")
}

// ListNames retrieves a list of models.
func (c *Client) ListModels(ctx context.Context) (*AvailableModels, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/models/names", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get models: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var availableModels AvailableModels
	if err := json.NewDecoder(resp.Body).Decode(&availableModels); err != nil {
		return nil, fmt.Errorf("failed to decode models: %w", err)
	}

	return &availableModels, nil
}

// ListPatterns retrieves the list of patterns.
func (c *Client) ListPatterns(ctx context.Context) ([]string, error) {
	return listEntity(c, ctx, "patterns")
}

// ListSessions retrieves the list of sessions.
func (c *Client) ListSessions(ctx context.Context) ([]string, error) {
	return listEntity(c, ctx, "sessions")
}

// ListStrategies retrieves a list of strategies.
func (c *Client) ListStrategies(ctx context.Context) ([]Strategy, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/strategies", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get strategies: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var strategies []Strategy
	if err := json.NewDecoder(resp.Body).Decode(&strategies); err != nil {
		return nil, fmt.Errorf("failed to decode strategies: %w", err)
	}

	return strategies, nil
}

// ContextExists checks if a context exists.
func (c *Client) ContextExists(ctx context.Context, name string) (bool, error) {
	return entityExists(c, ctx, EntityTypeContext, name)
}

// PatternExists checks if a pattern exists.
func (c *Client) PatternExists(ctx context.Context, name string) (bool, error) {
	return entityExists(c, ctx, EntityTypePattern, name)
}

// SessionExists checks if a session exists.
func (c *Client) SessionExists(ctx context.Context, name string) (bool, error) {
	return entityExists(c, ctx, EntityTypeSession, name)
}

// RenameContext renames a context.
func (c *Client) RenameContext(ctx context.Context, oldName string, newName string) error {
	return renameEntity(c, ctx, EntityTypeContext, oldName, newName)
}

// RenamePattern renames a pattern.
func (c *Client) RenamePattern(ctx context.Context, oldName string, newName string) error {
	return renameEntity(c, ctx, EntityTypePattern, oldName, newName)
}

// RenameSession renames a session.
func (c *Client) RenameSession(ctx context.Context, oldName string, newName string) error {
	return renameEntity(c, ctx, EntityTypeSession, oldName, newName)
}

// UpdateConfig updates the configuration of fabric.
func (c *Client) UpdateConfig(ctx context.Context, config *Config) error {
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPut, "/config/update", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	return nil
}
