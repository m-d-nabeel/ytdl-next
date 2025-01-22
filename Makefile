.PHONY: build build-web build-server run

build: build-web build-server

build-run: build run

build-web:
	cd website && npm run build

build-server:
	go build -o media-dl cmd/media-dl/main.go

run:
	./media-dl

run-debug:
	cd website && npm run dev &	go run cmd/media-dl/main.go
