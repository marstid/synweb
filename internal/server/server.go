package server

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/marstid/synweb/internal/logger"
	"github.com/marstid/synweb/internal/search"
)

type Server struct {
	mcpServer *server.MCPServer
	logger    *logger.Logger
}

func New(name, version string, log *logger.Logger, handler *search.Handler) *Server {
	s := server.NewMCPServer(
		name,
		version,
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
			mcp.DefaultNumber(1000),
			mcp.Description("Maximum number of characters to include in the 'text' field of each result. Defaults to 1000. Set higher if more detail is needed, or fetch the source URL directly for the full content."),
		),
	)

	s.AddTool(searchTool, handler.HandleSearch)

	return &Server{
		mcpServer: s,
		logger:    log,
	}
}

func (s *Server) MCPServer() *server.MCPServer {
	return s.mcpServer
}

func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("Starting MCP server", "name", "synweb")
	return server.ServeStdio(s.mcpServer)
}
