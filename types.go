package gofabric

// Entity is a type constraint for generic functions that operate on Pattern, Context, or Session types.
type Entity interface {
	Pattern | Context | Session
}

type EntityType string

const (
	EntityTypeContext EntityType = "context"
	EntityTypePattern EntityType = "pattern"
	EntityTypeSession EntityType = "session"
)

type StreamResponseType string

const (
	StreamResponseTypeComplete StreamResponseType = "complete"
	StreamResponseTypeContent  StreamResponseType = "content"
	StreamResponseTypeError    StreamResponseType = "error"
)

// AvailableModels contains a list of available model names and their vendors.
type AvailableModels struct {
	Models  []string            `json:"models"`  // Models is a list of model names.
	Vendors map[string][]string `json:"vendors"` // Vendors maps vendor names to their models.
}

type ChatOptions struct {
	Model              string  `json:"model"`              // Model is the name of the model to use for the chat.
	Temperature        float64 `json:"temperature"`        // Temperature controls the randomness of the model's responses.
	TopP               float64 `json:"topP"`               // TopP is the nucleus sampling parameter for controlling diversity.
	PresencePenalty    float64 `json:"presencePenalty"`    // PresencePenalty discourages repetition of tokens already present in the conversation.
	FrequencyPenalty   float64 `json:"frequencyPenalty"`   // FrequencyPenalty discourages repetition of tokens based on their frequency in the conversation.
	Raw                bool    `json:"raw"`                // Raw indicates whether to return raw model output without processing.
	Seed               int     `json:"seed"`               // Seed is used for random number generation to ensure reproducibility.
	ModelContextLength int     `json:"modelContextLength"` // ModelContextLength is the maximum context length for the model.
}

type ChatRequest struct {
	Prompts     []PromptRequest `json:"prompts"`     // Prompts is a list of prompt requests to be processed.
	Language    string          `json:"language"`    // Language specifies the language for the chat.
	ChatOptions ChatOptions     `json:"chatOptions"` // ChatOptions contains various options for the chat session.
}

// Config holds configuration values for various LLM providers.
type Config struct {
	Anthropic  string `json:"anthropic"`  // Anthropic API key.
	DeepSeek   string `json:"deepseek"`   // DeepSeek API key.
	Gemini     string `json:"gemini"`     // Gemini API key.
	Grokai     string `json:"grokai"`     // Grokai API key.
	Groq       string `json:"groq"`       // Groq API key.
	LMStudio   string `json:"lmstudio"`   // LMStudio API key.
	Mistral    string `json:"mistral"`    // Mistral API key.
	Ollama     string `json:"ollama"`     // Ollama API key.
	OpenAI     string `json:"openai"`     // OpenAI API key.
	OpenRouter string `json:"openrouter"` // OpenRouter API key.
	Silicon    string `json:"silicon"`    // Silicon API key.
}

// Context represents a named context file with its content.
type Context struct {
	Name    string `json:"name"`    // Name of the context.
	Content string `json:"content"` // Content is the text or data stored in the context.
}

// Pattern represents a reusable prompt pattern with a name, description, and pattern string.
type Pattern struct {
	Name        string `json:"name"`        // Name of the pattern.
	Description string `json:"description"` // Description provides details about the pattern's purpose.
	Pattern     string `json:"pattern"`     // Pattern is the actual prompt or template string.
}

type PromptRequest struct {
	UserInput    string `json:"userInput"`    // UserInput is the input provided by the user for the prompt.
	Vendor       string `json:"vendor"`       // Vendor is the name of the LLM vendor (e.g., OpenAI, Anthropic).
	Model        string `json:"model"`        // Model is the name of the model to use for the prompt.
	ContextName  string `json:"contextName"`  // ContextName is the name of the context to use.
	PatternName  string `json:"patternName"`  // PatternName is the name of the pattern to use.
	StrategyName string `json:"strategyName"` // StrategyName is the name of the strategy to use.
}

// Session represents a chat session with a name and a list of messages.
type Session struct {
	Name     string `json:"name"` // Name of the session.
	Messages []struct {
		Role    string `json:"role"`    // Role is the sender's role (e.g., user, assistant).
		Content string `json:"content"` // Content is the message text.
	} `json:"messages"`
}

// Strategy represents a named strategy with a description and associated pattern.
type Strategy struct {
	Name        string `json:"name"`        // Name of the strategy.
	Description string `json:"description"` // Description provides details about the strategy's purpose.
	Pattern     string `json:"pattern"`     // Pattern is the pattern associated with the strategy.
}

// StreamResponse represents the chat's streaming response
type StreamResponse struct {
	Type    string `json:"type"`    // "content", "error", "complete"
	Format  string `json:"format"`  // "markdown", "mermaid", "plain"
	Content string `json:"content"` // The actual content
}
