package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors" // standard library
	"fmt"
	"net/http"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/circuitbreaker"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
	"github.com/marstid/synweb/internal/logger"
	synweberrors "github.com/marstid/synweb/pkg/errors"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	logger     *logger.Logger
}

func NewClient(baseURL, apiKey string, log *logger.Logger) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
		apiKey:  apiKey,
		logger:  log,
	}
}

func (c *Client) Search(ctx context.Context, params *SearchParams) (*SearchResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	c.logger.Debug("Searching with params", "query", params.Query, "max_text_length", params.MaxTextLength)

	requestBody, err := json.Marshal(SearchRequest{Query: params.Query})
	if err != nil {
		return nil, synweberrors.ErrInvalidResponse.With(err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/search", bytes.NewReader(requestBody))
	if err != nil {
		return nil, synweberrors.ErrNetworkError.With(err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	retryPolicy := retrypolicy.NewBuilder[any]().
		HandleErrors(synweberrors.ErrNetworkError, synweberrors.ErrServerError).
		WithMaxRetries(3).
		WithBackoff(time.Millisecond*100, time.Millisecond*400).
		Build()

	circuitBreaker := circuitbreaker.NewBuilder[any]().
		HandleErrors(synweberrors.ErrNetworkError, synweberrors.ErrServerError).
		WithFailureThreshold(5).
		WithDelay(30 * time.Second).
		WithSuccessThreshold(2).
		Build()

	var response *http.Response

	err = failsafe.With(retryPolicy, circuitBreaker).Run(func() error {
		resp, err := c.httpClient.Do(req)
		if err != nil {
			c.logger.Error("HTTP request failed", "error", err)
			return synweberrors.ErrNetworkError.With(err)
		}
		response = resp
		return nil
	})

	if err != nil {
		c.logger.Error("Search failed", "error", err)
		if errors.Is(err, synweberrors.ErrNetworkError) {
			return nil, err
		}
		return nil, synweberrors.ErrNetworkError.With(err)
	}

	defer response.Body.Close()

	if response.StatusCode >= 500 {
		c.logger.Error("Server error", "status", response.Status)
		return nil, synweberrors.ErrServerError.With(fmt.Errorf("status: %s", response.Status))
	}

	if response.StatusCode == 429 {
		return nil, synweberrors.ErrRateLimited
	}

	if response.StatusCode >= 400 {
		body := ""
		return nil, synweberrors.NewAPIError(response.StatusCode, response.Status, body)
	}

	var searchResponse SearchResponse
	if err := json.NewDecoder(response.Body).Decode(&searchResponse); err != nil {
		return nil, synweberrors.ErrInvalidResponse.With(err)
	}

	c.logger.Debug("Search completed", "results_count", len(searchResponse.Results))

	return &searchResponse, nil
}

func (c *Client) TruncateResults(resp *SearchResponse, maxLength int) *SearchResponse {
	if maxLength <= 0 {
		return resp
	}

	truncated := make([]SearchResult, len(resp.Results))
	for i, result := range resp.Results {
		truncated[i] = result
		if len(result.Text) > maxLength {
			truncated[i].Text = result.Text[:maxLength] + "..."
		}
	}

	return &SearchResponse{Results: truncated}
}
