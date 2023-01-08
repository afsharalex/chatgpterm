package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiClient interface {
	Query(string) (string, error)
}

type Client struct {
	Model       string `json:"model"`
	MaxTokens   int    `json:"max_tokens"`
	Temperature int    `json:"temperature"`
}

func NewClient() *Client {
	return &Client{
		Model:       "text-davinci-003",
		MaxTokens:   300,
		Temperature: 0,
	}
}

func (c Client) Query(question string) (string, error) {
	jsonBody, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodGet, "https://api.openai.com/v1/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	token := "1234"
	bearerToken := fmt.Sprintf("Bearer %s", token)

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Authorization", bearerToken)

	return "", nil
}
