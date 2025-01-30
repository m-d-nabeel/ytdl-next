.PHONY: build build-web build-server run client-run server-run clean

# Build targets
build-server:
	@echo "Building server..."
	go build -o media-dl cmd/media-dl/main.go

build-web:
	@echo "Building website..."
	cd website && bun install && bun run build

build: build-web build-server

# Run targets
run: build
	@echo "Running application..."
	./media-dl

# Development targets
client-run:
	@echo "Starting client in development mode..."
	cd website && bun run dev

server-run:
	@echo "Starting server in development mode..."
	go run cmd/media-dl/main.go

# Cleanup
clean:
	@echo "Cleaning up..."
	rm -f media-dl
	cd website && rm -rf dist node_modules
