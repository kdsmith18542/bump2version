# bumpx Configuration

bumpx is configured via a `.bumpx.toml` file in your project root.

## Configuration File Format

### Project Section

```toml
[project]
name = "my-project"
current_version = "1.2.3"
version_scheme = "semver"
```

| Field | Required | Default | Description |
|-------|----------|---------|-------------|
| `name` | No | - | Project name (informational) |
| `current_version` | **Yes** | - | Current version string |
| `version_scheme` | No | `semver` | Version scheme: `semver`, `calver:YYYY.MM.DD`, or `custom:<template>` |

### Files Section

Define files to update when bumping versions.

```toml
[files.cargo]
paths = ["Cargo.toml"]
type = "cargo"

[files.npm]
paths = ["package.json"]
type = "npm"

[files.custom]
paths = ["VERSION.txt"]
type = "regex"
pattern = '^version=(?P<version>.+)$'
replacement = 'version={version}'
```

#### Supported File Types

| Type | Description | Files |
|------|-------------|-------|
| `cargo` | Rust Cargo manifest | `Cargo.toml` |
| `npm` | Node.js package manifest | `package.json` |
| `gomod` | Go module file | `go.mod` (with `// version: x.y.z` comment) |
| `pyproject` | Python project | `pyproject.toml` |
| `dotnet` | .NET project | `.csproj`, `AssemblyInfo.cs` |
| `regex` | Generic pattern matching | Any text file |

#### Regex Type Options

For `type = "regex"`:

| Field | Required | Description |
|-------|----------|-------------|
| `pattern` | **Yes** | Regular expression with named group `(?P<version>...)` |
| `replacement` | No | Replacement string with `{version}` placeholder |

### Git Section

```toml
[git]
enable = true
commit = true
commit_message = "chore(release): {version}"
tag = true
tag_name = "v{version}"
tag_message = "Release {version}"
allow_dirty = false
```

| Field | Default | Description |
|-------|---------|-------------|
| `enable` | `false` | Enable git integration |
| `commit` | `true` | Create a commit after bumping |
| `commit_message` | `chore(release): {version}` | Commit message template |
| `tag` | `true` | Create a tag after bumping |
| `tag_name` | `v{version}` | Tag name template |
| `tag_message` | `Release {version}` | Tag message template |
| `allow_dirty` | `false` | Allow bumping with uncommitted changes |

#### Template Variables

- `{version}` or `{new_version}`: The new version
- `{old_version}` or `{current_version}`: The previous version

### Hooks Section

```toml
[hooks]
pre_bump = "scripts/pre_bump.sh"
post_bump = "scripts/post_bump.sh"
```

| Field | Default | Description |
|-------|---------|-------------|
| `pre_bump` | - | Command to run before bumping |
| `post_bump` | - | Command to run after bumping |

#### Hook Environment Variables

The following environment variables are set when hooks run:

- `BUMPX_OLD_VERSION`: Previous version
- `BUMPX_NEW_VERSION`: New version
- `BUMPX_PART`: Part being bumped (e.g., `minor`)
- `BUMPX_CONFIG_PATH`: Path to config file

## Example Configuration

```toml
[project]
name = "my-monorepo"
current_version = "1.2.3"
version_scheme = "semver"

[files.cargo]
paths = ["backend/Cargo.toml"]
type = "cargo"

[files.npm]
paths = ["frontend/package.json", "cli/package.json"]
type = "npm"

[files.version_file]
paths = ["VERSION"]
type = "regex"
pattern = '(?P<version>\d+\.\d+\.\d+)'
replacement = '{version}'

[git]
enable = true
commit = true
commit_message = "release: v{version}"
tag = true
tag_name = "v{version}"

[hooks]
post_bump = "./scripts/generate-changelog.sh"
```
