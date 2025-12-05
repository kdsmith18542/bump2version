# bumpx

[![Go CI](https://github.com/kdsmith18542/bump2version/workflows/CI/badge.svg)](https://github.com/kdsmith18542/bump2version/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/kdsmith18542/bumpx.svg)](https://pkg.go.dev/github.com/kdsmith18542/bumpx)

A modern, fast version management CLI tool written in Go. Designed for polyglot repositories with support for multiple versioning schemes and file formats.

> **Note**: This is a complete rewrite of the original Python bump2version project. The Python version has been archived in `archive/python-legacy/` for reference.

## Features

- **Single static binary** - No runtime dependencies
- **Polyglot support** - Works with Cargo.toml, package.json, go.mod, pyproject.toml, and more
- **Multiple versioning schemes** - SemVer and CalVer support
- **Git integration** - Automatic commits and tags
- **Dry-run mode** - Preview changes before applying
- **JSON output** - Machine-readable output for CI/CD pipelines
- **Hooks** - Run custom scripts before/after version bumps

## Installation

### From Source

```bash
go install github.com/kdsmith18542/bumpx/cmd/bumpx@latest
```

### From Binary

Download the latest release from the [releases page](https://github.com/kdsmith18542/bump2version/releases).

## Quick Start

### Initialize a new project

```bash
bumpx init
```

This creates a `.bumpx.toml` configuration file with auto-detected project files.

### Show current version

```bash
bumpx show
```

### Bump version

```bash
# Bump major version (1.0.0 -> 2.0.0)
bumpx bump major

# Bump minor version (1.0.0 -> 1.1.0)
bumpx bump minor

# Bump patch version (1.0.0 -> 1.0.1)
bumpx bump patch

# Dry-run (preview changes without applying)
bumpx bump minor --dry-run
```

### Validate configuration

```bash
bumpx validate
```

## Configuration

bumpx uses a `.bumpx.toml` configuration file:

```toml
[project]
name = "my-project"
current_version = "1.2.3"
version_scheme = "semver"  # "semver" | "calver:YYYY.MM.DD"

[files]

# Rust project
[files.cargo]
paths = ["Cargo.toml"]
type = "cargo"

# Node.js project
[files.npm]
paths = ["package.json"]
type = "npm"

# Go module
[files.gomod]
paths = ["go.mod"]
type = "gomod"

# Python project
[files.python]
paths = ["pyproject.toml"]
type = "pyproject"

# Generic regex-based replacement
[files.version_txt]
paths = ["VERSION.txt"]
type = "regex"
pattern = '^version=(?P<version>[0-9A-Za-z\.\-\+]+)$'
replacement = 'version={version}'

# Git settings
[git]
enable = true
commit = true
commit_message = "chore(release): {version}"
tag = true
tag_name = "v{version}"
tag_message = "Release {version}"
allow_dirty = false

# Hooks (optional)
[hooks]
pre_bump = ""   # e.g., "scripts/pre_bump.sh"
post_bump = ""  # e.g., "scripts/post_bump.sh"
```

## Supported File Types

| Type | Files | Description |
|------|-------|-------------|
| `cargo` | `Cargo.toml` | Rust package manifest |
| `npm` | `package.json` | Node.js package manifest |
| `gomod` | `go.mod` | Go module file |
| `pyproject` | `pyproject.toml` | Python project file |
| `dotnet` | `*.csproj`, `AssemblyInfo.cs` | .NET project files |
| `regex` | Any | Custom regex-based replacement |

## Version Schemes

### SemVer (Semantic Versioning)

Format: `MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]`

Available parts: `major`, `minor`, `patch`, `pre`, `build`

```bash
bumpx bump major  # 1.2.3 -> 2.0.0
bumpx bump minor  # 1.2.3 -> 1.3.0
bumpx bump patch  # 1.2.3 -> 1.2.4
bumpx bump pre    # 1.2.3 -> 1.2.3-alpha.1
bumpx bump build  # 1.2.3 -> 1.2.3+1
```

### CalVer (Calendar Versioning)

Format: `YYYY.MM.DD[.RELEASE]`

Available parts: `date`, `build` (or `release`)

```toml
[project]
version_scheme = "calver:YYYY.MM.DD"
```

```bash
bumpx bump date   # Updates to current date
bumpx bump build  # 2024.01.15 -> 2024.01.15.1
```

## Command Line Options

### Global Options

```
-c, --config string   Config file (default ".bumpx.toml")
    --json            Output in JSON format
-h, --help            Help for bumpx
-v, --version         Version for bumpx
```

### Bump Command Options

```
-n, --dry-run         Don't write any files, just pretend
    --no-git          Ignore git settings in config
    --no-tag          Do not create a tag
    --no-commit       Do not commit
-m, --message string  Override commit message
    --tag-name string Override tag name
```

### Init Command Options

```
-y, --yes             Non-interactive, accept defaults
    --scheme string   Version scheme (default "semver")
```

### Validate Command Options

```
    --strict          Treat missing paths as errors
```

## Development

### Building

```bash
go build -o bumpx ./cmd/bumpx/
```

### Testing

```bash
go test -v ./...
```

## Migration from bump2version

If you're migrating from the Python bump2version, you'll need to:

1. Create a new `.bumpx.toml` configuration file (use `bumpx init`)
2. Map your existing `.bumpversion.cfg` settings to the new TOML format

The Python version is archived in `archive/python-legacy/` for reference during migration.

## License

bumpx is licensed under the MIT License - see the [LICENSE.rst](LICENSE.rst) file for details.

## Related Projects

See [RELATED.md](RELATED.md) for other version bumping tools
