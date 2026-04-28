# AGENTS.md — Developer onboarding for wiki2go

wiki2go is a Go CLI tool for creating self-contained wikis from Markdown files with GitHub-Flavored Markdown (GFM) support and static website serving capabilities.

## Quick Navigation

- **For users:** See [README.md](./README.md) — installation, usage, quick start
- **For developers:** This file — building, testing, extending
- **For security:** See [SECURITY.md](./SECURITY.md)

## Project Structure

```
cmd/wiki2go/           Entry point for CLI binary
internal/wiki/         Core wiki generation and serving logic
internal/markdown/     GFM parsing and rendering
internal/server/       Static file server
pkg/                   Public packages
Makefile               Build targets
go.mod                 Module definition
```

## Development Workflow

### Build

```bash
go build -o wiki2go ./cmd/wiki2go

# Or install to $GOPATH/bin
go install ./cmd/wiki2go
```

### Testing

```bash
go test ./...
go test -cover ./...
go test -v ./internal/wiki
```

Tests must pass before any PR merge. Aim for >80% coverage on modified packages.

### Running Locally

```bash
# Create a new wiki
go run ./cmd/wiki2go -- new my-wiki

# Serve locally
go run ./cmd/wiki2go -- serve --open

# Display help
go run ./cmd/wiki2go -- help
```

## Key Modules

### Wiki Generator (`internal/wiki/`)

Converts Markdown files into a static wiki structure:

- Recursively scans directory for `.md` files
- Parses GFM content
- Generates HTML output
- Creates index pages
- Maintains navigation structure

When modifying:
- Add tests in `internal/wiki/*_test.go`
- Test with various Markdown structures
- Verify navigation generation
- Check performance on large wikis (1000+ pages)

### Markdown Renderer (`internal/markdown/`)

Converts GitHub-Flavored Markdown to HTML:

- GFM support (tables, strikethrough, task lists, etc.)
- Code syntax highlighting
- Safe HTML escaping (XSS prevention)
- Link resolution within wiki

When modifying:
- Add test fixtures in `internal/markdown/testdata/`
- Test edge cases (empty content, nested elements)
- Verify XSS prevention
- Document new GFM features added

### Static Server (`internal/server/`)

Serves generated wiki as static website:

- HTTP server on configurable port
- Directory indexing
- MIME type detection
- Browser auto-open support

When modifying:
- Test server startup and shutdown
- Verify error handling (missing files, port conflicts)
- Test on various platforms (Windows/macOS/Linux)

## Adding Features

### New GFM Feature Support

To add support for a new GFM feature (e.g., footnotes):

1. Add parser in `internal/markdown/footnotes.go`
2. Add renderer implementation
3. Add tests in `internal/markdown/gfm_test.go`
4. Update [HUMANS.md](./HUMANS.md) § Markdown Features if user-visible
5. Commit: `feat: add GFM footnotes support`

### New CLI Command

To add a new command (e.g., `wiki2go validate`):

1. Add command handler in `cmd/wiki2go/commands/validate.go`
2. Register command in main CLI setup
3. Add help text and examples
4. Add tests covering the command
5. Update [HUMANS.md](./HUMANS.md) § Usage if user-facing
6. Commit: `feat: add validate command for wiki integrity checking`

### Wiki Configuration

To add config file support (e.g., `wiki.yaml`):

1. Define config struct in `internal/config/config.go`
2. Parse YAML/TOML configuration
3. Merge with CLI flags (CLI takes precedence)
4. Update wiki generation to use config
5. Add config parsing tests
6. Update [HUMANS.md](./HUMANS.md) § Configuration if user-facing
7. Commit: `feat: add wiki.yaml configuration file support`

## Dependency Management

Dependencies are minimal by design:

- **No heavy UI libraries** — Server uses standard `net/http`
- **Standard library** — Core Markdown parsing via `go-text/markdown` or similar lightweight library
- **Minimal external deps** — Focus on core functionality

Before adding a dependency:
1. Justify in commit message (why it's needed)
2. Run `go mod tidy`
3. Check for security issues: `govulncheck ./...`
4. Verify it doesn't bloat the binary (target <20MB)

## Testing Strategy

### Unit Tests

Test individual components:

```bash
go test -v ./internal/wiki       # Wiki generation tests
go test -v ./internal/markdown   # Markdown rendering tests
go test -v ./internal/server     # Server tests
```

### Integration Tests

Test end-to-end wiki creation:

```bash
go test -v ./tests/integration/  # Create wiki, verify output
```

### Test Fixtures

Sample Markdown in `internal/markdown/testdata/`:
- `simple.md` — Basic wiki page
- `gfm-tables.md` — GFM table examples
- `nested-structure.md` — Multi-level wiki
- etc.

Add fixtures for new features:
1. Create `internal/markdown/testdata/<feature>.md`
2. Add to test file
3. Verify expected HTML output

## Benchmarking

For performance-critical paths (Markdown parsing, wiki generation):

```bash
go test -bench ./internal/markdown -benchmem
go test -bench ./internal/wiki -benchmem
```

Target: parse 1000-page wiki in <5 seconds.

## Cross-Platform Considerations

wiki2go supports Windows, macOS, and Linux:

- **File paths** — Use `filepath` package (not hardcoded `/` or `\`)
- **Line endings** — Normalize to `\n` internally
- **File permissions** — Respect platform permissions
- **Browser opening** — Use `os/exec` with platform-specific commands

Test on each platform when:
- Adding file I/O
- Modifying server behavior
- Changing CLI flag handling

## Code Style

- **Go conventions:** Follow `go fmt` and `go vet`
- **Naming:** Clear, descriptive names
- **Comments:** Explain WHY, not WHAT
- **Error handling:** Always check errors
- **Tests:** Table-driven tests for multiple cases

## Performance Targets

- **Binary size** — <20MB
- **Memory usage** — <100MB for 1000-page wiki
- **Wiki generation** — <5 seconds for 1000 pages
- **Server startup** — <100ms
- **Page serving** — <100ms latency

## Debugging

Enable verbose output (if implemented):

```bash
go run ./cmd/wiki2go -- serve --verbose
```

Common issues:
- **Port already in use** — Try different `--port` value
- **File permissions** — Check read permissions on Markdown files
- **Missing GFM features** — Check feature support status

## References

- **[GitHub Flavored Markdown Spec](https://github.github.com/gfm/)** — GFM specification
- **[Go net/http package](https://golang.org/pkg/net/http/)** — HTTP server docs
- **[Go os/exec package](https://golang.org/pkg/os/exec/)** — Command execution docs

---

**Last updated:** 2026-04-28
