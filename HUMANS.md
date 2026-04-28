# HUMANS.md — User guide for wiki2go

This file covers **installation, usage, and examples**. For **developer context**, see [AGENTS.md](./AGENTS.md).

## Installation

### From Releases

Download pre-built binary:

```bash
wget https://github.com/Rethunk-Tech/wiki2go/releases/download/v1.0.0/wiki2go-linux-amd64
chmod +x wiki2go-linux-amd64
sudo mv wiki2go-linux-amd64 /usr/local/bin/wiki2go
```

### Build from Source

```bash
git clone https://github.com/Rethunk-Tech/wiki2go.git
cd wiki2go
go build -o wiki2go ./cmd/wiki2go
sudo cp wiki2go /usr/local/bin/
```

### Verify Installation

```bash
wiki2go --version
wiki2go --help
```

## Quick Start

### Create a New Wiki

```bash
wiki2go new my-wiki
cd my-wiki
```

This creates a new wiki structure:

```
my-wiki/
├── README.md
├── config.yaml
└── docs/
    └── index.md
```

### Edit Wiki Content

Open wiki in your editor:

```bash
wiki2go edit
```

Or manually edit Markdown files in `docs/`.

### Serve the Wiki

Start development server:

```bash
wiki2go serve --open
```

Options:
- `--port 8080` — Change port (default: 3000)
- `--open` — Auto-open in browser
- `--watch` — Auto-reload on file changes

Visit `http://localhost:3000` in your browser.

## Usage

### Basic Commands

#### Create Wiki

```bash
wiki2go new project-wiki
```

#### List Content

```bash
wiki2go list
```

Shows all pages and structure.

#### Validate Wiki

```bash
wiki2go validate
```

Checks for:
- Broken links
- Missing files
- Invalid Markdown

#### Generate Static Site

```bash
wiki2go build -o public/
```

Generates static HTML in `public/` directory.

### Directory Structure

Standard wiki layout:

```
wiki/
├── README.md           # Wiki home
├── config.yaml         # Settings
├── docs/
│   ├── index.md        # Getting started
│   ├── guide/
│   │   ├── intro.md
│   │   └── setup.md
│   └── api/
│       ├── overview.md
│       └── endpoints.md
└── assets/
    └── images/
```

### Front Matter

Add metadata to Markdown files:

```markdown
---
title: "Getting Started"
author: "Rethunk Tech"
date: 2026-04-28
---

# Getting Started

Content here...
```

Supported fields:
- `title` — Page title
- `author` — Page author
- `date` — Publication date
- `tags` — Comma-separated tags
- `draft` — true/false (exclude from output if true)

## Configuration

Edit `config.yaml`:

```yaml
site:
  title: "My Project Wiki"
  description: "Project documentation"
  url: "https://wiki.example.com"

build:
  output: public/
  theme: default

markdown:
  smartTypography: true
  syntax-highlight: true

server:
  port: 3000
  host: localhost
  auto-reload: true
```

## Examples

### Project Documentation

Create company/project documentation:

```bash
wiki2go new acme-corp-wiki
cd acme-corp-wiki

# Create structure
mkdir -p docs/{getting-started,api,deployment,troubleshooting}

# Add pages
echo "# Getting Started" > docs/getting-started/index.md
echo "# API Reference" > docs/api/index.md
```

### Knowledge Base

Build internal knowledge base:

```
docs/
├── engineering/
│   ├── standards.md
│   ├── code-review.md
│   └── testing.md
├── operations/
│   ├── runbooks/
│   │   ├── deployment.md
│   │   └── incident-response.md
│   └── monitoring.md
└── process/
    ├── onboarding.md
    └── communication.md
```

### Team Handbook

Create employee handbook:

```bash
wiki2go new company-handbook

# Handbook sections
mkdir -p docs/{policies,benefits,office,communication}
echo "# Time Off Policy" > docs/policies/time-off.md
echo "# Health Insurance" > docs/benefits/health.md
```

## Markdown Features

### GitHub Flavored Markdown

Full GFM support:

