# Command Line Interface

## Commands

### bumpx show

Show the current project version as understood from config.

**Usage:**
```bash
bumpx show
```

**Options:**
- `--json`: Output in JSON format
- `-c, --config`: Config file (default is .bumpx.toml)
- `-h, --help`: Help for show command

### bumpx bump <part>

Bump the version and update all configured files.

**Usage:**
```bash
bumpx bump <part>
```

For SemVer: major, minor, patch, pre, build
For CalVer: date, build

**Options:**
- `-n, --dry-run`: Don't write any files, just pretend
- `--no-git`: Ignore git settings in config
- `--git-commit`: Create a git commit (overrides config)
- `--git-tag`: Create a git tag (overrides config)
- `--no-tag`: Do not create a tag
- `--no-commit`: Do not commit
- `--allow-dirty`: Allow operation on dirty working directory
- `-m, --message`: Override commit message
- `--tag-name`: Override tag name
- `--format`: Output format: human or json
- `-c, --config`: Config file (default is .bumpx.toml)
- `-h, --help`: Help for bump command

**Examples:**
```bash
bumpx bump major    # 1.0.0 -> 2.0.0
bumpx bump minor    # 1.0.0 -> 1.1.0
bumpx bump patch    # 1.0.0 -> 1.0.1
bumpx bump pre      # 1.0.0 -> 1.0.0-alpha.1
bumpx bump build    # 1.0.0 -> 1.0.0+1
```

### bumpx init

Create a starter .bumpx.toml in the current directory.

**Usage:**
```bash
bumpx init
```

Scans common files (Cargo.toml, package.json, go.mod, etc.) and attempts to infer current version.

**Options:**
- `-y, --yes`: Non-interactive, accept defaults
- `--scheme`: Version scheme (default "semver")
- `-c, --config`: Config file (default is .bumpx.toml)
- `-h, --help`: Help for init command

### bumpx validate

Validate .bumpx.toml and exit.

**Usage:**
```bash
bumpx validate
```

**Options:**
- `--strict`: Treat missing paths as errors
- `--json`: Output in JSON format
- `-c, --config`: Config file (default is .bumpx.toml)
- `-h, --help`: Help for validate command

## Global Options

- `-c, --config`: Config file (default is .bumpx.toml)
- `--json`: Output in JSON format
- `-h, --help`: Help for bumpx
- `-v, --version`: Version for bumpx

## Exit Codes

- 0: Success
- 1: Usage or validation error
- 2: Operation failed (file I/O, parse errors, git errors)

## Examples

```bash
# Show current version
bumpx show

# Show in JSON format
bumpx show --json

# Bump minor version
bumpx bump minor

# Dry-run bump (preview changes)
bumpx bump minor --dry-run

# Bump with custom commit message
bumpx bump patch -m "Release patch version"

# Validate configuration
bumpx validate

# Validate with strict mode
bumpx validate --strict
```