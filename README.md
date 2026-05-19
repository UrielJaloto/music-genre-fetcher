# 🎵 music-genre-fetcher

A highly concurrent, SOLID-compliant CLI automation tool built in Go to fetch, curate, and catalog accurate music genres for local audio tracks using the Last.fm API.

<p align="center">
  <img src="https://s6.ezgif.com/tmp/ezgif-6ae82a565797634e.gif" width="450" height="110" alt="Music Visualizer">
</p>

## Overview

Managing metadata for local music libraries can be a tedious process. This project automates genre retrieval by querying track and artist details asynchronously, extracting the top 3 most accurate crowdsourced tags while intelligently filtering out irrelevant numerical metadata like release years.

The application is engineered for high throughput and reliability, functioning as a highly scalable, production-ready pipeline for metadata aggregation capable of resolving up to 1,500 tracks per minute while maintaining a memory footprint well under 15MB.

## Key Features

* **High-Performance Concurrency:** Utilizes a robust worker-pool pipeline to safely process thousands of records asynchronously, maximizing network utilization and minimizing idle time.
* **Rate-Limit Resilience:** Built-in backoff mechanism that handles HTTP `429` and `5xx` server strains gracefully with automatic retries, ensuring zero dropped payloads during sustained execution.
* **Polymorphic Input Parsers:** Seamlessly switches between custom structured TXT datasets (`|||` delimited) and standard JSON structures based on file extensions.
* **Pluggable Architecture:** Features modular multi-format data exporters (JSON, CSV, TXT) that allow seamless extension for future file formats without touching the core logic.
* **Interactive CLI Context:** A clean console interface guiding the operator through ingestion and distribution with real-time feedback.

## Architecture & Internal Design

This codebase is structured around pragmatic, scalable software engineering practices:

* **Decoupled Components:** The core application logic depends exclusively on abstract interfaces. Implementations are separated into distinct presentation and infrastructure layers for maximum testability.
* **Dependency Injection:** The main execution serves purely as an assembly point where dependencies are instantiated and injected into the execution context, ensuring modularity and isolated state management.
* **Extensibility:** Output formatters and data providers are designed as discrete, interchangeable objects. Adding a new API provider or output format requires zero modifications to the core engine.
* **Memory Optimization:** Heavy data transformations avoid forcing pointers to the heap prematurely, utilizing Go’s stack allocations effectively during large I/O operations.

## Installation

### Option 1: Pre-compiled Binaries (Recommended)
You do not need Go installed to run this tool. 

1. Head over to the [Releases](https://github.com/UrielJaloto/music-genre-fetcher/releases) page.
2. Download the executable tailored for your Operating System (Windows, macOS, or Linux).
3. Place the executable in your desired working directory and proceed to the **Configuration Setup** step.

### Option 2: Build from Source
If you prefer to compile the application yourself:

1. Clone the repository:
   ```bash
   git clone [https://github.com/UrielJaloto/music-genre-fetcher.git](https://github.com/UrielJaloto/music-genre-fetcher.git)
   cd music-genre-fetcher
2. Build the executable:
    ```bash
    go build -o mgf internal/cmd/mgf/main.go
## Configuration Setup

Regardless of your installation method, the application requires a specific folder structure relative to the executable for safe I/O operations.

1. Create the following directory tree next to your executable:
    ```text
    env/
    ├── config.json
    ├── input/
    └── output/
2. Configure your Last.fm API credentials inside `env/config.json`:
    ```json
    {
      "lastfm_api_key": "YOUR_LASTFM_API_KEY"
    }
## Usage

Place your input files inside the `env/input/` directory.

### Supported Formats

#### Structured TXT (`env/input/mp3tag_data.txt`)

```text
Path|||Title|||Artist
C:\Music\Rock\Queen\bohemian_rhapsody.mp3|||Bohemian Rhapsody|||Queen
```

#### JSON (`env/input/data.json`)

```json
[
  {
    "path": "C:\\Music\\Pop\\Thriller.mp3",
    "title": "Thriller",
    "artist": "Michael Jackson"
  }
]
```

### Execution

Launch the executable from your terminal. If you built from source:

```bash
go run internal/cmd/mgf/main.go
```

If using the downloaded binary:

```bash
./mgf
```

*(On Windows, just double-click the `.exe` or run `mgf.exe` in CMD/PowerShell).*

The application will prompt you via terminal to specify your file path or let you press `Enter` to use the fallback default configuration.

## The Story Behind the Project

In an era dominated by algorithmic cloud streaming platforms, maintaining a deeply curated, local music collection offers a unique level of ownership. To me, managing an offline music archive is the modern equivalent of organizing a physical vinyl shelf.

However, anyone who preserves local files knows the persistent friction of chaotic metadata—specifically, genre tagging. Music downloaded or ripped from disparate sources rarely comes with standardized genre attributes. Manually analyzing and rewriting tags for thousands of songs quickly becomes an exhausting, bottlenecked task.

I engineered this tool to solve my own problem. By turning a repetitive manual chore into an optimized Go application, I reduced hours of manual data entry into a process that completes in a matter of seconds. What started as a personal utility evolved into a highly reliable pipeline—providing a practical environment to apply advanced concurrency, strict architectural decoupling, and robust network resilience.
