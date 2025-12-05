# bumpx Specification

This document defines the technical specification for bumpx, a version management tool for polyglot repositories.

## Overview

bumpx is a single static CLI binary designed for version management and release chores across polyglot repositories. It supports multiple ecosystems (Rust, Go, Node, Python, .NET) and provides a unified interface for version bumping.

## Version Schemes

### Semantic Versioning (SemVer)

Format: `MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]`

Examples:
- `1.0.0`
- `1.0.0-alpha.1`
- `1.0.0+build.123`
- `1.0.0-alpha.1+build.123`

Supported parts:
- `major`: Increments major version, resets minor and patch to 0
- `minor`: Increments minor version, resets patch to 0
- `patch`: Increments patch version
- `pre` / `prerelease`: Increments prerelease identifier
- `build`: Increments build metadata

### Calendar Versioning (CalVer)

Format: `YYYY.MM.DD[.RELEASE]`

Examples:
- `2024.01.15`
- `2024.01.15.1`

Supported parts:
- `date`: Updates to current date, resets release counter
- `build` / `release`: Increments release number

### Custom Versioning

Format: User-defined template

Custom versioning allows for arbitrary version formats with support for bumping embedded semver-like patterns.

## Configuration

### File Format

Configuration is stored in `.bumpx.toml` using TOML format.

### Required Fields

- `project.current_version`: Current version string

### Optional Fields

- `project.name`: Project name
- `project.version_scheme`: Version scheme (`semver`, `calver:FORMAT`, `custom:TEMPLATE`)
- `files.paths`: Array of file paths to auto-detect handlers
- `handlers.<name>.paths`: Array of file paths for specific handler
- `handlers.<name>.type`: Handler type
- `handlers.<name>.pattern`: Regex pattern (for regex type)
- `handlers.<name>.replacement`: Replacement string (for regex type)
- `git.enable`: Enable git integration
- `git.commit`: Create commits on bump
- `git.commit_message`: Commit message template
- `git.tag`: Create tags on bump
- `git.tag_name`: Tag name template
- `git.tag_message`: Tag message template
- `git.allow_dirty`: Allow operation on dirty working directory
- `hooks.pre_bump`: Pre-bump hook command
- `hooks.post_bump`: Post-bump hook command

### Placeholders

Templates support the following placeholders:
- `{version}`: New version
- `{new_version}`: New version (alias)
- `{old_version}`: Previous version
- `{current_version}`: Previous version (alias)

## File Handlers

### Built-in Handlers

| Type | Files | Pattern |
|------|-------|---------|
| `cargo` | `Cargo.toml` | `version = "x.y.z"` |
| `npm` | `package.json` | JSON `version` field |
| `gomod` | `go.mod` | `// version: x.y.z` comment |
| `pyproject` | `pyproject.toml` | `version = "x.y.z"` |
| `dotnet` | `*.csproj`, `AssemblyInfo.cs` | `<Version>` or `AssemblyVersion` |
| `regex` | Any | Custom regex pattern |

### Regex Handler

The regex handler requires:
- `pattern`: Regex with capturing group for version
- `replacement`: Replacement string with `{version}` placeholder

Named group `(?P<version>...)` is preferred but first capturing group is used as fallback.

## CLI Commands

### bumpx show

Display current version.

Exit codes:
- 0: Success
- 1: Usage error
- 2: Operation error

### bumpx bump \<part\>

Bump version by specified part.

Options:
- `--dry-run`, `-n`: Preview changes without writing
- `--no-git`: Disable git operations
- `--git-commit`: Force git commit
- `--git-tag`: Force git tag
- `--no-commit`: Disable commit
- `--no-tag`: Disable tag
- `--allow-dirty`: Allow dirty working directory
- `--message`, `-m`: Custom commit message
- `--tag-name`: Custom tag name

Exit codes:
- 0: Success
- 1: Usage error (invalid part)
- 2: Operation error (file I/O, git, parse)

### bumpx init

Initialize configuration file.

Options:
- `--yes`, `-y`: Non-interactive mode
- `--scheme`: Version scheme

Exit codes:
- 0: Success
- 1: Usage error
- 2: File already exists

### bumpx validate

Validate configuration file.

Options:
- `--strict`: Treat missing files as errors

Exit codes:
- 0: Valid configuration
- 2: Invalid configuration

## Output Formats

### Human Format (default)

Human-readable text output.

### JSON Format

Structured JSON output for CI/CD integration.

Show result:
```json
{
  "version": "1.0.0"
}
```

Bump result:
```json
{
  "old_version": "1.0.0",
  "new_version": "1.0.1",
  "part": "patch",
  "files_changed": ["package.json"],
  "git_committed": "abc1234",
  "git_tagged": "v1.0.1"
}
```

Validation result:
```json
{
  "valid": true,
  "errors": [],
  "warnings": []
}
```

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | Usage or validation error |
| 2 | Operation failed (file I/O, parse errors, git errors) |

## Environment Variables

Hooks receive the following environment variables:
- `BUMPX_OLD_VERSION`: Version before bump
- `BUMPX_NEW_VERSION`: Version after bump
- `BUMPX_PART`: Bump part (major, minor, patch, etc.)
- `BUMPX_CONFIG_PATH`: Path to configuration file

## Security Considerations

- Hook commands are executed via shell and should only come from trusted config files
- File operations preserve original file permissions
- No network operations are performed (except git if configured)
