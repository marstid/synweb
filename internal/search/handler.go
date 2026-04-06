package search

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/marstid/synweb/internal/logger"
	"github.com/marstid/synweb/pkg/errors"
)

type Handler struct {
	client        *Client
	logger        *logger.Logger
	maxTextLength int
	maxResults    int
}

func NewHandler(client *Client, log *logger.Logger, maxTextLength int, maxResults int) *Handler {
	return &Handler{
		client:        client,
		logger:        log,
		maxTextLength: maxTextLength,
		maxResults:    maxResults,
	}
}

func (h *Handler) HandleSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := request.RequireString("query")
	if err != nil {
		h.logger.Warn("Missing query parameter", "error", err)
		return mcp.NewToolResultError("Query parameter is required"), nil
	}

	maxTextLength := request.GetFloat("max_text_length", float64(h.maxTextLength))
	maxResults := request.GetFloat("max_results", float64(h.maxResults))

	params := &SearchParams{
		Query:         query,
		MaxTextLength: int(maxTextLength),
		MaxResults:    int(maxResults),
	}

	h.logger.Info("Processing search request", "query", query)

	response, err := h.client.Search(ctx, params)
	if err != nil {
		h.logger.Error("Search failed", "error", err)
		if apiErr, ok := err.(*errors.APIError); ok {
			return mcp.NewToolResultError(apiErr.Error()), nil
		}
		return mcp.NewToolResultError("Search failed: " + err.Error()), nil
	}

	limitedResponse := h.client.LimitResults(response, params.MaxResults)
	truncatedResponse := h.client.TruncateResults(limitedResponse, params.MaxTextLength)

	resultJSON, err := json.MarshalIndent(truncatedResponse.Results, "", "  ")
	if err != nil {
		h.logger.Error("Failed to marshal results", "error", err)
		return mcp.NewToolResultError("Failed to format results"), nil
	}

	h.logger.Info("Search completed successfully", "results", len(truncatedResponse.Results))

	return mcp.NewToolResultText(string(resultJSON)), nil
}
