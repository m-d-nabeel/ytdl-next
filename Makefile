.PHONY: all build build-web build-server run run-dev run-prod client-run server-run clean

# Default target
all: build

# Parallelizable build targets
build: build-web build-server

build-server: copy-static-files
	@echo "Building server..."
	go build -o media-dl cmd/media-dl/main.go

build-web:
	@echo "Building website..."
	cd website && bun install && bun run build

copy-static-files: build-web
	@echo "Copying static files for embedding..."
	mkdir -p internal/embedfs/web
	cp -r website/dist/* internal/embedfs/web/

# Run targets
run: run-prod

run-dev: build
	@echo "Running application in development mode..."
	DEV_MODE=true ./media-dl

run-prod: build
	@echo "Running application in production mode..."
	DEV_MODE=false ./media-dl

# Development targets
client-run:
	@echo "Starting client in development mode..."
	cd website && bun run dev

server-run:
	@echo "Starting server in development mode..."
	DEV_MODE=true go run cmd/media-dl/main.go

# Combined development mode (requires tmux or multiple terminals)
dev:
	@echo "Starting development environment (run client and server in separate terminals)..."
	@echo "In terminal 1: make client-run"
	@echo "In terminal 2: make server-run"

# Cleanup
clean:
	@echo "Cleaning up..."
	rm -f media-dl
	cd website && rm -rf dist node_modules
	rm -rf internal/embedfs/web
