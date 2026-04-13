# SmirnovTorrent Usage Guide

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/ainsonet/smirnovtorrent.git
cd smirnovtorrent

# Build
make build

# Run
make run
```

### Pre-built Binaries

Download from [releases](https://github.com/ainsonet/smirnovtorrent/releases):
- **Linux**: `smirnovtorrent-linux`
- **Windows**: `smirnovtorrent.exe`
- **macOS**: `smirnovtorrent-macos`

## Quick Start

### Download a Torrent

```bash
# Using torrent file
./smirnovtorrent download path/to/file.torrent

# Using magnet link
./smirnovtorrent download "magnet:?xt=urn:btih:HASH"

# Specify output directory
./smirnovtorrent download file.torrent -o ~/Downloads
```

### Start Web UI

```bash
# Start with Web UI on port 8080
./smirnovtorrent webui --port 8080

# Access at http://localhost:8080
```

## Command Line Interface

### Commands

#### `download` - Download a torrent

```bash
smirnovtorrent download <torrent|magnet> [flags]
```

**Flags:**
- `-o, --output string` - Output directory (default: current directory)
- `-download-limit int` - Download speed limit in bytes/sec (0 = unlimited)
- `-upload-limit int` - Upload speed limit in bytes/sec (0 = unlimited)
- `--dht` - Enable DHT peer discovery (default: true)
- `--pex` - Enable Peer Exchange (default: true)
- `--encrypt` - Enable MSE encryption (default: true)

**Examples:**

```bash
# Basic download
smirnovtorrent download ubuntu-22.04.torrent

# Download with custom output
smirnovtorrent download magnet:?xt=urn:btih:HASH -o ~/Downloads

# Download with speed limits (1 MB/s download, 512 KB/s upload)
smirnovtorrent download file.torrent \
  --download-limit 1048576 \
  --upload-limit 524288

# Download with DHT and PEX enabled
smirnovtorrent download file.torrent --dht --pex

# Download with all options
smirnovtorrent download file.torrent \
  -o ~/Downloads \
  --download-limit 1048576 \
  --upload-limit 524288 \
  --dht \
  --pex \
  --encrypt
```

#### `webui` - Start Web UI server

```bash
smirnovtorrent webui [flags]
```

**Flags:**
- `--port int` - Web UI port (default: 8080)
- `--host string` - Web UI host (default: "localhost")
- `--download-dir string` - Default download directory

**Examples:**

```bash
# Start Web UI on default port
smirnovtorrent webui

# Start Web UI on custom port
smirnovtorrent webui --port 9090

# Start Web UI accessible from network
smirnovtorrent webui --host 0.0.0.0 --port 8080
```

#### `info` - Show torrent information

```bash
smirnovtorrent info <torrent|magnet>
```

**Output:**
- Name
- Size
- Files
- Trackers
- Hash
- Piece count and size

**Example:**

```bash
smirnovtorrent info ubuntu-22.04.torrent
```

#### `version` - Show version information

```bash
smirnovtorrent version
```

## Configuration

### Resume Data

Resume data is automatically saved in:
- **Linux/macOS**: `~/.smirnovtorrent/`
- **Windows**: `%APPDATA%\.smirnovtorrent\`

Files:
- `<infohash>.resume` - Resume data (completed pieces, peers, etc.)

### Default Directories

- **Downloads**: Current directory or specified with `-o`
- **Resume data**: `~/.smirnovtorrent/`

## Features

### Peer Discovery

SmirnovTorrent supports multiple peer discovery methods:

1. **Tracker** - Traditional HTTP/UDP trackers
2. **DHT** - Distributed Hash Table (BEP 5)
3. **PEX** - Peer Exchange (BEP 11)

Enable all methods:
```bash
smirnovtorrent download file.torrent --dht --pex
```

### Encryption

Message Stream Encryption (MSE) is enabled by default for compatibility with ISPs that throttle BitTorrent traffic.

Disable encryption:
```bash
smirnovtorrent download file.torrent --encrypt=false
```

### Speed Limits

Limit download and upload speeds:

```bash
# Limit download to 1 MB/s, upload to 512 KB/s
smirnovtorrent download file.torrent \
  --download-limit 1024 \
  --upload-limit 512
```

### Seed Mode

After download completes, SmirnovTorrent automatically switches to seed mode to share with other peers.

Stop seeding and exit:
```bash
# Press Ctrl+C to stop gracefully
# Resume data is automatically saved
```

## Web UI

The Web UI provides a modern interface for managing downloads:

### Features

- Real-time progress updates
- Download/upload speed monitoring
- Peer list with encryption status
- Activity log
- Start/Stop controls

### API Endpoints

```bash
# Get current status
curl http://localhost:8080/api/status

# Start download
curl -X POST http://localhost:8080/api/start

# Stop download
curl -X POST http://localhost:8080/api/stop
```

### Screenshots

Access the Web UI at: `http://localhost:8080`

## Advanced Usage

### Magnet Links

SmirnovTorrent fully supports magnet links (BEP 9):

```bash
smirnovtorrent download "magnet:?xt=urn:btih:HASH&dn=Name"
```

The client will:
1. Connect to DHT to find peers
2. Download metadata from peers
3. Start downloading files

### Multi-file Torrents

For torrents with multiple files:

```bash
# All files downloaded to specified directory
smirnovtorrent download multi-file.torrent -o ~/Downloads
```

Files maintain their directory structure.

### Resume Interrupted Downloads

Resume is automatic. Just run the same command:

```bash
# First run (interrupted)
smirnovtorrent download large-file.torrent
# Press Ctrl+C

# Second run (resumes automatically)
smirnovtorrent download large-file.torrent
```

## Troubleshooting

### Slow Download Speed

1. Enable DHT and PEX for more peers:
   ```bash
   smirnovtorrent download file.torrent --dht --pex
   ```

2. Check firewall settings for port 6881

3. Try with encryption disabled (some peers don't support it):
   ```bash
   smirnovtorrent download file.torrent --encrypt=false
   ```

### Connection Issues

1. Check tracker status
2. Enable DHT as fallback:
   ```bash
   smirnovtorrent download file.torrent --dht
   ```

### Resume Data Issues

Clear resume data and start fresh:
```bash
# Linux/macOS
rm -rf ~/.smirnovtorrent/

# Windows
rmdir /s %APPDATA%\.smirnovtorrent
```

## Performance Tips

1. **Increase peer connections**: Edit `maxPeers` in code (default: 50)
2. **Adjust workers**: Edit `numWorkers` in code (default: 4)
3. **Use SSD**: Faster piece verification and file I/O
4. **Enable all peer sources**: `--dht --pex`

## Support

- **Issues**: https://github.com/ainsonet/smirnovtorrent/issues
- **Discussions**: https://github.com/ainsonet/smirnovtorrent/discussions

## License

MIT License - see LICENSE file for details
