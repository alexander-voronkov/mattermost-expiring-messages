.PHONY: all build clean test deploy dist watch

# Build variables
BINARY_NAME=plugin
GO=go
GOFLAGS=-mod=mod
SERVER_DIR=server
WEBAPP_DIR=webapp

# Build the plugin
all: build

build:
	@echo "Building server..."
	cd $(SERVER_DIR) && $(GO) build $(GOFLAGS) -o dist/$(BINARY_NAME)-linux-amd64
	@echo "Building webapp..."
	cd $(WEBAPP_DIR) && npm run build

clean:
	@echo "Cleaning..."
	rm -rf $(SERVER_DIR)/dist
	rm -rf $(WEBAPP_DIR)/dist

test:
	@echo "Running tests..."
	cd $(SERVER_DIR) && $(GO) test ./...

deploy: build
	@echo "Deploying plugin..."
ifdef MM_LOCALSOCKETPATH
		@echo "Using local socket: $(MM_LOCALSOCKETPATH)"
		mkdir -p $(MM_LOCALSOCKETPATH)/plugins/com.fambear.expiring-messages
		cp $(SERVER_DIR)/dist/$(BINARY_NAME)-linux-amd64 $(MM_LOCALSOCKETPATH)/plugins/com.fambear.expiring-messages/plugin-linux-amd64
		cp $(WEBAPP_DIR)/dist/main.js $(MM_LOCALSOCKETPATH)/plugins/com.fambear.expiring-messages/
		cp plugin.json $(MM_LOCALSOCKETPATH)/plugins/com.fambear.expiring-messages/
else
		@echo "MM_LOCALSOCKETPATH not set. Skipping deployment."
endif

dist:
	@echo "Building distribution..."
	@mkdir -p dist
	@cp plugin.json dist/
	@cd $(SERVER_DIR) && $(GO) build $(GOFLAGS) -o ../dist/$(BINARY_NAME)-linux-amd64
	@cd $(SERVER_DIR) && GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o ../dist/$(BINARY_NAME)-darwin-amd64
	@cd $(SERVER_DIR) && GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -o ../dist/$(BINARY_NAME)-darwin-arm64
	@cd $(SERVER_DIR) && GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o ../dist/$(BINARY_NAME)-windows-amd64.exe
	@cd $(WEBAPP_DIR) && npm run build
	@cp $(WEBAPP_DIR)/dist/main.js dist/
	@cd dist && tar -czvf mattermost-expiring-messages.tar.gz plugin.json $(BINARY_NAME)-* main.js
	@echo "Distribution created in dist/"

watch:
	@echo "Watching webapp for changes..."
ifdef MM_SERVICESETTINGS_SITEURL
ifdef MM_ADMIN_TOKEN
		cd $(WEBAPP_DIR) && npm run watch
else
		@echo "MM_ADMIN_TOKEN not set. Cannot watch."
endif
else
		@echo "MM_SERVICESETTINGS_SITEURL not set. Cannot watch."
endif

lint:
	@echo "Running linters..."
	cd $(SERVER_DIR) && golangci-lint run
