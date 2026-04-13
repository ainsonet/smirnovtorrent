# SmirnovTorrent

<img src="https://raw.githubusercontent.com/ainsonet/smirnovtorrent/master/logo.png" alt="SmirnovTorrent Logo" width="200"/>

🌊 Lightweight BitTorrent client written in Go

## 🚀 Features

- ✅ Parse .torrent files
- ✅ Bencode encoding/decoding
- ✅ Tracker protocol support
- ✅ DHT network support (BEP 5)
- ✅ Magnet link support (BEP 9)
- ✅ Multi-file torrents
- ✅ Piece verification with SHA-1
- ✅ Rate limiting (token bucket)
- ✅ Encryption (MSE/PE - BEP 47)
- ✅ CLI with progress bar
- ✅ Web UI (real-time)
- ✅ Desktop GUI (Tauri) ✨ NEW
- ✅ PEX (Peer Exchange - BEP 11)
- ✅ Resume support
- ✅ Structured logging
- ✅ JSON configuration

## 📦 Installation

```bash
git clone https://github.com/ainsonet/smirnovtorrent.git
cd smirnovtorrent
go build -o smirnovtorrent.exe ./cmd/smirnovtorrent
```

## 💻 Usage

```bash
# Show torrent information
smirnovtorrent info example.torrent

# Download a torrent
smirnovtorrent download example.torrent

# Download from magnet link (experimental)
smirnovtorrent download "magnet:?xt=urn:btih:..."

# Start Web UI (default port 8080)
smirnovtorrent webui [port]

# Show version
smirnovtorrent version

# Show help
smirnovtorrent help
```

### Web Interface

Open http://localhost:8080 in your browser to access the Web UI with:
- Real-time progress monitoring
- Download/upload statistics
- Active peers count
- Start/stop controls
- Activity log

### Desktop GUI (NEW!)

Launch the desktop application for a native experience:

```bash
cd gui
npm install
npm run tauri dev
```

Features:
- 🎨 Modern dark theme UI
- 📊 Real-time statistics
- 🎯 Visual progress bars
- ⏯️ Pause/resume controls
- 🔍 File browser integration
- 📝 Activity log

See [gui/README.md](gui/README.md) for detailed setup instructions.

## 🏗️ Project Structure

```
smirnovtorrent/
├── cmd/
│   └── smirnovtorrent/     # CLI application
├── gui/                    # Desktop GUI (Tauri) ✨
│   ├── src-tauri/         # Rust backend
│   ├── public/            # Frontend assets
│   └── README.md          # GUI documentation
├── internal/
│   ├── dht/               # DHT network client
│   ├── engine/            # Download engine
│   ├── encryption/        # MSE/PE encryption
│   ├── magnet/            # Magnet link parser
│   ├── parser/            # Torrent file parser
│   ├── peer/              # BitTorrent peer protocol
│   ├── tracker/           # Tracker protocol client
│   ├── config/            # Configuration management
│   ├── logger/            # Structured logging
│   └── ratelimit/         # Rate limiting
└── pkg/
    └── bencode/           # Public bencode package
```

## 🧪 Running Tests

```bash
# Run all unit tests
go test ./...

# Run with coverage
go test ./... -cover

# Run E2E tests (requires TORRENT_FILE environment variable)
go test -tags e2e -v -timeout 5m ./cmd/smirnovtorrent
```

## 📈 Development Status

| Module | Status | Tests |
|--------|--------|-------|
| Parser | ✅ Complete | 6/6 + 2 bench |
| Engine | ✅ Complete | 24/24 |
| Encryption | ✅ Full Integration | 6/6 |
| Tracker | ✅ Working | ✓ |
| DHT | ✅ Kademlia + Iterative | ✓ |
| Magnet | ✅ Metadata (BEP 9) | ✓ |
| PEX | ✅ BEP 11 | ✓ |
| Peer | ✅ Working | ✓ |
| Web UI | ✅ v0.8.0 | - |
| **Desktop GUI** | ✅ **v1.0.0 (Tauri)** | - |
| **Resume** | ✅ **Graceful Shutdown** | ✓ |
| **Config** | ✅ **Integrated** | 5/5 |
| **Logger** | ✅ **Structured** | 4/4 |
| **RateLimit** | ✅ **Token Bucket** | 8/8 + 2 bench |

**Total: 50+ tests passing** (45 unit + 5 E2E)

**Current version: v1.0.0** 🎉

## 📝 Roadmap

- [x] Core torrent parsing
- [x] Tracker protocol
- [x] Peer protocol
- [x] Piece management
- [x] Download engine
- [x] Multi-file support
- [x] Rarest-first algorithm
- [x] Seed mode
- [x] Magnet links (BEP 9)
- [x] DHT bootstrap (BEP 5)
- [x] Kademlia routing table
- [x] DHT iterative lookup
- [x] PEX (Peer Exchange - BEP 11)
- [x] BitTorrent encryption (BEP 47)
- [x] Graceful shutdown & resume
- [x] Rate limiting (token bucket)
- [x] Configuration system
- [x] Structured logging
- [x] Web UI
- [x] **Desktop GUI (Tauri)** ✨
- [x] **E2E testing**
- [x] **CI/CD pipeline**
- [ ] Advanced statistics dashboard
- [ ] WebTorrent support
- [ ] IPv6 support
- [ ] Plugin system

## 📄 License

MIT
