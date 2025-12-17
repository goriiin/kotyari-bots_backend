package otvet

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-faster/errors"
)

const (
	createPostEndpoint        = "/api/topic/question"
	predictTagsSpacesEndpoint = "/api/tags_spaces/predict"
)

// CreatePost creates a new post/question on otvet.mail.ru
func (c *OtvetClient) CreatePost(ctx context.Context, req *CreatePostRequest) (*CreatePostResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request")
	}

	url := c.baseURL + createPostEndpoint
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")

	// Set authentication cookie
	httpReq.AddCookie(&http.Cookie{
		Name:  "Auth-Token",
		Value: c.config.AuthToken,
	})

	// Perform request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request")
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log error if needed
			_ = closeErr
		}
	}()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	// Check status code
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.Errorf("otvet.mail.ru returned non-2xx response status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var createResp CreatePostResponse
	if err := json.Unmarshal(respBody, &createResp); err != nil {
		// If response is not JSON or doesn't match expected structure, return raw response info
		return nil, errors.Wrapf(err, "failed to unmarshal response, status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return &createResp, nil
}

// TextContentNode represents a simple text content node
type TextContentNode struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// CreatePostSimple is a helper function to create a simple post with minimal required fields
func (c *OtvetClient) CreatePostSimple(ctx context.Context, title string, contentText string, topicType int, spaces []Space) (*CreatePostResponse, error) {
	req := &CreatePostRequest{
		Title:     title,
		TopicType: topicType,
		Version:   0,
		VisibleTo: 0,
		Content: &Content{
			Type: "doc",
			Content: []ContentNode{
				{
					Type: "paragraph",
					Content: []interface{}{
						TextContentNode{
							Type: "text",
							Text: contentText,
						},
					},
				},
			},
		},
		Tags:   []Tag{},
		Spaces: spaces,
	}

	return c.CreatePost(ctx, req)
}

// NewTextContent creates a simple text content structure
func NewTextContent(text string) *Content {
	return &Content{
		Type: "doc",
		Content: []ContentNode{
			{
				Type: "paragraph",
				Content: []interface{}{
					TextContentNode{
						Type: "text",
						Text: text,
					},
				},
			},
		},
	}
}

// PredictTagsSpaces predicts tags and spaces for given text
func (c *OtvetClient) PredictTagsSpaces(ctx context.Context, text string) (*PredictTagsSpacesResponse, error) {
	req := PredictTagsSpacesRequest{
		Data: []PredictDataItem{
			{
				ID:   "0",
				Text: text,
			},
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request")
	}

	url := c.baseURL + predictTagsSpacesEndpoint
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")

	// Set authentication cookie
	httpReq.AddCookie(&http.Cookie{
		Name:  "Auth-Token",
		Value: c.config.AuthToken,
	})

	// Perform request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request")
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log error if needed
			_ = closeErr
		}
	}()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	// Check status code
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.Errorf("otvet.mail.ru returned non-2xx response status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var predictResp PredictTagsSpacesResponse
	if err := json.Unmarshal(respBody, &predictResp); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal response, status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return &predictResp, nil
}
