package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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

type Client struct {
	Model       string
	MaxTokens   int
	Temperature int
	key         string
}

func NewClient(apiKey string) *Client {
	return &Client{
		Model:       "text-davinci-003",
		MaxTokens:   300,
		Temperature: 0,
		key:         apiKey,
	}
}

func (c Client) Query(question string) (string, error) {
	apiRequest := ApiRequest{
		Model:       c.Model,
		MaxTokens:   c.MaxTokens,
		Temperature: c.Temperature,
		Prompt:      question,
	}

	jsonBody, err := json.Marshal(apiRequest)
	if err != nil {
		log.Printf("Received error: %s", err)
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Received error: %s", err)
		return "", err
	}

	bearerToken := fmt.Sprintf("Bearer %s", c.key)

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Authorization", bearerToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Received error: %s", err)
		return "", err
	}

	return c.parseResponse(res)
}

func (c Client) parseResponse(res *http.Response) (string, error) {
	apiResponse := &ApiResponse{}
	err := json.NewDecoder(res.Body).Decode(apiResponse)
	defer res.Body.Close()
	if err != nil {
		log.Printf("Received error: %+v", err)
		return "", err
	}

	// TODO: Handle errors and display them to user. I think this should
	// be some kind of popup, or a fatal log since the user may need to
	// change something with their account...

	// fmt.Printf("API Response: %+v", res)
	// fmt.Printf("API Response Decoded: %+v", apiResponse)

	if len(apiResponse.Choices) > 0 {
		return apiResponse.Choices[0].Text, nil
	}

	return "", nil
}
