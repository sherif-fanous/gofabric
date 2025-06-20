package gofabric_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sherif-fanous/gofabric"
)

func TestChat(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/chat" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	chatRequest := &gofabric.ChatRequest{
		Prompts: []gofabric.PromptRequest{
			{
				UserInput:   "A Hacker News clone",
				Vendor:      "Gemini",
				Model:       "gemini-2.5-flash-preview-05-20",
				PatternName: "create_prd",
			},
		},
		Language:    "en",
		ChatOptions: gofabric.ChatOptions{},
	}

	client := gofabric.NewClient(ts.URL)
	_, err := client.Chat(context.Background(), chatRequest)
	if err != nil {
		t.Fatalf("Failed to chat: %v", err)
	}
}

func TestCreateContext(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/contexts/test" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	err := client.CreateContext(context.Background(), "test", nil)
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}
}

func TestCreatePattern(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/patterns/test" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	err := client.CreatePattern(context.Background(), "test", nil)
	if err != nil {
		t.Fatalf("Failed to create pattern: %v", err)
	}
}

func TestCreateSession(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/sessions/test" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	err := client.CreateSession(context.Background(), "test", nil)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}
}

func TestDeleteContext(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/contexts/test" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	err := client.DeleteContext(context.Background(), "test")
	if err != nil {
		t.Fatalf("Failed to delete context: %v", err)
	}
}

func TestDeletePattern(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/patterns/test" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	err := client.DeletePattern(context.Background(), "test")
	if err != nil {
		t.Fatalf("Failed to delete pattern: %v", err)
	}
}

func TestDeleteSession(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/sessions/test" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	err := client.DeleteSession(context.Background(), "test")
	if err != nil {
		t.Fatalf("Failed to delete session: %v", err)
	}
}

