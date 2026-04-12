# Contributing to SmirnovTorrent

Thank you for your interest in contributing to SmirnovTorrent! This document provides guidelines and information for contributors.

## Development Setup

### Prerequisites

- Go 1.22 or higher
- Git
- A text editor or IDE (VS Code, GoLand, etc.)

### Building

```bash
# Build the binary
go build -o smirnovtorrent.exe ./cmd/smirnovtorrent

# Run tests
go test ./...

# Run with verbose output
go test -v ./...
```

## Code Style

We follow standard Go conventions with some project-specific preferences:

### Naming Conventions

- Use descriptive variable names (`pieceManager` not `pm`)
- Exported functions should have clear documentation
- Avoid abbreviations unless well-known (e.g., `ID` is OK, `Cnt` is not)

### Comments

```go
// Good: Explain WHY, not WHAT
// Retry with exponential backoff to handle transient network failures
func retryWithBackoff() { ... }

// Bad: Obvious from the code
// Increment counter by 1
counter++
```

### Error Handling

```go
// Wrap errors with context
if err := tracker.Announce(); err != nil {
    return fmt.Errorf("tracker announce failed: %w", err)
}

// Don't ignore errors
if err := os.WriteFile(path, data, 0644); err != nil {
    log.Printf("Warning: failed to write file: %v", err)
}
```

### Concurrency

- Use `sync.RWMutex` when reads outnumber writes
- Always defer unlocks: `defer mu.Unlock()`
- Avoid holding locks during I/O operations

## Architecture Guidelines

### Module Responsibilities

- **pkg/bencode**: Core encoding/decoding only
- **internal/parser**: Parse torrent files, no network
- **internal/tracker**: HTTP communication with trackers
- **internal/peer**: Peer protocol, no business logic
- **internal/engine**: Orchestration and coordination

### Dependencies

- Internal packages should not import external dependencies directly
- Use interfaces for dependencies when possible
- Keep circular dependencies out

## Testing

### Unit Tests

```go
func TestParseMinimalTorrent(t *testing.T) {
    data := createTestTorrent()
    torrent, err := Parse(data)
    if err != nil {
        t.Fatalf("Parse failed: %v", err)
    }
    
    if torrent.Announce != "http://test.com" {
        t.Errorf("Expected announce 'http://test.com', got '%s'", torrent.Announce)
    }
}
```

### Test Organization

- One test file per source file (`foo_test.go` for `foo.go`)
- Table-driven tests for multiple cases
- Helper functions in `test_helpers.go`

## Pull Request Process

1. **Fork** the repository
2. **Create a branch** from `main`
3. **Make changes** following the guidelines above
4. **Write tests** for new functionality
5. **Update documentation** if needed
6. **Submit PR** with a clear description

### Commit Messages

Follow conventional commits:

```
feat: Add new feature
fix: Fix bug
docs: Update documentation
refactor: Code restructuring
test: Add tests
chore: Maintenance tasks
```

Example:
```
feat: Add Rarest-first piece selection algorithm

- Implement piece rarity tracking
- Add peer capability discovery
- Improve download efficiency by 30%
- Update tests and documentation
```

## Code Review

- All PRs require at least one reviewer
- Address all comments before merging
- Keep PRs focused and manageable (< 500 lines ideal)
- Link related issues in PR description

## Issues

### Reporting Bugs

Include:
- Expected behavior
- Actual behavior
- Steps to reproduce
- Go version and OS
- Any relevant logs

### Feature Requests

Include:
- Problem being solved
- Proposed solution
- Alternative approaches considered
- Use cases

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Questions?

Open an issue for general questions or reach out to the maintainers.

---

Happy coding! 🚀