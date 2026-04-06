# synweb

MCP server for Synthetic Web Search API, written in Go.

## Features

- **MCP Protocol**: Implements Model Context Protocol for AI integration
- **Circuit Breaker**: Protects against cascading failures using failsafe-go
- **Retry Logic**: Automatic retries with exponential backoff
- **Structured Logging**: DEBUG/INFO/WARN/ERROR levels via slog

## Requirements

- Go 1.21+
- Synthetic API key

## Installation

```bash
git clone https://github.com/martinstidelius/synweb.git
cd synweb
go mod tidy
make build
```

## Configuration

Copy `.env.example` to `.env` and set your API key:

```bash
cp .env.example .env
# Edit .env with your SYNTHETIC_API_KEY
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SYNTHETIC_API_KEY` | Your Synthetic API key | (required) |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |
| `API_BASE_URL` | Synthetic API endpoint | `https://api.synthetic.new/v2` |

## Usage

### Build and Run

```bash
make build
./synweb
```

### Development

```bash
make run-dev
```

### Available Commands

```bash
make build    # Build the binary
make run      # Build and run
make run-dev  # Run in development mode
make lint     # Run linter
make test     # Run tests
make clean    # Clean build artifacts
make tidy     # Update dependencies
```

## MCP Tool

### search_web

Search the web using Synthetic API.

**Parameters:**
- `query` (required): Search query string
- `max_text_length` (optional): Maximum characters in result text (default: 1000)

**Example:**

```json
{
  "name": "search_web",
  "arguments": {
    "query": "Go programming best practices",
    "max_text_length": 500
  }
}
```

## Testing with MCP Inspector

```bash
npx @modelcontextprotocol/inspector ./synweb
```

## Opencode Configuration

To use this MCP server with Opencode, add it to your Opencode configuration file (`~/.config/opencode/opencode.json`):

```json
{
  "mcp": {
    "synweb": {
      "command": ["/absolute/path/to/synweb"],
      "enabled": true,
      "type": "local",
      "environment": {
        "SYNTHETIC_API_KEY": "your-api-key-here",
        "LOG_LEVEL": "info"
      }
    }
  }
}
```

Replace `/absolute/path/to/synweb` with the actual path to the synweb binary. You can get the path by running:

```bash
pwd # if you're in the synweb directory
# or
echo "$(pwd)/synweb"
```

See [Environment Variables](#environment-variables) above for available configuration options.

## Example Prompts

Here are some example prompts you can use when interacting with the MCP server through your AI assistant:

### Basic Search

> "Search the web for the latest Go 1.26 release notes"

> "Find information about the best practices for REST API design in 2024"

### Technical Research

> "Search for articles about PostgreSQL performance optimization techniques"

> "Find documentation on building MCP servers in Go"

### Current Events

> "Search for recent news about artificial intelligence developments"

> "Find the latest information about TypeScript 6.0 features"

### Code Examples

> "Search for React hooks tutorial with useEffect examples"

> "Find Python asyncio best practices for production code"

### Troubleshooting

> "Search for solutions to fix 'connection refused' errors in Docker"

> "Find how to resolve merge conflicts in Git"

### Customizing Results

You can also use the `max_text_length` parameter to get more or less detail:

> "Search for 'Kubernetes deployment strategies' and give me the full text (max_text_length: 5000)"

## License

MIT