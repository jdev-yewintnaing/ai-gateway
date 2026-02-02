package anthropic

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/yewintnaing/ai-gateway/internal/providers"
)

type Provider struct {
	apiKey string
	client *http.Client
}

func NewProvider(apiKey string) *Provider {
	return &Provider{
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *Provider) Chat(req providers.ChatRequest) (*providers.ChatResponse, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY is not set")
	}

	if p.apiKey == "mock" {
		return &providers.ChatResponse{
			ID:      "mock-anthropic-id",
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   req.Model,
			Choices: []struct {
				Index        int               `json:"index"`
				Message      providers.Message `json:"message"`
				FinishReason string            `json:"finish_reason"`
			}{
				{
					Index: 0,
					Message: providers.Message{
						Role:    "assistant",
						Content: fmt.Sprintf("Anthropic Mock: %s", req.Messages[len(req.Messages)-1].Content),
					},
					FinishReason: "stop",
				},
			},
			Usage: providers.Usage{
				PromptTokens:     15,
				CompletionTokens: 25,
				TotalTokens:      40,
			},
		}, nil
	}

	return nil, fmt.Errorf("anthropic real API call not implemented in MVP (use mock)")
}

func (p *Provider) ChatStream(req providers.ChatRequest) (<-chan providers.ChatChunk, <-chan error) {
	chunkCh := make(chan providers.ChatChunk)
	errCh := make(chan error, 1)

	if p.apiKey == "mock" {
		go func() {
			defer close(chunkCh)
			defer close(errCh)
			content := fmt.Sprintf("Anthropic Mock Stream: %s", req.Messages[len(req.Messages)-1].Content)
			words := strings.Split(content, " ")
			for i, word := range words {
				chunkCh <- providers.ChatChunk{
					ID:      "mock-anth-stream-id",
					Object:  "chat.completion.chunk",
					Created: time.Now().Unix(),
					Model:   req.Model,
					Choices: []struct {
						Index int `json:"index"`
						Delta struct {
							Content string `json:"content"`
						} `json:"delta"`
						FinishReason string `json:"finish_reason"`
					}{
						{
							Index: 0,
							Delta: struct {
								Content string `json:"content"`
							}{Content: word + " "},
						},
					},
				}
				if i == len(words)-1 {
					chunkCh <- providers.ChatChunk{
						ID:      "mock-anth-stream-id",
						Object:  "chat.completion.chunk",
						Created: time.Now().Unix(),
						Model:   req.Model,
						Choices: []struct {
							Index int `json:"index"`
							Delta struct {
								Content string `json:"content"`
							} `json:"delta"`
							FinishReason string `json:"finish_reason"`
						}{
							{
								Index:        0,
								FinishReason: "stop",
							},
						},
					}
				}
				time.Sleep(50 * time.Millisecond)
			}
		}()
		return chunkCh, errCh
	}

	errCh <- fmt.Errorf("anthropic real streaming not implemented")
	return chunkCh, errCh
}
