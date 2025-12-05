# bumpx - Version Management for Polyglot Repositories

A single static CLI binary for version management and release chores across polyglot repositories (Rust, Go, Node, Python, .NET, etc.), inspired by bumpversion and bump-my-version.

## Features

- **Language-agnostic**: Supports multiple ecosystems (Cargo, npm, Go modules, Python, .NET, custom regex)
- **Static binary**: No runtime dependencies - single executable
- **Config-driven**: TOML configuration with flexible file handlers
- **CI-friendly**: Designed for automated release workflows
- **Git integration**: Optional commit and tag creation
- **Extensible**: Hooks for pre/post bump operations

## Quick Start

### Installation

Download from [GitHub Releases](https://github.com/kdsmith18542/bumpx/releases) or build from source:

```bash
go build -ldflags='-s -w' -o bumpx ./cmd/bumpx/
```

### Basic Usage

1. **Initialize a project:**
   ```bash
   bumpx init --yes
   ```

2. **Show current version:**
   ```bash
   bumpx show
   ```

3. **Bump version:**
   ```bash
   bumpx bump patch
   bumpx bump minor --git-commit --git-tag
   ```

4. **Validate configuration:**
   ```bash
   bumpx validate
   ```

## Configuration

Create a `.bumpx.toml` file in your project root:

```toml
[project]
name = "my-project"
current_version = "1.2.3"
version_scheme = "semver"

[files.npm]
paths = ["package.json"]
type = "npm"

[files.cargo]
paths = ["Cargo.toml"]
type = "cargo"

[files.gomod]
paths = ["go.mod"]
type = "gomod"

[git]
enable = true
commit = true
tag = true
```

See `.bumpx.toml.example` for a complete configuration reference.

## Documentation

- [Quickstart Guide](docs/quickstart.md)
- [Configuration Reference](docs/config.md)
- [CLI Reference](docs/cli.md)
- [Specification](docs/specification.md)
- [Contributing](CONTRIBUTING.md)

## License

MIT