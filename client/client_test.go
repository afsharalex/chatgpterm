package client

import "testing"

func TestClientDefaults(t *testing.T) {
	expectedModel := "text-davinci-003"
	expectedTemperature := 0
	expectedMaxTokens := 300
	c := NewClient()

	if c.Model != expectedModel {
		t.Errorf("expected model %q, got %q", expectedModel, c.Model)
	}

	if c.Temperature != expectedTemperature {
		t.Errorf("expected temperature %q, got %q", expectedTemperature, c.Temperature)
	}

	if c.MaxTokens != expectedMaxTokens {
		t.Errorf("expected max_tokens %q, got %q", expectedMaxTokens, c.MaxTokens)
	}
}

type MockClient struct {
}

func (c MockClient) Query(question string) (string, error) {
	return "Blah", nil
}

func TestClient(t *testing.T) {
	c := &MockClient{}
	query := "Explain life"

	response, err := c.Query(query)
	if err != nil {
		t.Fatalf("Error received from query %q, %v", query, err)
	}

	if response == "" {
		t.Error("expected response, got nil")
	}
}
