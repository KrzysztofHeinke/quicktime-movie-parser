BINARY_NAME=quicktime-movie-parser
DOCKER_IMAGE=quicktime-movie-parser-builder
CONTAINER_NAME=quicktime-movie-parser-container
OUTPUT_DIR=bin
BUILD_WINDOWS=$(OUTPUT_DIR)/windows
BUILD_LINUX=$(OUTPUT_DIR)/linux

.PHONY: build
build: clean build-linux build-windows

.PHONY: build-linux
build-linux: clean
	mkdir -p $(BUILD_LINUX)
	@echo "Building the Go application in Docker..."

	# Check if the container exists and remove it if it does
	@docker rm -f $(CONTAINER_NAME) 2>/dev/null || true

	# Build the Docker image
	@docker build --build-arg GOOS=linux --build-arg GOARCH=amd64  --target builder -t $(DOCKER_IMAGE) .

	# Create the container without starting it
	@docker create --name $(CONTAINER_NAME) $(DOCKER_IMAGE)

	# Copy the binary from the container to the host
	@docker cp $(CONTAINER_NAME):/app/bin/$(BINARY_NAME) $(BUILD_LINUX)/$(BINARY_NAME)

	# Remove the container
	@docker rm $(CONTAINER_NAME)

	# Set execute permissions on the binary
	@chmod +x $(BUILD_LINUX)/$(BINARY_NAME)

	@echo "Binary copied to $(BUILD_LINUX)/$(BINARY_NAME) and permissions set."

.PHONY: build-windows
build-windows: clean
	mkdir -p $(BUILD_WINDOWS)
	@echo "Building the Go application in Docker..."

	# Check if the container exists and remove it if it does
	@docker rm -f $(CONTAINER_NAME) 2>/dev/null || true

	# Build the Docker image
	@docker build --build-arg GOOS=windows --build-arg GOARCH=amd64 --target builder -t $(DOCKER_IMAGE) .

	# Create the container without starting it
	@docker create --name $(CONTAINER_NAME) $(DOCKER_IMAGE)

	# Copy the binary from the container to the host
	@docker cp $(CONTAINER_NAME):/app/bin/$(BINARY_NAME) $(BUILD_WINDOWS)/$(BINARY_NAME).exe

	# Remove the container
	@docker rm $(CONTAINER_NAME)

	# Set execute permissions on the binary
	@chmod +x $(BUILD_WINDOWS)/$(BINARY_NAME).exe

	@echo "Binary copied to $(BUILD_WINDOWS)/$(BINARY_NAME) and permissions set."

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -f $(OUTPUT_DIR)/$(BINARY_NAME)
	@docker image rm $(DOCKER_IMAGE) 2>/dev/null || true
	rm -rf $(OUTPUT_DIR)
.PHONY: clean-image
clean-image:
	@echo "Removing Docker image..."
	@docker image rm $(DOCKER_IMAGE) 2>/dev/null || true
	rm -rf 

.PHONY: build-and-clean
build-and-clean: build clean-image
	@echo "Build complete and Docker image removed."

.PHONY: build-local
build-local: clean-local
	@echo "Building the Go application locally..."
	@go build -o $(OUTPUT_DIR)/$(BINARY_NAME) ./main.go
	@chmod +x $(OUTPUT_DIR)/$(BINARY_NAME)
	@echo "Local binary built and placed in $(OUTPUT_DIR)/$(BINARY_NAME)."

.PHONY: clean-local
clean-local:
	@echo "Cleaning up local build..."
	@rm -f $(OUTPUT_DIR)/$(BINARY_NAME)

.PHONY: test
test:
	@echo "Running tests locally..."
	@go test ./... -v -coverprofile=coverage.out
	@go tool cover -func=coverage.out

.PHONY: test-docker
test-docker: build
	@echo "Running tests in Docker..."
	@docker run --rm -v $(PWD):/app -w /app $(DOCKER_IMAGE) go test ./... -v -coverprofile=coverage.out
	@go tool cover -func=coverage.out

.PHONY: coverage
coverage:
	@echo "Generating code coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Code coverage report generated at coverage.html"

.PHONY: lint
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

.PHONY: lint-docker
lint-docker:
	@echo "Running linter in Docker..."
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint golangci-lint run -v ./...