func TestGetConfig(t *testing.T) {
	t.Parallel()

	want := &gofabric.Config{
		Anthropic:  "anthropic_key",
		DeepSeek:   "deepseek_key",
		Gemini:     "gemini_key",
		Grokai:     "grokai_key",
		Groq:       "groq_key",
		LMStudio:   "lmstudio_key",
		Mistral:    "mistral_key",
		Ollama:     "ollama_key",
		OpenAI:     "openai_key",
		OpenRouter: "openrouter_key",
		Silicon:    "silicon_key",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/config" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(want)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	config, err := client.GetConfig(context.Background())
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if diff := cmp.Diff(want, config); diff != "" {
		t.Fatalf("Mismatch (-want +got):\n%s", diff)
	}
}

func TestGetContextMetadata(t *testing.T) {
	t.Parallel()

	want := &gofabric.Context{Name: "test", Content: "content"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/contexts/test" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		_ = json.NewEncoder(w).Encode(want)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	context, err := client.GetContextMetadata(context.Background(), "test")
	if err != nil {
		t.Fatalf("Failed to get contexts: %v", err)
	}

	if diff := cmp.Diff(want, context); diff != "" {
		t.Fatalf("Mismatch (-want +got):\n%s", diff)
	}
}

func TestGetPatternMetadata(t *testing.T) {
	t.Parallel()

	want := &gofabric.Pattern{Name: "test", Description: "description", Pattern: "pattern"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/patterns/test" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		_ = json.NewEncoder(w).Encode(want)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	pattern, err := client.GetPatternMetadata(context.Background(), "test")
	if err != nil {
		t.Fatalf("Failed to get patterns: %v", err)
	}

	if diff := cmp.Diff(want, pattern); diff != "" {
		t.Fatalf("Mismatch (-want +got):\n%s", diff)
	}
}

func TestGetSession(t *testing.T) {
	t.Parallel()

	want := &gofabric.Session{
		Name: "test",
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{Role: "role", Content: "content"},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/sessions/test" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(want)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	session, err := client.GetSessionMetadata(context.Background(), "test")
	if err != nil {
		t.Fatalf("Failed to get sessions: %v", err)
	}

	if diff := cmp.Diff(want, session); diff != "" {
		t.Fatalf("Mismatch (-want +got):\n%s", diff)
	}
}

func TestListContexts(t *testing.T) {
	t.Parallel()

	want := []string{"test_1", "test_2"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/contexts/names" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		_ = json.NewEncoder(w).Encode(want)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	contexts, err := client.ListContexts(context.Background())
	if err != nil {
		t.Fatalf("Failed to get contexts: %v", err)
	}

	if diff := cmp.Diff(want, contexts); diff != "" {
		t.Fatalf("Mismatch (-want +got):\n%s", diff)
	}
}

func TestListModels(t *testing.T) {
	t.Parallel()

	want := &gofabric.AvailableModels{
		Models:  []string{"test_1", "test_2"},
		Vendors: map[string][]string{"vendor_1": {"test_1"}},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/models/names" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		_ = json.NewEncoder(w).Encode(want)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	models, err := client.ListModels(context.Background())
	if err != nil {
		t.Fatalf("Failed to get models: %v", err)
	}

	if diff := cmp.Diff(want, models); diff != "" {
		t.Fatalf("Mismatch (-want +got):\n%s", diff)
	}
}

func TestListPatterns(t *testing.T) {
	t.Parallel()

	want := []string{"test_1", "test_2"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/patterns/names" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		_ = json.NewEncoder(w).Encode(want)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	patterns, err := client.ListPatterns(context.Background())
	if err != nil {
		t.Fatalf("Failed to get patterns: %v", err)
	}

	if diff := cmp.Diff(want, patterns); diff != "" {
		t.Fatalf("Mismatch (-want +got):\n%s", diff)
	}
}

func TestListSessions(t *testing.T) {
	t.Parallel()

	want := []string{"test_1", "test_2"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/sessions/names" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		_ = json.NewEncoder(w).Encode(want)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	sessions, err := client.ListSessions(context.Background())
	if err != nil {
		t.Fatalf("Failed to get sessions: %v", err)
	}

	if diff := cmp.Diff(want, sessions); diff != "" {
		t.Fatalf("Mismatch (-want +got):\n%s", diff)
	}
}

func TestListStrategies(t *testing.T) {
	t.Parallel()

	want := []gofabric.Strategy{{Name: "test", Description: "description", Pattern: "pattern"}}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/strategies" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		_ = json.NewEncoder(w).Encode(want)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	strategies, err := client.ListStrategies(context.Background())
	if err != nil {
		t.Fatalf("Failed to get strategies: %v", err)
	}

	if diff := cmp.Diff(want, strategies); diff != "" {
		t.Fatalf("Mismatch (-want +got):\n%s", diff)
	}
}

func TestContextExists(t *testing.T) {
	t.Parallel()

	want := true

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/contexts/exists/test" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		_ = json.NewEncoder(w).Encode(want)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	exists, err := client.ContextExists(context.Background(), "test")
	if err != nil {
		t.Fatalf("Failed to check if context exists: %v", err)
	}

	if diff := cmp.Diff(want, exists); diff != "" {
		t.Fatalf("Mismatch (-want +got):\n%s", diff)
	}
}

func TestPatternExists(t *testing.T) {
	t.Parallel()

	want := true

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/patterns/exists/test" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		_ = json.NewEncoder(w).Encode(want)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	exists, err := client.PatternExists(context.Background(), "test")
	if err != nil {
		t.Fatalf("Failed to check if pattern exists: %v", err)
	}

	if diff := cmp.Diff(want, exists); diff != "" {
		t.Fatalf("Mismatch (-want +got):\n%s", diff)
	}
}

func TestSessionExists(t *testing.T) {
	t.Parallel()

	want := true

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/sessions/exists/test" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		_ = json.NewEncoder(w).Encode(want)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	exists, err := client.SessionExists(context.Background(), "test")
	if err != nil {
		t.Fatalf("Failed to check if session exists: %v", err)
	}

	if diff := cmp.Diff(want, exists); diff != "" {
		t.Fatalf("Mismatch (-want +got):\n%s", diff)
	}
}

func TestRenameContext(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/contexts/rename/test/test-2" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	err := client.RenameContext(context.Background(), "test", "test-2")
	if err != nil {
		t.Fatalf("Failed to rename context: %v", err)
	}
}

func TestRenamePattern(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/patterns/rename/test/test-2" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	err := client.RenamePattern(context.Background(), "test", "test-2")
	if err != nil {
		t.Fatalf("Failed to rename pattern: %v", err)
	}
}

func TestRenameSession(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/sessions/rename/test/test-2" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	err := client.RenameSession(context.Background(), "test", "test-2")
	if err != nil {
		t.Fatalf("Failed to rename session: %v", err)
	}
}

func TestUpdateConfig(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/config/update" {
			t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := gofabric.NewClient(ts.URL)
	err := client.UpdateConfig(context.Background(), &gofabric.Config{
		Anthropic:  "new_anthropic_key",
		DeepSeek:   "new_deepseek_key",
		Gemini:     "new_gemini_key",
		Grokai:     "new_grokai_key",
		Groq:       "new_groq_key",
		LMStudio:   "new_lmstudio_key",
		Mistral:    "new_mistral_key",
		Ollama:     "new_ollama_key",
		OpenAI:     "new_openai_key",
		OpenRouter: "new_openrouter_key",
		Silicon:    "new_silicon_key",
	})
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}
}
