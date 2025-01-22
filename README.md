# Media Stream Downloader

A high-performance media streaming application that leverages yt-dlp for extraction and provides direct browser streaming through a Go backend. Built with a focus on memory efficiency and real-time delivery.

## Core Features

- Direct memory-to-browser streaming without local storage
- Real-time media extraction and delivery
- Audio quality selection and format controls
- React-based web interface with TypeScript
- Efficient Go backend with concurrent stream handling

## Architecture

- **Frontend**: React + TypeScript application built with Vite
- **Backend**: Go server handling media extraction and streaming
- **Core**: yt-dlp integration for reliable media source extraction

## Prerequisites

- Go 1.x
- yt-dlp installed and in system PATH
- JavaScript runtime (Node.js 16+, Deno, or Bun)

## Quick Start

1. Install yt-dlp
2. Install frontend dependencies:
```bash
cd website
npm install    # or bun install / yarn
```

3. Build and run:
```bash
make build-run
```

Development mode with hot reloading:
```bash
make run-debug
```

## System Design

```
Browser <-> Go Server <-> yt-dlp
   ^          |
   |          v
   └── Streaming Response
```

