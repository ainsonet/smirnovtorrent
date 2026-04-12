# SmirnovTorrent

<img src="https://raw.githubusercontent.com/ainsonet/smirnovtorrent/master/logo.png" alt="SmirnovTorrent Logo" width="200"/>

🌊 Lightweight BitTorrent client written in Go

## 🚀 Features

- ✅ Parse .torrent files
- ✅ Bencode encoding/decoding
- ✅ Tracker protocol support
- ✅ DHT network support (experimental)
- ✅ Magnet link support (partial)
- ✅ Multi-file torrents
- ✅ Piece verification with SHA-1
- ✅ Rate limiting
- ✅ Encryption (MSE/PE)
- ✅ CLI with progress bar

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

## 🏗️ Project Structure

```
smirnovtorrent/
├── cmd/
│   └── smirnovtorrent/     # CLI application
├── internal/
│   ├── dht/               # DHT network client
│   ├── engine/            # Download engine
│   ├── encryption/        # MSE/PE encryption
│   ├── magnet/            # Magnet link parser
│   ├── parser/            # Torrent file parser
│   ├── peer/              # BitTorrent peer protocol
│   ├── tracker/           # Tracker protocol client
│   └── ...
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
| Parser | ✅ Complete | 6/6 |
| Engine | ✅ Complete | 16/16 |
| Encryption | ✅ Complete | 6/6 |
| Tracker | ✅ Working | ✓ |
| **DHT** | ✅ **Kademlia** | ✓ |
| Magnet | ✅ Parsing | ✓ |
| Peer | ✅ Working | ✓ |
| **Web UI** | ✅ **v0.8.0** | - |

**Total: 28+ tests passing**

**Current version: v0.9.0**

## 📝 Roadmap

- [x] Core torrent parsing
- [x] Tracker protocol
- [x] Peer protocol
- [x] Piece management
- [x] Download engine
- [x] Multi-file support
- [x] Rarest-first algorithm
- [x] Seed mode
- [x] Magnet links (parse)
- [x] DHT bootstrap
- [x] **Kademlia routing table** (v0.9.0)
- [x] BitTorrent encryption
- [x] Rate limiting
- [x] Resume support
- [x] **Web UI** (v0.8.0)
- [ ] Full DHT iterative lookup
- [ ] PEX (Peer Exchange)
- [ ] Magnet metadata download
- [ ] Desktop GUI (Tauri)

## 📄 License

MIT
