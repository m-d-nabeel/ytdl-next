# Mastering Command Line Programming in Go: Building a Video Downloader

## Table of Contents
1. Introduction
2. Understanding Process Management in Go
3. Command Execution Fundamentals
4. Streaming and Real-time Output Processing
5. Interactive CLI Applications
6. Advanced Topics and Best Practices
7. Building Our Video Downloader
8. Error Handling and Recovery
9. Performance Optimization
10. Deployment and Distribution

## 1. Introduction

Welcome to our comprehensive guide on building command-line applications in Go. We'll be creating a video downloader, but the concepts we'll cover extend far beyond this specific application. By the end of this guide, you'll understand not just how to write command-line tools, but why certain patterns exist and when to use them.

### Why Go for CLI Applications?

Go excels at building command-line tools for several reasons. Its standard library provides robust process management capabilities, it compiles to a single binary, and it offers excellent cross-platform support. The language's concurrent features make it particularly well-suited for handling multiple streams of input and output, which is crucial for interactive CLI applications.

### Prerequisites

Before diving in, you should have:
- Basic Go programming knowledge
- Go installed on your system (version 1.11 or later recommended)
- Understanding of terminal/command prompt usage
- Familiarity with basic concurrency concepts

## 2. Understanding Process Management in Go

### The Process Model

At its core, every command-line program we run is a process. When our Go program executes another program (like a video downloader), it creates a child process. Understanding this relationship is crucial for building reliable CLI applications.

Let's explore what happens when we run a command:

1. The operating system creates a new process
2. The process gets its own memory space
3. Standard streams (stdin, stdout, stderr) are established
4. The program begins execution

In Go, the `os/exec` package manages this process creation and management. Here's a simple example:

```go
cmd := exec.Command("ls", "-l")
```

This line doesn't actually run the command yet - it just creates a `Cmd` struct that represents our future process. Think of it as a blueprint for what we want to run.

### Process States and Lifecycle

A process in Go goes through several states:
1. Created (New)
2. Running
3. Waiting
4. Terminated

Understanding these states is crucial because different methods are available at different states:

```go
cmd := exec.Command("echo", "hello") // State: Created
err := cmd.Start()                   // State: Running
err = cmd.Wait()                     // State: Waiting -> Terminated
```

Common mistake: Trying to access output before starting the process or trying to start a process twice.

## 3. Command Execution Fundamentals

### Simple Command Execution

Let's start with the most basic way to run a command. While not suitable for our final video downloader, understanding this helps build a foundation:

```go
output, err := exec.Command("youtube-dl", url).CombinedOutput()
```

Why this isn't ideal for a video downloader:
1. Blocks until the command completes
2. Stores all output in memory
3. Provides no real-time feedback
4. Cannot interact with the process

### Understanding Command Construction

The `exec.Command` function is more sophisticated than it appears. Let's break down what it does:

1. Path Resolution: It looks for the executable in the system PATH
2. Argument Processing: It handles argument escaping and quoting
3. Environment Setup: It inherits the current process's environment

Common pitfall: Assuming the command will run in a specific directory. Always be explicit about working directories:

```go
cmd := exec.Command("youtube-dl", url)
cmd.Dir = "/downloads"  // Explicitly set working directory
```

### Environment and Context

Every command runs within an environment. Go gives us fine-grained control over this:

```go
cmd := exec.Command("youtube-dl", url)
cmd.Env = append(os.Environ(),
    "DOWNLOAD_QUALITY=best",
    "DOWNLOAD_PATH=/videos",
)
```

The environment isn't just about variables - it's about the complete context in which our command runs. This includes:
- Working directory
- Environment variables
- User permissions
- Resource limits

## 4. Streaming and Real-time Output Processing

### Understanding I/O Streams

Every process has three standard streams:
1. Standard Input (stdin)
2. Standard Output (stdout)
3. Standard Error (stderr)

Think of these streams as pipes connecting processes. In our video downloader, we need to:
- Read progress information from stdout
- Handle errors from stderr
- Potentially send commands through stdin

