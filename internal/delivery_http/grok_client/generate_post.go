package grok_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/pkg/grok"
)

const (
	defaultRetriesNum = 3
	defaultRetryDelay = 2 * time.Second
)

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

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, grok.GrokTargetUrl, bytes.NewBuffer(body))
	if err != nil {
		return "", errors.Wrap(err, "failed to create request")
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.ApiKey)

	var resp *http.Response
	for i := 0; i < defaultRetriesNum; i++ {
		resp, err = c.httpClient.Do(httpReq)
		if err == nil {
			break
		}

		backoff := time.Duration(math.Pow(2, float64(i))) * time.Second
		// log failed attempt
		fmt.Printf("Request failed: %v. Retrying in %v... (attempt %d/%d)", err, backoff, i+1, defaultRetriesNum)
		time.Sleep(defaultRetryDelay)
	}

	if err != nil {
		return "", errors.Wrapf(err, "failed to perform request after %d attempts", defaultRetriesNum)
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