```markdown
# Tables

| Header 1 | Header 2 |
|----------|----------|
| Cell 1   | Cell 2   |

# Strikethrough

~~This is struck through~~

# Task Lists

- [x] Completed task
- [ ] Pending task
```

### Code Blocks with Syntax Highlighting

```markdown
\`\`\`javascript
function hello() {
  console.log("Hello, world!");
}
\`\`\`
```

### Links

Internal links:

```markdown
[See Getting Started](getting-started/index.md)
```

External links:

```markdown
[Visit Rethunk](https://rethunk.tech)
```

## Building & Deployment

### Generate Static Site

```bash
wiki2go build -o public/
```

Output is static HTML in `public/`.

### Deploy to GitHub Pages

```bash
# Build site
wiki2go build -o docs/

# Commit and push
git add docs/
git commit -m "docs: rebuild wiki"
git push

# Enable GitHub Pages in repo settings
# Source: gh-pages or docs/ folder
```

### Deploy to Web Server

```bash
# Build locally
wiki2go build -o dist/

# Upload to server
scp -r dist/* user@server:/var/www/wiki/

# Or use CI/CD pipeline
```

### Docker Deployment

```dockerfile
FROM golang:1.25 as builder
WORKDIR /app
COPY . .
RUN go build -o wiki2go ./cmd/wiki2go

FROM scratch
COPY --from=builder /app/wiki2go /
COPY --from=builder /app/docs /docs
CMD ["/wiki2go", "build", "-o", "/public/"]
```

Build and run:

```bash
docker build -t wiki2go .
docker run -v $(pwd)/docs:/docs -v $(pwd)/public:/public wiki2go
```

## Development Workflow

### Edit and Preview

```bash
# Terminal 1: Live server
wiki2go serve --watch

# Terminal 2: Edit files in editor
# Server auto-reloads on save
```

### Check for Issues

```bash
wiki2go validate

# Output:
# ✓ All links valid
# ✗ Broken link: docs/api/users.md (referenced but not found)
# ✗ Unlinked page: docs/archive.md (not referenced anywhere)
```

### Build for Production

```bash
wiki2go build --minify -o public/
```

Options:
- `--minify` — Minify HTML/CSS
- `-o` — Output directory

## Troubleshooting

### Server Won't Start

**Symptom:** "Port already in use"

**Solution:**
```bash
wiki2go serve --port 8080  # Use different port
```

### Links Broken After Build

**Symptom:** Links work in development, broken after `build`

**Solutions:**
1. Use relative links: `[Link](../other/page.md)`
2. Validate before building: `wiki2go validate`
3. Check output structure: `ls -la public/`

### Markdown Not Rendering

**Symptom:** Raw Markdown appears in output

**Solutions:**
1. Verify file extension is `.md`
2. Check front matter syntax (if used)
3. Run `wiki2go validate`

### Missing Images

**Symptom:** Images not showing after build

**Solutions:**
1. Copy assets: `cp -r assets public/`
2. Use relative paths: `![Alt text](../assets/image.png)`
3. Check image format is supported (PNG, JPG, GIF, SVG)

## Best Practices

- **Organize hierarchically:** Use directory structure to organize topics
- **Use front matter:** Add metadata for search/filtering
- **Link internally:** Cross-link related pages
- **Keep updated:** Schedule regular content reviews
- **Validate regularly:** Run `wiki2go validate` before deploying
- **Use drafts:** Mark incomplete pages with `draft: true`

## Tips

### Navigation

Use `index.md` files for directory landing pages:

```
docs/
├── api/
│   ├── index.md      ← User lands here when clicking API
│   ├── users.md
│   └── posts.md
```

### Table of Contents

Create manual TOC in README:

```markdown
# My Wiki

- [Getting Started](getting-started.md)
- [API Reference](api/index.md)
  - [Users](api/users.md)
  - [Posts](api/posts.md)
```

### Search

Generated sites may include search (depending on theme).

Enable search in config:

```yaml
features:
  search: true
```

## Support

For issues:
1. Check `wiki2go validate` output
2. Review Markdown syntax
3. Check config.yaml
4. File issue with: wiki2go version, example content, error message

---

**Last updated:** 2026-04-28
