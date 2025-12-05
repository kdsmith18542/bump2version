# Configuration Guide

## Overview

bumpx uses a `.bumpx.toml` configuration file to define project settings, versioning scheme, and file targets.

## File Format

The configuration uses TOML format with the following sections:

### [project]
Defines the project-level settings:

```toml
[project]
name = "my-project"              # Project name (optional)
current_version = "1.2.3"       # Current version
version_scheme = "semver"       # Versioning scheme: "semver" or "calver:YYYY.MM.DD"
```

### [files]
Defines the files to update during version bumps:

```toml
[files.cargo]
paths = ["Cargo.toml"]
type = "cargo"

[files.npm]
paths = ["package.json"]
type = "npm"

[files.gomod]
paths = ["go.mod"]
type = "gomod"

[files.python]
paths = ["pyproject.toml"]
type = "pyproject"

[files.generic_version_txt]
paths = ["VERSION.txt"]
type = "regex"
pattern = '^version=(?P<version>[0-9A-Za-z\.\-\+]+)$'
replacement = 'version={version}'
```

File types supported:
- `cargo`: For Cargo.toml files
- `npm`: For package.json files
- `gomod`: For go.mod files
- `pyproject`: For pyproject.toml files
- `regex`: For generic regex-based replacement

### [git]
Controls Git integration:

```toml
[git]
enable = true                    # Enable Git integration
commit = true                    # Create commits
commit_message = "chore(release): {version}"
tag = true                       # Create tags
tag_name = "v{version}"
tag_message = "Release {version}"
allow_dirty = false              # Allow operation on dirty working directory
```

### [hooks]
Defines external commands to run:

```toml
[hooks]
pre_bump = "scripts/pre_bump.sh"
post_bump = "scripts/post_bump.sh"
```

## Placeholders

The configuration supports the following placeholders:

- `{version}`: New version after bump
- `{new_version}`: New version after bump
- `{old_version}`: Current version before bump
- `{current_version}`: Current version before bump

## Example Configuration

```toml
[project]
name = "my-monorepo"
current_version = "1.2.3"
version_scheme = "semver"

[files.cargo]
paths = ["backend/Cargo.toml"]
type = "cargo"

[files.npm_app]
paths = ["frontend/package.json"]
type = "npm"

[files.go_module]
paths = ["tools/go.mod"]
type = "gomod"

[git]
enable = true
commit = true
commit_message = "chore(release): {version}"
tag = true
tag_name = "v{version}"
tag_message = "Release {version}"
allow_dirty = false

[hooks]
pre_bump = ""
post_bump = ""
```