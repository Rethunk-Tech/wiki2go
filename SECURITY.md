# Security Policy

## Reporting Security Vulnerabilities

**DO NOT** open a public GitHub issue for security vulnerabilities. Instead, please report them responsibly to:

**Email:** security@rethunk.tech  
**Response SLA:** We aim to respond to security reports within 24 hours.

When reporting a vulnerability, please include:
- Description of the vulnerability
- Affected component(s) and version(s)
- Steps to reproduce (if applicable)
- Potential impact
- Suggested fix (optional)

## Supported Versions

wiki2go is an active development project. Security updates are applied to:

| Version | Support Status | Update Cadence |
|---------|----------------|---|
| Latest | Active | Continuous |

Only the latest version receives security updates. Users are encouraged to upgrade.

## Security Practices

### File Handling Security

- **Path validation** — All file paths validated to prevent directory traversal
  - No `../` sequences allowed
  - Paths restricted to wiki directory and subdirectories
  - Symlink handling is explicit (can be controlled)
- **File permissions** — Respects OS file permissions
- **No arbitrary file deletion** — Operations limited to wiki directory
- **Safe temp files** — Temporary files cleaned up properly

### Markdown Security

- **XSS prevention** — HTML output sanitized
  - Raw HTML in Markdown is escaped (not executed)
  - Event handlers in HTML are stripped
  - Dangerous tags (script, iframe, etc.) are removed
- **Link safety** — External links validated
  - JavaScript URLs blocked (no `javascript:`)
  - File protocol restricted where appropriate
  - Relative links checked for path traversal
- **GFM safety** — GitHub Flavored Markdown parsed safely
  - Tables, code blocks, strikethrough all HTML-escaped
  - Task lists rendered without code execution

### Code Execution Prevention

- **No code execution** — Markdown rendering is text-only
- **Syntax highlighting** — Highlight.js or similar (read-only, no execution)
- **No macro expansion** — Markdown content is never evaluated
- **No template injection** — Content is never passed to template engines

### Server Security

- **No directory listing** — Index pages only; no file browser
- **MIME type validation** — Correct MIME types on served files
- **No symlink traversal** — Symlinks either explicitly allowed or denied
- **Port validation** — Port numbers in valid range (1024-65535)
- **Localhost binding** — Option to bind to localhost only (not default)

### Input Validation

- **CLI arguments** — Validated and sanitized
- **File paths** — Checked for traversal attempts
- **Directory names** — Valid characters only
- **Wiki titles** — No special characters that could cause issues

## Testing & Validation

- **Unit tests** — Path handling, XSS prevention, link validation
- **Integration tests** — Full wiki generation with security test cases
- **Fuzzing** — Malformed Markdown tested against parser
- **Linting** — `go vet` and `golangci-lint` catch common issues
- **Dependency scanning** — `govulncheck` checks for vulnerabilities

## Known Vulnerabilities

None currently known. Reports are welcome via security@rethunk.tech.

## Path Traversal Prevention

All file operations validated:

```go
// ✗ BAD: Not validated
path := filepath.Join(wikiDir, userInput)

// ✓ GOOD: Validated
cleanPath := filepath.Clean(filepath.Join(wikiDir, userInput))
absWikiDir, _ := filepath.Abs(wikiDir)
absPath, _ := filepath.Abs(cleanPath)
if !strings.HasPrefix(absPath, absWikiDir) {
    return fmt.Errorf("path traversal attempt")
}
```

- No `../` sequences allowed
- Absolute paths used for validation
- Symlinks handled explicitly

## XSS Prevention

Markdown to HTML conversion:

```go
// ✓ GOOD: HTML escaped
html := html.EscapeString(markdownContent)

// ✓ GOOD: Dangerous tags removed
html := sanitizer.Sanitize(html, allowlist)

// ✗ BAD: Raw HTML
html := markdownContent // Don't do this
```

- All special characters HTML-escaped
- Dangerous tags removed (script, iframe, event handlers)
- Links validated (no javascript:)

## Dependency Management

wiki2go dependencies:
- **Go standard library** — Core functionality
- **Markdown parser** — Vetted for security
- **Server library** — Go's `net/http` (battle-tested)

**Security checks:**
- `go mod verify` — Verify checksums
- `govulncheck ./...` — Check for vulnerabilities
- **Dependabot** — Automated alerts (if enabled)

## Threat Model

### Attack Vectors Considered

| Vector | Risk | Mitigation |
|--------|------|-----------|
| **Path Traversal** | High | Absolute path validation, no `../` allowed |
| **XSS via Markdown** | High | HTML escaping, dangerous tags removed |
| **File Permissions** | Medium | Respects OS permissions; errors on access denied |
| **Resource Exhaustion** | Medium | File size limits (if implemented), timeouts |
| **Symlink Attacks** | Medium | Explicit symlink handling; default: deny |
| **Denial of Service** | Low | Server timeouts, resource limits |

### Attack Vectors NOT Applicable

- **Remote code execution** — No code execution; text-only processing
- **SQL injection** — No database; file-based only
- **Authentication bypass** — No authentication; local tool only
- **Privilege escalation** — Runs with user permissions only
- **Network attacks** — No network communication (local server only)

## Security Best Practices

### For Users

- **Source directory** — Only include Markdown files from trusted sources
- **File permissions** — Restrict wiki directory to appropriate users
- **Output directory** — Verify output directory permissions before deployment
- **Server binding** — Use `--bind localhost` for local-only wikis
- **Keep updated** — Upgrade to latest version for security patches

### For Developers

- **Path handling** — Always validate file paths; use `filepath.Clean()` + absolute path checks
- **HTML escaping** — Always escape user content before rendering
- **Link validation** — Whitelist safe schemes (`http`, `https`, `file` for local)
- **Testing** — Include security test cases (path traversal, XSS payloads)

## Incident Response

In the event of a confirmed security vulnerability:
1. Impact assessment (what's affected, severity)
2. Fix development (likely path handling or HTML escaping)
3. Testing with security test cases
4. Security update release
5. User notification (if critical)
6. Post-incident review

## Security Checklist

Before generating a wiki for production use:
- [ ] Markdown source directory from trusted source only
- [ ] Output directory has correct permissions (not world-writable)
- [ ] Server bound to localhost or restricted network
- [ ] Test with sample Markdown containing special characters
- [ ] Verify links are correct after generation
- [ ] No symlinks in source directory (or explicitly allowed)

## License

wiki2go is released under the [GNU Affero General Public License (AGPL)](LICENSE). The AGPL requires that modifications distributed in a network service be made available to users.

## Contact

- **Security Issues:** security@rethunk.tech
- **General Support:** support@rethunk.tech
- **Website:** https://rethunk.tech

---

**Last updated:** 2026-04-28
