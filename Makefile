.PHONY: all build clean test dist dist-linux lint

# Build variables
PLUGIN_ID=com.fambear.expiring-messages
PLUGIN_VERSION=$(shell jq -r '.version' plugin.json)
BINARY_NAME=plugin
GO=go
GOFLAGS=-mod=mod -trimpath
SERVER_DIR=server
WEBAPP_DIR=webapp
DIST_DIR=dist

# Build the plugin
all: build

build:
	@echo "Building server..."
	cd $(SERVER_DIR) && $(GO) build $(GOFLAGS) -o ../$(DIST_DIR)/$(BINARY_NAME)-linux-amd64
	@echo "Building webapp..."
	cd $(WEBAPP_DIR) && npm ci && npm run build
	@mkdir -p $(DIST_DIR)/webapp/dist
	@cp $(WEBAPP_DIR)/dist/main.js $(DIST_DIR)/webapp/dist/

clean:
	@echo "Cleaning..."
	rm -rf $(DIST_DIR)
	rm -rf $(WEBAPP_DIR)/dist
	rm -rf $(WEBAPP_DIR)/node_modules

test:
	@echo "Running server tests..."
	cd $(SERVER_DIR) && $(GO) test -v ./...
	@echo "Running webapp tests..."
	cd $(WEBAPP_DIR) && npm ci && npm run test

# Build for Linux AMD64 only (CI)
dist-linux: clean
	@echo "Building distribution for Linux AMD64..."
	@mkdir -p $(DIST_DIR)/webapp/dist
	@cp plugin.json $(DIST_DIR)/
	@cd $(SERVER_DIR) && $(GO) build $(GOFLAGS) -o ../$(DIST_DIR)/$(BINARY_NAME)-linux-amd64
	@cd $(WEBAPP_DIR) && npm ci && npm run build
	@cp $(WEBAPP_DIR)/dist/main.js $(DIST_DIR)/webapp/dist/
	@cd $(DIST_DIR) && tar -czvf expiring-messages-$(PLUGIN_VERSION).tar.gz plugin.json $(BINARY_NAME)-linux-amd64 webapp
	@echo "Distribution created: $(DIST_DIR)/expiring-messages-$(PLUGIN_VERSION).tar.gz"

# Build for all platforms
dist: clean
	@echo "Building distribution for all platforms..."
	@mkdir -p $(DIST_DIR)/webapp/dist
	@cp plugin.json $(DIST_DIR)/
	@cd $(SERVER_DIR) && $(GO) build $(GOFLAGS) -o ../$(DIST_DIR)/$(BINARY_NAME)-linux-amd64
	@cd $(SERVER_DIR) && GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o ../$(DIST_DIR)/$(BINARY_NAME)-darwin-amd64
	@cd $(SERVER_DIR) && GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -o ../$(DIST_DIR)/$(BINARY_NAME)-darwin-arm64
	@cd $(SERVER_DIR) && GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o ../$(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe
	@cd $(WEBAPP_DIR) && npm ci && npm run build
	@cp $(WEBAPP_DIR)/dist/main.js $(DIST_DIR)/webapp/dist/
	@cd $(DIST_DIR) && tar -czvf expiring-messages-$(PLUGIN_VERSION).tar.gz plugin.json $(BINARY_NAME)-* webapp
	@echo "Distribution created: $(DIST_DIR)/expiring-messages-$(PLUGIN_VERSION).tar.gz"

lint:
	@echo "Running Go linters..."
	cd $(SERVER_DIR) && golangci-lint run
	@echo "Running TypeScript linters..."
	cd $(WEBAPP_DIR) && npm run lint
