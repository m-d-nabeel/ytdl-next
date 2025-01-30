# Media Stream Downloader  

A high-performance media streaming application that leverages [yt-dlp](https://github.com/yt-dlp/yt-dlp) for media extraction and delivers streams directly to the browser through a Go backend. Designed with a focus on memory efficiency and real-time delivery.  

## Core Features  

- **Direct Streaming:** Memory-to-browser streaming without intermediate local storage.  
- **Real-Time Delivery:** Instant media extraction and streaming.  
- **Flexible Controls:** Audio quality and format selection.  
- **Modern Frontend:** React-based web interface with TypeScript.  
- **Efficient Backend:** Concurrent stream handling using Go.  

## Architecture  

- **Frontend:** React + TypeScript application built with Vite.  
- **Backend:** Go server managing media extraction and streaming.  
- **Core Integration:** Reliable media extraction via [yt-dlp](https://github.com/yt-dlp/yt-dlp).  

## Prerequisites  

- **Go 1.x** installed.  
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) installed and accessible in the system PATH.  
- **JavaScript runtime**: Node.js 16+, Deno, or Bun.  

## Quick Start  

1. Install frontend dependencies:  
   ```bash
   cd website
   npm install    # or bun install / yarn
   ```  

2. Build and run the project:  
   ```bash
   make build-run
   ```  

## System Design  

```plaintext
Browser <-> Go Server <-> yt-dlp
   ^          |
   |          v
   └── Streaming Response
```  

## License  

Note: This project uses [yt-dlp](https://github.com/yt-dlp/yt-dlp), which is licensed under the Unlicense/Public Domain (source code) and GNU GPL v3+ (binary/installer).  

For more details, see the [yt-dlp license](https://github.com/yt-dlp/yt-dlp/blob/master/LICENSE).  

