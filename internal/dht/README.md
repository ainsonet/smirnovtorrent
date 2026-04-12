# DHT Module Documentation

## Overview

The DHT (Distributed Hash Table) module provides decentralized peer discovery for BitTorrent without relying on trackers. It implements a simplified version of the Kademlia DHT protocol.

## Architecture

```
┌─────────────────────────────────────────────┐
│              DHTClient                       │
│  - UDP listener                              │
│  - Node ID management                        │
│  - Request/response handling                 │
└─────────────────────────────────────────────┘
                    │
         ┌──────────┴──────────┐
         │                     │
         ▼                     ▼
┌─────────────────┐   ┌─────────────────┐
│  KademliaTable  │   │  Bootstrap Nodes │
│  - Node storage │   │  - Initial peers │
│  - Routing      │   │  - Discovery     │
└─────────────────┘   └─────────────────┘
```

## Usage

### Basic Example

```go
package main

import (
    "log"
    "smirnovtorrent/internal/dht"
)

func main() {
    // Bootstrap nodes (public DHT routers)
    bootstrap := []string{
        "router.bittorrent.com:6881",
        "dht.transmissionbt.com:6881",
        "router.utorrent.com:6881",
    }

    // Create DHT client on random port
    client, err := dht.NewDHTClient(bootstrap, 0)
    if err != nil {
        log.Fatal(err)
    }

    // Start DHT
    if err := client.Start(); err != nil {
        log.Fatal(err)
    }
    defer client.Stop()

    // Find peers for a torrent
    peers, err := client.FindPeer("73ef7ed9f70e94f1e3a4b8b5c2d1e0f9a8b7c6d5")
    if err != nil {
        log.Printf("Peer discovery failed: %v", err)
    }

    log.Printf("Found %d peers", len(peers))
}
```

### Advanced Usage

```go
// Get continuous peer updates
go func() {
    for peers := range client.GetPeersFound() {
        log.Printf("New peers: %v", peers)
        // Add peers to download engine
    }
}()

// Check node count
nodeCount := client.GetNodeCount()
log.Printf("Connected to %d DHT nodes", nodeCount)
```

## API Reference

### DHTClient

Main interface for DHT operations.

#### NewDHTClient

```go
func NewDHTClient(bootstrap []string, port uint16) (*DHTClient, error)
```

Creates a new DHT client.

- `bootstrap`: List of initial DHT nodes to connect to
- `port`: Local UDP port (0 for random)

#### Start

```go
func (d *DHTClient) Start() error
```

Starts the DHT client and begins listening for responses.

#### Stop

```go
func (d *DHTClient) Stop()
```

Stops the DHT client and closes all connections.

#### FindPeer

```go
func (d *DHTClient) FindPeer(infoHash string) ([]string, error)
```

Finds peers sharing a torrent with the given info hash.

- Returns list of peer addresses in `ip:port` format
- Blocks until peers are found or timeout (5 seconds)

#### GetPeersFound

```go
func (d *DHTClient) GetPeersFound() <-chan []string
```

Returns a channel for continuous peer discovery updates.

#### GetNodeCount

```go
func (d *DHTClient) GetNodeCount() int
```

Returns the number of nodes in the Kademlia table.

### Peer Utilities

#### ParsePeerAddress

```go
func ParsePeerAddress(peerStr string) (string, uint16, error)
```

Parses a peer address string into IP and port.

```go
ip, port, err := dht.ParsePeerAddress("192.168.1.100:6881")
// ip = "192.168.1.100", port = 6881
```

#### PeerToString

```go
func PeerToString(ip string, port uint16) string
```

Converts IP and port to address string.

```go
addr := dht.PeerToString("10.0.0.1", 6882)
// addr = "10.0.0.1:6882"
```

#### EncodePeerInfo / DecodePeerInfo

```go
func EncodePeerInfo(ip string, port uint16) ([]byte, error)
func DecodePeerInfo(data []byte) (string, uint16, error)
```

Binary encoding/decoding for DHT peer information.

## Protocol Details

### Message Format

DHT messages use Bencode encoding:

```
Request:
d8:ping14:0200000000000000000014:aae

Response:
d8:ping14:020000000000000000014:aae
```

### Node ID

- 20-byte random identifier
- Used for routing in Kademlia tree
- Generated on client startup

### UDP Protocol

- Default port: 6881
- MTU: 1400 bytes
- Timeout: 5 seconds per request

## Bootstrap Nodes

Common public DHT routers:

```
- router.bittorrent.com:6881
- dht.transmissionbt.com:6881
- router.utorrent.com:6881
- dht.libtorrent.org:25401
- dht.aelitis.com:6881
```

## Integration with DownloadEngine

```go
// Create DHT client
dhtClient, _ := dht.NewDHTClient(bootstrap, 0)
dhtClient.Start()
defer dhtClient.Stop()

// Use in download engine
peers, _ := dhtClient.FindPeer(torrent.InfoHash)
for _, peerAddr := range peers {
    ip, port, _ := dht.ParsePeerAddress(peerAddr)
    // Connect to peer
}
```

## Limitations

Current implementation is simplified:

- ✅ Basic DHT node discovery
- ✅ Peer lookup by info hash
- ✅ UDP communication
- ⏳ Full Kademlia routing table
- ⏳ NAT traversal (STUN/TURN)
- ⏳ Message signing
- ⏳ IPv6 support

## Security Considerations

- Validate peer addresses before connecting
- Rate limit peer discovery requests
- Handle malicious DHT responses
- Consider using trusted bootstrap nodes

## Testing

```bash
go test ./internal/dht -v
```

## Future Enhancements

1. Complete Kademlia algorithm implementation
2. Peer caching and eviction policies
3. Integration with main download loop
4. Support for extended DHT protocol (BEP 5)
5. NAT hole-punching support