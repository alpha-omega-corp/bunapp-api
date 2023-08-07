package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alpha-omega-corp/bunapp-api/types"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type ClientOptions struct {
	host  string
	token string
}

type GptClient struct {
	client   *http.Client
	cContext *types.CompletionContext
	options  *ClientOptions
}

type RoundTripper struct {
	next   http.RoundTripper
	logger io.Writer
}

func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	_, err := fmt.Fprintf(rt.logger, "[%s] %s %s\n",
		time.Now().Format(time.ANSIC), req.Method, req.URL.String())

	if err != nil {
		return nil, err
	}

	return rt.next.RoundTrip(req)
}

func (app *App) initClient() {
	c := &http.Client{
		Transport: &RoundTripper{
			next:   http.DefaultTransport,
			logger: os.Stdout,
		},
		Timeout: 100 * time.Second,
	}

	app.gptClient = &GptClient{
		client: c,
		cContext: &types.CompletionContext{
			Model:    "gpt-3.5-turbo",
			Messages: []types.CompletionMessage{},
		},
		options: &ClientOptions{
			token: app.config.GPT.TOKEN,
			host:  app.config.GPT.HOST,
		},
	}
}

func (gpt *GptClient) UserRequest(prompt string) ([]types.CompletionMessage, error) {
	message := types.CompletionMessage{
		Role:    "user",
		Content: prompt,
	}

	return gpt.request(message)
}

func (gpt *GptClient) request(m types.CompletionMessage) ([]types.CompletionMessage, error) {
	gpt.cContext.Messages = append(gpt.cContext.Messages, m)

	req, err := http.NewRequest(http.MethodPost, gpt.reqUrl(), gpt.completionBody())
	if err != nil {
		return nil, err
	}

	gpt.setToken(req)

	res, err := gpt.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	var completionRes types.CompletionResponse
	cErr := json.Unmarshal(body, &completionRes)
	if cErr != nil {
		return nil, cErr
	}

	gpt.cContext.Messages = append(gpt.cContext.Messages, completionRes.Choices[0].Message)
	return gpt.cContext.Messages, nil
}

func (gpt *GptClient) completionBody() *bytes.Reader {
	marshalled, err := json.Marshal(gpt.cContext)
	if err != nil {
		panic(err)
	}

	return bytes.NewReader(marshalled)
}

func (gpt *GptClient) setToken(req *http.Request) {
	bearer := []string{"Bearer", gpt.options.token}
	req.Header.Set("Authorization", strings.Join(bearer, " "))
	req.Header.Set("Content-Type", "application/json")
}

func (gpt *GptClient) reqUrl() string {
	return gpt.options.host + "/chat/completions"
}
