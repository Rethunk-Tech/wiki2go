# Wiki2go

## Features

- [ ] Create a self-contained Wiki from Markdown files.
  - [ ] Support for metadata from YAML front matter.
  - [ ] Support for GitHub-Flavoured-Markdown.

- [ ] Serve the Wiki as a static website.

- [ ] Use as a Github Action.

## Quick Start

```bash
# Clone the repository
git clone https://github.com/Rethunk-Tech/wiki2go.git
cd wiki2go

# Build the binary locally
go install ./cmd/wiki2go

# Run the binary
wiki2go --help
```

## Usage

```bash
# Create a new Wiki-based project.
wiki2go new my-wiki

# Open the Wiki's contents in your default editor.
wiki2go edit

# Serve the Wiki. (Open in your default browser.)
wiki2go serve --open
```

## License

This repository's binaries and source code are provided under the:  
[GNU Affero General Public License (AGPL)](LICENSE)
