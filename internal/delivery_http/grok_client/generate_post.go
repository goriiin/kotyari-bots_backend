package grok_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-faster/errors"
)

const targetUrl = "https://api.x.ai/v1/chat/completions"

func (c *GrokClient) GeneratePost(ctx context.Context, botPrompt, profilePrompt string) (string, error) {
	botMessage := messageFromPrompt(systemRole, botPrompt)
	userMessage := messageFromPrompt(userRole, profilePrompt)

	req := GrokRequest{
		Model:    defaultModel,
		Messages: messagesFromGrokMessage(botMessage, userMessage),
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal request")
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, targetUrl, bytes.NewBuffer(body))
	if err != nil {
		return "", errors.Wrap(err, "failed to create request")
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.ApiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", errors.Wrap(err, "failed to perform request")
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			// log err (линтер ругается)
			fmt.Println("failed to close body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("grok returned non-200 response status: %v", resp.StatusCode)
	}

	var grokResp GrokResponse
	if err := json.NewDecoder(resp.Body).Decode(&grokResp); err != nil {
		return "", errors.Errorf("failed to unmarshal.\n Returned body: %v", resp.Body)
	}

	return grokResp.Choices[0].Message.Content, nil
}
