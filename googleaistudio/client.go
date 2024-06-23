package googleaistudio

import (
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"io"
	"net/http"
	"strings"
)

const (
	ModelGemini15Flash  = "gemini-1.5-flash"
	ModelGemini15Pro    = "gemini-1.5-pro"
	ModelGemini10Pro    = "gemini-1.0-pro"
	ModelGemini10Vision = "gemini-1.0-pro-vision"

	ModelTextEmbeddingGecko001 = "textembedding-gecko@001"
	ModelTextEmbeddingGecko002 = "textembedding-gecko@002"
	ModelTextEmbeddingGecko003 = "textembedding-gecko@003"
	ModelTextEmbedding004      = "text-embedding-004"

	ModelTextEmbeddingGeckoMultilingual001 = "textembedding-gecko-multilingual@001"
	ModelTextMultilingualEmbedding001      = "text-multilingual-embedding-002"
	ModelEmbeddingsForMultilingual         = "multimodalembedding"

	ModelImgGen2                = "imagegeneration@006"
	ModelCodeyForCodeCompletion = "code-gecko"

	ModelMedLMMedium = "medlm-medium"
	ModelMedLMLarge  = "medlm-large"
)

type Client struct {
	client *genai.Client
}

type Config struct {
	ApiKey string
}

func New(ctx context.Context, cfg Config) (*Client, error) {
	c, err := genai.NewClient(ctx, option.WithAPIKey(cfg.ApiKey))
	if err != nil {
		return nil, err
	}

	context.AfterFunc(ctx, func() {
		c.Close()
	})

	return &Client{client: c}, nil
}

func mergeResponse(resp *genai.GenerateContentResponse) string {
	if resp == nil {
		return ""
	}

	builder := strings.Builder{}
	for _, candidate := range resp.Candidates {
		if candidate.Content != nil {
			for _, part := range candidate.Content.Parts {
				builder.WriteString(fmt.Sprint(part))
			}
		}
	}

	return builder.String()
}

func (c *Client) Generate(ctx context.Context, model, prompt string) (string, error) {
	resp, err := c.client.GenerativeModel(model).GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	return mergeResponse(resp), nil
}

type ChatClient struct {
	session *genai.ChatSession
}

func (c *Client) StartChat(ctx context.Context, model string, basePrompt ...string) (*ChatClient, error) {
	sess := c.client.GenerativeModel(model).StartChat()
	pars := make([]genai.Part, len(basePrompt))
	for i, prompt := range basePrompt {
		pars[i] = genai.Text(prompt)
	}

	if _, err := sess.SendMessage(ctx, pars...); err != nil {
		return nil, err
	}

	return &ChatClient{session: sess}, nil
}

func Text(text string) genai.Part {
	return genai.Text(text)
}

func Image(data []byte) genai.Part {
	return genai.ImageData(http.DetectContentType(data), data)
}

func File(ctx context.Context, c *Client, name string, reader io.Reader) (genai.Part, error) {
	f, err := c.client.UploadFile(ctx, name, reader, nil)
	if err != nil {
		return nil, err
	}

	context.AfterFunc(ctx, func() {
		c.client.DeleteFile(ctx, f.URI)
	})

	fd := genai.FileData{
		URI:      f.URI,
		MIMEType: f.MIMEType,
	}

	return fd, nil
}

func FileWithoutAutoRemove(ctx context.Context, c *Client, name string, reader io.Reader) (genai.Part, error) {
	f, err := c.client.UploadFile(ctx, name, reader, nil)
	if err != nil {
		return nil, err
	}

	fd := genai.FileData{
		URI:      f.URI,
		MIMEType: f.MIMEType,
	}

	return fd, nil
}

func (c *Client) Embedding(ctx context.Context, model string, data ...genai.Part) ([]float32, error) {
	resp, err := c.client.EmbeddingModel(model).EmbedContent(ctx, data...)
	if err != nil {
		return nil, err
	}

	return resp.Embedding.Values, nil
}
