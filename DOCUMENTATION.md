# YTDL-Next: Comprehensive Documentation

## Overview

YTDL-Next is a high-performance media streaming application that leverages [yt-dlp](https://github.com/yt-dlp/yt-dlp) for media extraction and delivers streams directly to the browser through a Go backend. The application is designed with a focus on memory efficiency, real-time delivery, and a modern user interface.

## Table of Contents

1. [Architecture](#architecture)
2. [Backend Components](#backend-components)
3. [Frontend Components](#frontend-components)
4. [Key Features](#key-features)
5. [Implementation Details](#implementation-details)
6. [Download Process](#download-process)
7. [Error Handling](#error-handling)
8. [Troubleshooting](#troubleshooting)
9. [Future Enhancements](#future-enhancements)

## Architecture

### Overview

YTDL-Next follows a client-server architecture:

```
Browser <-> Go HTTP Server <-> yt-dlp <-> Video Platforms
   ^             |
   |             v
   └─── Streaming Response
```

### Technologies

- **Frontend**: React 18+, TypeScript, Vite, Framer Motion for animations
- **Backend**: Go 1.x with standard HTTP libraries
- **Media Extraction**: [yt-dlp](https://github.com/yt-dlp/yt-dlp) CLI tool
- **Styling**: TailwindCSS

## Backend Components

### Core Modules

1. **Server Package** (`internal/server/`)
   - `server.go`: Main HTTP server implementation
   - `routes.go`: API route definitions
   - `middleware.go`: HTTP middleware for logging, CORS, etc.
   - `health-route.go`: Health check endpoint
   - `handle-yt-info.go`: Handler for video information fetching
   - `handle-yt-download.go`: Handler for video downloading
   - `download_manager.go`: Concurrent download management system

2. **Download API** (`internal/dl-api/`)
   - `cmd-interface.go`: Interface for command execution
   - `media.go`: Media extraction interface
   - `media-impl.go`: Implementation of media extraction

3. **Cache System** (`internal/cache/`)
   - `cache.go`: In-memory caching for video metadata

4. **Types** (`internal/types/`)
   - `yt.go`: Type definitions for YouTube data structures

5. **Embedded Files** (`internal/embedfs/`)
   - `embedfs.go`: Embedding of frontend assets
   - `web/`: Compiled frontend assets

### Main Entry Point

- `cmd/media-dl/main.go`: Application entry point that initializes and starts the server

## Frontend Components

1. **App Component** (`website/src/App.tsx`): 
   - Main application component
   - URL input and validation
   - Video information display
   - Format selection interface

2. **AudioSelector Component** (`website/src/AudioSelector.tsx`):
   - Modal for selecting audio formats for video-only downloads

3. **Utilities** (`website/src/util.ts`):
   - Helper functions for format conversion
   - Type definitions for media formats

## Key Features

### 1. Direct Memory-to-Browser Streaming

The application streams content directly from memory to the browser without saving intermediate files on the server's filesystem, enabling efficient resource usage.

**Implementation**: 
- Uses Go's `io.Pipe` to create an in-memory pipe between the yt-dlp process and the HTTP response
- Employs goroutines to handle the concurrent processing

### 2. Concurrent Download Management

A worker pool-based download manager handles multiple simultaneous download requests efficiently.

**Implementation**:
- `DownloadManager` struct manages a pool of worker goroutines
- Worker queue system for managing concurrent downloads
- Configurable maximum number of concurrent downloads

### 3. Format Selection

Users can select specific video/audio formats based on quality, codec, and file size.

**Implementation**:
- Frontend displays available formats with details like resolution, codec, and size
- Audio-only formats are separated for video-only selections
- Visual indicators for compatible formats

### 4. Temporary File Extension System

Downloads in progress use a `.ytdlp` extension, which helps users identify incomplete downloads.

**Implementation**:
- Server sets appropriate Content-Disposition headers
- Files are downloaded with temporary extension `.ytdlp`
- Users can manually rename files after download completion

### 5. Progress Feedback

The application provides visual feedback during downloads.

**Implementation**:
- Status messages in the UI
- Console messages indicating download status

## Implementation Details

### Download Manager

The Download Manager is implemented as a worker pool system in Go:

- Uses channels for communication between components
- Worker goroutines pick up download requests from a queue
- Active download count tracking
- Graceful shutdown capabilities

```go
// Core structures
type DownloadRequest struct {
    MediaURL     string
    FormatID     string
    Title        string
    ResponseChan chan *DownloadResponse
    Context      context.Context
}

type DownloadManager struct {
    maxWorkers       int
    downloadQueue    chan *DownloadRequest
    activeDownloads  int
    activeDownloadMu sync.Mutex
    workerWg         sync.WaitGroup
    shutdownCh       chan struct{}
}
```

### Download Process

The download process follows these steps:

1. User submits a URL for video information retrieval
2. Backend fetches video metadata using yt-dlp
3. Frontend displays available formats
4. User selects a format to download
5. For video-only formats, an audio format can be selected for combined download
6. Download request is enqueued in the Download Manager
7. Worker processes the download, streaming directly to the client
8. Download is delivered with a `.ytdlp` temporary extension

### Format Handling

The system supports various media formats:

- Single format downloads (audio or video)
- Combined format downloads (audio + video)
- Format compatibility checking

For combined formats, FFmpeg is used to mux the streams together.

### Error Handling

The application implements comprehensive error handling:

- Client disconnection detection
- Download process cancellation
- Invalid URL validation
- Process execution error handling

### Content Length and Chunked Transfer Encoding

For combined formats (audio+video), the system uses chunked transfer encoding instead of Content-Length headers to prevent browser download issues:

```go
// For compound formats or when size is unknown, use chunked transfer encoding
if isCompound || downloadResp.FileSize <= 0 {
    w.Header().Set("Transfer-Encoding", "chunked")
} else {
    // Set Content-Length when we have an accurate size for non-compound formats
    w.Header().Set("Content-Length", fmt.Sprintf("%d", downloadResp.FileSize))
}
```

## Download Process Details

### 1. Video Information Retrieval

```
Client → GET /api/yt/info?url={url} → Server → yt-dlp → Response
```

The server executes yt-dlp with the `-J` flag to get JSON output containing all available formats.

### 2. Format Selection

The client displays formats sorted by:
1. Compatibility (compatible formats first)
2. File size (largest to smallest)

### 3. Download Initialization

```
Client → GET /api/yt/download?url={url}&format_id={format_id} → Server → Queue → Worker → yt-dlp → Streaming Response
```

### 4. Download Streaming

For single formats:
```
yt-dlp → stdout → io.Pipe → HTTP Response → Browser Download
```

For combined formats:
```
yt-dlp (audio) + yt-dlp (video) → FFmpeg muxing → stdout → io.Pipe → HTTP Response → Browser Download
```

### 5. Client Disconnection Handling

If the client disconnects during download:
1. The server detects the disconnection through context cancellation
2. The download process is killed
3. Resources are cleaned up

## Troubleshooting

### Common Issues

1. **Error: `command with a non-nil Cancel was not created with CommandContext`**
   - Cause: Incorrect command creation for cancellable processes
   - Solution: Use `exec.CommandContext()` for processes that need cancellation

2. **Error: `http: wrote more than the declared Content-Length`**
   - Cause: Inaccurate file size estimation for combined formats
   - Solution: Use chunked transfer encoding for combined formats

3. **JavaScript Error: `Cannot set properties of null (setting 'innerHTML')`**
   - Cause: Issues with iframe content manipulation
   - Solution: Use direct link creation for downloads

## Future Enhancements

1. **Server-side file renaming**: Automatically rename files after download completion
2. **Download progress indicator**: Real-time progress updates during download
3. **Download queue management UI**: Interface for managing download queue
4. **Persistent settings**: User preferences storage
5. **More platform support**: Beyond YouTube to other video platforms

---

This documentation is maintained as of May 7, 2025. For the latest information, check the project repository.