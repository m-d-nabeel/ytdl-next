.PHONY: build build-web build-server run debug-run

build-server:
	go build -o media-dl cmd/media-dl/main.go

build-web:
	cd website && bun install && bun run build

build: build-web build-server

build-run: build
	./media-dl

debug-run:
	cd website && bun run dev & go run cmd/media-dl/main.go
