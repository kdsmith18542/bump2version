# bumpx - Version Management for Polyglot Repositories

A single static CLI binary for version management and release chores across polyglot repositories.

## Quick Start

### Installation

Download the latest release from [GitHub Releases](https://github.com/kdsmith18542/bumpx/releases) or build from source:

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

### Configuration

Create a `.bumpx.toml` file:

```toml
[project]
name = "my-project"
current_version = "1.2.3"
version_scheme = "semver"

[files]
paths = ["package.json", "Cargo.toml", "go.mod"]

[git]
commit = false
tag = false
```

### Supported File Types

- **npm**: `package.json`
- **cargo**: `Cargo.toml`
- **gomod**: `go.mod`
- **python**: `pyproject.toml`, `setup.cfg`
- **dotnet**: `.csproj`, `AssemblyInfo.cs`
- **regex**: Custom patterns for any file

### Version Schemes

- **semver**: `MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]`
- **calver**: Calendar versioning (e.g., `YYYY.MM.DD`)
- **custom**: Template-based schemes

### CLI Reference

```bash
bumpx --help
bumpx show --format=json
bumpx bump patch --dry-run
bumpx bump minor --git-commit --git-tag
bumpx validate
bumpx init --yes
```

### Exit Codes

- `0`: Success
- `1`: Usage error
- `2`: Operation failed

### CI/CD Integration

```yaml
# GitHub Actions example
- name: Bump version
  run: |
    bumpx bump patch --git-commit --git-tag
    version=$(bumpx show --format=json | jq -r .version)
    echo "VERSION=$version" >> $GITHUB_ENV
```

## Development

```bash
# Run tests
go test ./...

# Build release binary
make release

# Run integration tests
go test ./cmd/bumpx/...
```

## License

See LICENSE.rst