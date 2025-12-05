# bumpx CLI Reference

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--config` | `-c` | Config file path (default: `.bumpx.toml`) |
| `--json` | | Output in JSON format |
| `--help` | `-h` | Show help |
| `--version` | `-v` | Show version |

## Commands

### `bumpx init`

Initialize a new bumpx configuration file.

```bash
bumpx init [flags]
```

**Flags:**
| Flag | Short | Description |
|------|-------|-------------|
| `--yes` | `-y` | Non-interactive, accept defaults |
| `--scheme` | | Version scheme (default: `semver`) |

**Examples:**
```bash
# Create config interactively
bumpx init

# Create config with defaults
bumpx init --yes

# Create config with CalVer
bumpx init --scheme "calver:YYYY.MM.DD"
```

### `bumpx show`

Display the current project version.

```bash
bumpx show [flags]
```

**Examples:**
```bash
# Show version
bumpx show
# Output: 1.2.3

# Show version as JSON
bumpx show --json
# Output: {"version": "1.2.3"}
```

### `bumpx bump`

Bump the version and update all configured files.

```bash
bumpx bump <part> [flags]
```

**Arguments:**
| Argument | Description |
|----------|-------------|
| `part` | Version part to bump |

**Version Parts by Scheme:**

| Scheme | Parts |
|--------|-------|
| SemVer | `major`, `minor`, `patch`, `pre`, `build` |
| CalVer | `date`, `build` |

**Flags:**
| Flag | Short | Description |
|------|-------|-------------|
| `--dry-run` | `-n` | Preview changes without writing |
| `--no-git` | | Skip git operations |
| `--no-tag` | | Skip tag creation |
| `--no-commit` | | Skip commit creation |
| `--message` | `-m` | Override commit message |
| `--tag-name` | | Override tag name |

**Examples:**
```bash
# Bump minor version
bumpx bump minor

# Bump major version (dry run)
bumpx bump major --dry-run

# Bump patch with custom commit message
bumpx bump patch -m "Release v{version}"

# Bump without git operations
bumpx bump minor --no-git

# Bump with JSON output
bumpx bump patch --json
# Output: {"old_version":"1.2.3","new_version":"1.2.4","part":"patch","changed_files":["Cargo.toml"],"git_commit":"abc1234","git_tag":"v1.2.4"}
```

### `bumpx validate`

Validate the configuration file.

```bash
bumpx validate [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--strict` | Treat missing file paths as errors |

**Examples:**
```bash
# Validate config
bumpx validate

# Validate with strict mode
bumpx validate --strict

# Validate with JSON output
bumpx validate --json
```

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | Generic failure |
| 2 | Configuration error |
| 3 | Version parsing/bumping error |
| 4 | File IO or pattern-matching error |
| 5 | Git-related error |

## Environment Variables

bumpx respects the following environment variables:

| Variable | Description |
|----------|-------------|
| `BUMPX_CONFIG` | Path to config file (overrides `--config`) |

When running hooks, bumpx sets these variables:

| Variable | Description |
|----------|-------------|
| `BUMPX_OLD_VERSION` | Previous version |
| `BUMPX_NEW_VERSION` | New version |
| `BUMPX_PART` | Part being bumped |
| `BUMPX_CONFIG_PATH` | Path to config file |
