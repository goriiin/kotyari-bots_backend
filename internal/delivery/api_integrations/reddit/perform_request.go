package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: move
const redditAPIString = "reddit"

func (r *RedditAPIDelivery) performRequests() ([]RedditAPIResponse, error) {
	// TODO: log

	ctx := context.Background()

	integrations, err := r.integration.GetIntegrations(ctx, redditAPIString)
	if err != nil {
		// TODO: log, err
		return nil, err
	}

	var redditAPIResponses []RedditAPIResponse

	// TODO: async requests
	for _, integration := range integrations {
		req, err := http.NewRequest(http.MethodGet, integration.Url, http.NoBody)
		if err != nil {
			// TODO: log, err

			return nil, err
		}

		resp, err := r.client.Do(req)
		if err != nil {
			// TODO: log, err

			return nil, err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			// TODO: log, err

			return nil, fmt.Errorf("респонс бади хуйня из жопы %s\n%s", err.Error(), string(body))
		}

		err = resp.Body.Close()
		if err != nil {
			// TODO: log, err

			return nil, err
		}

		var redditAPIResponse RedditAPIResponse
		if err := json.Unmarshal(body, &redditAPIResponse); err != nil {
			// TODO: log, err

			return nil, fmt.Errorf("анмаршал хуйни УРЛ: %s", integration.Url)
		}
		redditAPIResponses = append(redditAPIResponses, redditAPIResponse)
	}

	return redditAPIResponses, nil
}
