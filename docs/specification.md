# bumpx Specification

## Overview

**Project Name (working):** bumpx  
**Language:** Go (latest stable release)  
**Purpose:** A single static CLI binary for version management and release chores across polyglot repositories (Rust, Go, Node, Python, .NET, etc.), without requiring Python/Node runtimes.

bumpx is inspired by bumpversion / bump2version and modern successors such as bump-my-version, but is explicitly:

- Language-agnostic
- Polyrepo/monorepo aware
- Static-binary / CI-friendly
- Designed for multi-ecosystem projects (e.g., Cargo + npm + go.mod in one repo)

## Goals

### Single binary, multi-platform
- Build and distribute static binaries for Linux, macOS, and Windows (amd64 + arm64)
- Zero required runtime dependencies (no Python, Node, etc.)

### Multi-ecosystem version bumping
First-class support for:
- `Cargo.toml` (Rust)
- `package.json` (Node/NPM/PNPM/Yarn)
- `go.mod` (Go modules)
- `pyproject.toml` / `setup.cfg` (Python)
- `.csproj` / `AssemblyInfo.cs` (C#/.NET, basic)
- Generic "search & replace with regex" for arbitrary files

### Multiple versioning schemes
- Semantic Versioning (SemVer: MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD])
- Calendar Versioning (CalVer, e.g. YYYY.MM.DD or YYYY.MM.DD.N)
- Optional custom schemes via templates

### Project-config-driven
- Project config file (`.bumpx.toml`) living at repo root
- CLI commands read config and operate deterministically
- Support environment-variable overrides for CI

### Git integration
- Optional:
  - Commit changes with configurable message template
  - Create annotated tags with configurable format
- Dry-run mode to preview diffs without writing or committing

### CI & automation friendly
- Machine-readable output (`--json`)
- Exit codes for typical scenarios:
  - 0: Success
  - 1: Generic failure
  - 2: Config error
  - 3: Version parsing/bumping error
  - 4: File IO or pattern-matching error
  - 5: Git-related error

### Extensibility
- Config-driven file handlers for 80–90% of cases
- External "plugin commands" for advanced/organization-specific workflows

## Non-Goals (for v1)

- Not trying to be a full release orchestrator (no pipelines, no deployment logic)
- Not trying to replace ecosystem-specific advanced tools (e.g., poetry, cargo-release, lerna)
- No native shared-library plugin system

## Directory Structure

```
bumpx/
  cmd/
    bumpx/
      main.go           # CLI entry point

  internal/
    config/             # Config structs, load/validate
    version/            # Version parsing and bumping (SemVer, CalVer)
    files/              # File handlers (cargo, npm, gomod, etc.)
    git/                # Git integration
    hooks/              # External command runner
    output/             # Human and JSON output

  docs/
    specification.md    # This document
    config.md
    cli.md

  .bumpx.toml.example   # Example configuration
  go.mod
  README.md
```

## CLI Commands

### `bumpx init`
Create a starter `.bumpx.toml` in the current directory.

Options:
- `--yes`: non-interactive, accept defaults
- `--scheme`: semver or calver:YYYY.MM.DD

### `bumpx show`
Show the current project version as understood from config.

Options:
- `--json`: output `{ "version": "1.2.3" }`

### `bumpx bump <part>`
Bump the version and update all configured files.

Parts:
- For SemVer: `major`, `minor`, `patch`, `pre`, `build`
- For CalVer: `date`, `build`

Options:
- `--dry-run`: no file writes or git operations
- `--no-git`: ignore git settings in config
- `--no-tag`: do not create a tag
- `--no-commit`: do not commit
- `--message`: override commit message template
- `--tag-name`: override tag name template
- `--json`: machine-readable summary

### `bumpx validate`
Validate `.bumpx.toml` and exit.

Options:
- `--strict`: treat missing paths as error
- `--json`: output validation results in JSON
