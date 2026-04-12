# SmirnovTorrent Architecture

## Overview

SmirnovTorrent is a lightweight BitTorrent client written in Go. The project follows a modular architecture with clear separation of concerns.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI (main.go)                       │
│  - Command parsing                                          │
│  - User interface                                           │
│  - Progress display                                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Download Engine (engine)                  │
│  - Orchestrates the download process                        │
│  - Manages peer connections                                 │
│  - Coordinates piece downloading                            │
│  - File assembly                                            │
└─────────────────────────────────────────────────────────────┘
         │                   │                    │
         ▼                   ▼                    ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│   Tracker       │ │     Peer        │ │   Piece         │
│   (tracker)     │ │   (peer)        │ │   Manager       │
│                 │ │                 │ │   (engine)      │
│ - HTTP announce │ │ - Handshake     │ │ - Piece storage │
│ - Peer discovery│ │ - Messages      │ │ - Validation    │
│ - Peer list     │ │ - Connection    │ │ - Progress      │
└─────────────────┘ └─────────────────┘ └─────────────────┘
         │                   │                    │
         ▼                   ▼                    ▼
┌─────────────────────────────────────────────────────────────┐
│                   Parser (internal)                         │
│  - .torrent file parsing                                    │
│  - Bencode serialization/deserialization                    │
│  - Info hash calculation                                    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Bencode (pkg)                            │
│  - Core bencode encoding/decoding                           │
│  - Value types (Int, String, List, Dict)                    │
└─────────────────────────────────────────────────────────────┘
```

## Modules

### `pkg/bencode`
Core bencode serialization library.
- **Responsibility**: Encode/decode bencode format
- **Key types**: `Int`, `String`, `List`, `Dict`
- **Key functions**: `Marshal`, `Unmarshal`

### `internal/parser`
Torrent file parsing.
- **Responsibility**: Parse .torrent files into usable structures
- **Key types**: `Torrent`, `TorrentInfo`, `FileInfo`
- **Key functions**: `Parse`, `ParseFile`, `CalculateInfoHash`

### `internal/tracker`
Tracker client implementation.
- **Responsibility**: Communicate with BitTorrent trackers
- **Key types**: `Tracker`, `AnnounceParams`, `AnnounceResponse`, `PeerInfo`
- **Key functions**: `Announce`, `GetPeers`

### `internal/peer`
Peer protocol implementation.
- **Responsibility**: Handle peer connections and messaging
- **Key types**: `Peer`, `PeerConnection`
- **Key functions**: `SendHandshake`, `ReadHandshake`, `SendMessage`, `ReadMessage`
- **Message types**: Choke, Unchoke, Interested, Have, Bitfield, Request, Piece, Cancel

### `internal/engine`
Download orchestration.
- **Responsibility**: Coordinate the download process
- **Key types**: `DownloadEngine`, `PieceManager`, `DownloadStatus`, `SeedManager`
- **Key functions**: `Start`, `GetNextPiece`, `MarkPieceComplete`, `AssembleFile`
- **Modes**: Download mode, Seed mode

## Data Flow

1. **Torrent Loading**
   ```
   CLI → parser.ParseFile() → Torrent struct
   ```

2. **Peer Discovery**
   ```
   DownloadEngine → tracker.Announce() → []PeerInfo
   ```

3. **Peer Connection**
   ```
   DownloadEngine → peer.Connect() → PeerConnection
   PeerConnection → SendHandshake() → Handshake exchange
   ```

4. **Piece Download**
   ```
   PieceManager.GetNextPiece() → Piece
   PeerConnection.SendRequest() → Request piece
   PeerConnection.ReadMessage() → Piece data
   PieceManager.MarkPieceComplete() → Validate & store
   ```

5. **File Assembly**
   ```
   PieceManager.AssembleFile() → Raw data
   DownloadEngine.assembleFiles() → Write to disk
   ```

6. **Seed Mode**
   ```
   Download complete → SeedManager.StartSeed()
   Incoming peer connections → HandlePeerConnection()
   Piece requests → handleRequest() → SendPiece()
   ```

## Concurrency

The engine uses mutex-based synchronization:
- `PieceManager` uses `sync.RWMutex` for piece data
- `DownloadEngine` uses separate mutexes for peers and status
- Peer connections run in goroutines

## Future Enhancements

1. **Parallel Downloads**: Download multiple pieces simultaneously from different peers
2. **Piece Prioritization**: Priority-based piece selection (rarest first)
3. **Seed Mode**: Upload pieces after download completion
4. **DHT Support**: Decentralized peer discovery
5. **Magnet Links**: Direct support for magnet URIs

## Error Handling

All modules return typed errors with context:
```go
func Download(source string) error {
    // Wrap low-level errors with context
    if err := tracker.Announce(); err != nil {
        return fmt.Errorf("tracker announce failed: %w", err)
    }
    return nil
}
```

## Testing

Each module has comprehensive tests:
- Unit tests for core functionality
- Integration tests for module interactions
- Mock servers for network testing

Run tests with:
```bash
go test ./...
```