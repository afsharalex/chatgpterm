package client

type ApiClient interface {
	Query(string) (string, error)
}

type ApiRequest struct {
	Model       string `json:"model"`
	MaxTokens   int    `json:"max_tokens"`
	Temperature int    `json:"temperature"`
	Prompt      string `json:"prompt"`
}

type ApiChoice struct {
	Text         string      `json:"text"`
	Index        int         `json:"index"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

type ApiUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ApiError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

type ApiResponse struct {
	Id      string      `json:"id"`
	Object  string      `json:"object"`
	Created int         `json:"created"`
	Model   string      `json:"model"`
	Choices []ApiChoice `json:"choices"`
	Usage   ApiUsage    `json:"usage"`
	Error   ApiError    `json:"error"`
}
