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
- `--no-tag`: Do not create a tag
- `--no-commit`: Do not commit
- `-m, --message`: Override commit message
- `--tag-name`: Override tag name
- `--json`: Output in JSON format
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
- 1: Generic failure
- 2: Configuration error (invalid config, missing fields)
- 3: Version parsing/bumping error
- 4: File IO or pattern-matching error
- 5: Git-related error (dirty repo, commit/tag failure)

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