package types

type CompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionContext struct {
	Model    string              `json:"model"`
	Messages []CompletionMessage `json:"messages"`
}

type CompletionResponse struct {
	Id      string                     `json:"id"`
	Object  string                     `json:"object"`
	Created int64                      `json:"created"`
	Model   string                     `json:"model"`
	Choices []CompletionChoiceResponse `json:"choices"`
}

type CompletionChoiceResponse struct {
	Index      int               `json:"index"`
	Message    CompletionMessage `json:"message"`
	ExitReason string            `json:"finish_reason"`
}
