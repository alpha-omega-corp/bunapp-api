package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type GptClient struct {
	client *http.Client
	token  string
	host   string
}

type CompletionMessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionRequest struct {
	Model    string                     `json:"model"`
	Messages []CompletionMessageRequest `json:"messages"`
}

func (app *App) initClient() {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	app.gptClient = &GptClient{
		client: &http.Client{Transport: tr},
		token:  app.config.GPT.BEARER,
		host:   app.config.GPT.HOST,
	}
}

func (g *GptClient) Request(path string) (*http.Response, error) {
	CompletionMessageRequest := []CompletionMessageRequest{
		{
			Role:    "Chad",
			Content: "I want to lose weight",
		},
		{
			Role:    "Chad",
			Content: "I want to lose weight",
		},
	}

	unbufferedBody := &CompletionRequest{
		Model:    "gpt-3.5",
		Messages: CompletionMessageRequest,
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(unbufferedBody)

	req, err := http.NewRequest("POST", g.host+path, &body)
	if err != nil {
		return nil, err
	}

	res, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, err
}
