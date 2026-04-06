.PHONY: build run lint test clean

BINARY_NAME=synweb
GO_CMD=go
MAIN_PATH=./cmd/synweb

build:
	$(GO_CMD) build -o $(BINARY_NAME) $(MAIN_PATH)

run: build
	./$(BINARY_NAME)

run-dev:
	SYNTHETIC_API_KEY=test LOG_LEVEL=debug $(GO_CMD) run $(MAIN_PATH)

lint:
	golangci-lint run ./...

test:
	$(GO_CMD) test -v -race -cover ./...

clean:
	rm -f $(BINARY_NAME)

tidy:
	$(GO_CMD) mod tidy

vet:
	$(GO_CMD) vet ./...

fmt:
	$(GO_CMD) fmt ./...

all: fmt vet lint test build