### Real-time Output Processing

The key to building a responsive CLI application is processing output in real-time. Let's break down the pattern:

```go
stdout, err := cmd.StdoutPipe()
if err != nil {
    return err
}

scanner := bufio.NewScanner(stdout)
for scanner.Scan() {
    line := scanner.Text()
    // Process each line as it arrives
}
```

Why use a scanner instead of direct reading?
1. Automatic line buffering
2. Handles different line endings
3. More efficient than reading byte-by-byte
4. Cleaner error handling

### Concurrent Stream Processing

In real applications, we need to handle multiple streams concurrently. Go's goroutines make this elegant:

```go
go processOutput(stdout, "OUT")
go processOutput(stderr, "ERR")
```

But concurrent processing introduces challenges:
1. Synchronization: When can we consider the command complete?
2. Resource management: How do we ensure all goroutines clean up?
3. Error handling: How do we propagate errors from goroutines?

The solution involves careful orchestration with WaitGroups and channels:

```go
var wg sync.WaitGroup
wg.Add(2)  // One for stdout, one for stderr

go func() {
    defer wg.Done()
    processOutput(stdout)
}()

go func() {
    defer wg.Done()
    processOutput(stderr)
}()

wg.Wait()  // Wait for both streams to complete
```

## 5. Interactive CLI Applications

### User Input Handling

Interactive CLI applications need to handle user input while simultaneously processing output. This creates an interesting challenge: how do we handle multiple input sources concurrently?

The solution involves multiplexing:
1. Input from the user (stdin)
2. Output from the process (stdout/stderr)
3. System signals (like Ctrl+C)

Here's a pattern for handling this:

```go
func handleUserInput(stdin io.WriteCloser, done chan bool) {
    defer stdin.Close()
    
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        select {
        case <-done:
            return
        default:
            fmt.Fprintln(stdin, scanner.Text())
        }
    }
}
```

### Progress Reporting

For a video downloader, progress reporting is crucial. Let's explore different ways to show progress:

1. Simple percentage:
```go
fmt.Printf("\rProgress: %d%%", percentage)
```

2. Progress bar:
```go
const barLength = 50
complete := int(float64(barLength) * (float64(downloaded) / float64(total)))
fmt.Printf("\r[%s%s] %d%%",
    strings.Repeat("=", complete),
    strings.Repeat(" ", barLength-complete),
    percentage,
)
```

The key to good progress reporting is:
1. Update frequency (not too often, not too rare)
2. Visual clarity
3. Information density
4. Terminal compatibility

## 6. Advanced Topics and Best Practices

### Resource Management

Long-running processes can consume significant resources. Best practices include:

1. Setting resource limits:
```go
cmd.SysProcAttr = &syscall.SysProcAttr{
    Setpgid: true,
}
```

2. Implementing timeouts:
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
defer cancel()
cmd := exec.CommandContext(ctx, "youtube-dl", url)
```

3. Proper cleanup on exit:
```go
defer func() {
    if cmd.Process != nil {
        cmd.Process.Kill()
    }
}()
```

### Error Handling Strategies

Error handling in CLI applications requires special attention:

1. Distinguishing between different types of errors:
   - Process errors (exit code != 0)
   - I/O errors
   - User interruption
   - Resource exhaustion

2. Providing meaningful error messages:
```go
if err != nil {
    if exitErr, ok := err.(*exec.ExitError); ok {
        return fmt.Errorf("download failed with code %d: %s",
            exitErr.ExitCode(),
            string(exitErr.Stderr),
        )
    }
    return fmt.Errorf("unexpected error: %v", err)
}
```

### Signal Handling

Proper signal handling ensures our application behaves well when interrupted:

```go
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

