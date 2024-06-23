package ollama

import (
	"context"
	"github.com/ollama/ollama/api"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ModelLlama3    = "llama3"
	ModelQwen2     = "qwen2"
	ModelPhi3      = "phi3"
	ModelAya       = "aya"
	ModelMistral   = "mistral"
	ModelGemma     = "gemma"
	ModelMixtral   = "mixtral"
	ModelCodeGemma = "codegemma"
)

func Box[T any](value T) *T {
	return &value
}

type Client struct {
	client *api.Client
}

type Config struct {
	Endpoint string
	Timeout  time.Duration
}

func New(cfg Config) (*Client, error) {
	u, err := url.Parse(cfg.Endpoint)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: api.NewClient(u, &http.Client{
			Timeout: cfg.Timeout,
		}),
	}, nil
}

type ChatClient struct {
	client   *Client
	model    string
	messages []api.Message
	format   string
}

func (c *Client) StartChat(ctx context.Context, model string, basePrompt []string) *ChatClient {
	messages := make([]api.Message, len(basePrompt))
	for i, prompt := range basePrompt {
		messages[i] = api.Message{
			Content: prompt,
			Role:    "system",
		}
	}

	return &ChatClient{
		client:   c,
		model:    model,
		messages: messages,
		format:   "",
	}
}

func (c *ChatClient) SendMessage(ctx context.Context, message string) (string, error) {
	wait := make(chan struct{}, 1)
	responseMessage := strings.Builder{}
	c.messages = append(c.messages, api.Message{
		Content: message,
		Role:    "user",
	})
	if err := c.client.client.Chat(ctx, &api.ChatRequest{
		Model:     c.model,
		Messages:  c.messages,
		Stream:    Box(true),
		Format:    c.format,
		KeepAlive: Box(api.Duration{Duration: 30 * time.Second}),
		Options:   nil,
	}, func(response api.ChatResponse) error {
		responseMessage.WriteString(response.Message.Content)

		if response.Done {
			close(wait)
		}

		return nil
	}); err != nil {
		return "", err
	}

	for range wait {
	}

	return responseMessage.String(), nil
}

func (c *Client) Generate(ctx context.Context, model string, prompt string) (string, error) {
	responseMessage := strings.Builder{}
	wait := make(chan struct{}, 1)
	if err := c.client.Generate(ctx, &api.GenerateRequest{
		Model:     model,
		Prompt:    prompt,
		Stream:    Box(true),
		KeepAlive: Box(api.Duration{Duration: 30 * time.Second}),
		Options:   nil,
	}, func(response api.GenerateResponse) error {
		responseMessage.WriteString(response.Response)

		if response.Done {
			close(wait)
		}

		return nil
	}); err != nil {
		return "", err
	}

	return responseMessage.String(), nil
}

func (c *Client) Embedding(ctx context.Context, model string, prompt string) ([]float64, error) {
	response, err := c.client.Embeddings(ctx, &api.EmbeddingRequest{
		Model:  model,
		Prompt: prompt,
	})
	if err != nil {
		return nil, err
	}

	return response.Embedding, nil
}
