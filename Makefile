# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Build the application
build:
	mkdir -p bin
	$(GOBUILD) -o $(BINARY_PATH) -v ./cmd/server

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_PATH)