go func() {
    <-sigChan
    fmt.Println("\nReceived interrupt signal. Cleaning up...")
    if cmd.Process != nil {
        cmd.Process.Kill()
    }
    os.Exit(1)
}()
```

## 7. Building Our Video Downloader

Now that we understand the fundamentals, let's build our video downloader. We'll use a structured approach:

### Command Executor Interface

First, define our interface:

```go
type Downloader interface {
    Download(url string) error
    Cancel() error
    Progress() float64
}
```

This interface abstraction allows us to:
1. Switch between different download tools
2. Mock the downloader for testing
3. Provide a clean API for other parts of our application

### Implementation Considerations

When implementing the downloader:

1. Progress parsing:
   - Different tools provide progress in different formats
   - Need to handle incomplete/corrupt output
   - Consider rate limiting progress updates

2. Error recovery:
   - Network interruptions
   - Disk space issues
   - Invalid URLs
   - Rate limiting from services

3. User feedback:
   - Download speed
   - ETA
   - File size
   - Quality information

## 8. Error Handling and Recovery

### Robust Error Handling

In a production-grade CLI application, error handling needs to be comprehensive:

1. Network errors:
```go
if err != nil {
    if netErr, ok := err.(net.Error); ok {
        if netErr.Temporary() {
            // Implement retry logic
            return retry(func() error {
                return Download(url)
            }, 3, time.Second)
        }
    }
}
```

2. Resource errors:
```go
if err != nil {
    if pathErr, ok := err.(*os.PathError); ok {
        switch {
        case os.IsNotExist(pathErr):
            // Handle missing file/directory
        case os.IsPermission(pathErr):
            // Handle permission issues
        }
    }
}
```

### Recovery Strategies

For long-running downloads, implement recovery mechanisms:

1. Checkpoint system:
```go
type DownloadState struct {
    URL           string
    BytesReceived int64
    TotalBytes    int64
    Timestamp     time.Time
}

func saveCheckpoint(state DownloadState) error {
    // Save state to disk
}

func resumeDownload(state DownloadState) error {
    // Resume from saved state
}
```

2. Partial file handling:
```go
func resumableDownload(url string, existingFile string) error {
    info, err := os.Stat(existingFile)
    if err == nil {
        // Add range header to request
        headers["Range"] = fmt.Sprintf("bytes=%d-", info.Size())
    }
    // Continue download
}
```

## 9. Performance Optimization

### Buffering Strategies

Proper buffering is crucial for performance:

```go
// Create a buffered reader for better performance
reader := bufio.NewReaderSize(stdout, 16*1024)  // 16KB buffer

// Use a buffer pool for output processing
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 32*1024)
    },
}
```

### Concurrent Downloads

For multiple files:

```go
func downloadWorker(jobs <-chan string, results chan<- error) {
    for url := range jobs {
        results <- Download(url)
    }
}

func downloadMany(urls []string, concurrency int) error {
    jobs := make(chan string, len(urls))
    results := make(chan error, len(urls))

    // Start workers
    for i := 0; i < concurrency; i++ {
        go downloadWorker(jobs, results)
    }

    // Send jobs
    for _, url := range urls {
        jobs <- url
    }
    close(jobs)

    // Collect results
    for range urls {
        if err := <-results; err != nil {
            return err
        }
    }
    return nil
}
```

## 10. Deployment and Distribution

### Building for Distribution

Create a proper build process:

```makefile
.PHONY: build
build:
    go build -ldflags="-s -w" -o downloader main.go

.PHONY: release
release:
    GOOS=linux GOARCH=amd64 go build -o downloader-linux-amd64
    GOOS=darwin GOARCH=amd64 go build -o downloader-darwin-amd64
    GOOS=windows GOARCH=amd64 go build -o downloader-windows-amd64.exe
```

### Configuration Management

Use configuration files for flexibility:

```go
type Config struct {
    DefaultQuality string `json:"default_quality"`
    DownloadPath   string `json:"download_path"`
    Concurrent     int    `json:"concurrent_downloads"`
    RetryAttempts  int    `json:"retry_attempts"`
}

func loadConfig() (*Config, error) {
    data, err := os.ReadFile("config.json")
    if err != nil {
        return defaultConfig()
    }
    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }
    return &config, nil
}
```
