package server

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/marstid/synweb/internal/config"
	"github.com/marstid/synweb/internal/logger"
	"github.com/marstid/synweb/internal/search"
)

type Server struct {
	mcpServer *server.MCPServer
	logger    *logger.Logger
	cfg       *config.Config
}

func New(name string, log *logger.Logger, cfg *config.Config, handler *search.Handler) *Server {
	s := server.NewMCPServer(
		name,
		"", // version - not used
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	searchTool := mcp.NewTool(
		"search_web",
		mcp.WithDescription("Search the web using Synthetic API"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query string"),
		),
		mcp.WithNumber("max_text_length",
			mcp.DefaultNumber(float64(cfg.MaxTextLength)),
			mcp.Description(fmt.Sprintf("Maximum number of characters to include in the 'text' field of each result. Defaults to %d. Set higher if more detail is needed, or fetch the source URL directly for the full content.", cfg.MaxTextLength)),
		),
		mcp.WithNumber("max_results",
			mcp.DefaultNumber(float64(cfg.MaxResults)),
			mcp.Description(fmt.Sprintf("Maximum number of results to return. Defaults to %d.", cfg.MaxResults)),
		),
	)

	s.AddTool(searchTool, handler.HandleSearch)

	return &Server{
		mcpServer: s,
		logger:    log,
		cfg:       cfg,
	}
}

func (s *Server) MCPServer() *server.MCPServer {
	return s.mcpServer
}

func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("Starting MCP server", "name", "synweb")
	return server.ServeStdio(s.mcpServer)
}
