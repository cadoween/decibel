# decibel

[![Go Version](https://img.shields.io/github/go-mod/go-version/cadoween/decibel)](go.mod)

decibel is a overpowered command-line tool for analyzing and managing your music listening history, with current support for Spotify data. It processes your Spotify streaming history data and provides insights into your listening habits.

## Features

- Import Spotify streaming history from extended history files.
- Store listening data in SQLite database, so you don't have to spawn any kind of external processes.
- Support for concurrent and batch processing of large datasets.
- Detailed verbose logging.

## Prerequisites

- Go 1.23.4 or higher.
- SQLite database.
- Spotify account data export for extended history.

## Installation

```bash
# Clone the repository
git clone https://github.com/cadoween/decibel.git
cd decibel

# Install dependencies
go mod download

# Build the project
go build -o decibel cmd/main.go
```

## Usage

### Basic Commands

```bash
# Import Spotify streaming history
decibel spotify seeder run --db ./path/to/database.db --dir "./path/to/spotify/data" --verbose

# Using make command (predefined paths)
make spotify-seeder-run
```

### Command Structure

```
decibel
├── spotify
│   └── seeder
│       └── run [flags]
```

### Available Flags

- `--db`: Path to the SQLite database file (required)
- `--dir`: Directory containing Spotify Extended Streaming History (required)
- `--verbose, -v`: Enable verbose logging (optional)

## Data Structure

### Spotify Stream Data

The tool processes the following data points for each stream:

- Timestamp
- Username
- Platform
- Play Duration (ms)
- Connection Country
- IP Address (if available)
- User Agent
- Track Information
  - Track Name
  - Artist Name
  - Album Name
  - Spotify URI
- Playback Information
  - Start/End Reason
  - Shuffle Status
  - Skip Status
  - Offline Status
  - Incognito Mode

## Development

### Adding New Features

1. Create new command in `cmd/` directory.
2. Implement business logic in `internal/` directory.
3. Update documentation in `README.md`.
