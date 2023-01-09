package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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
