# How to contribute

We'd love to accept your patches and contributions to this project. There are just a few small guidelines you need to follow.

## Guidelines

1. Write your patch
2. Add a test case to your patch
3. Make sure that `make test` runs properly
4. Send your patch as a PR

## Setup

1. Fork & clone the repo
2. Install [Go](https://go.dev/doc/install) (1.21 or later)
3. Run `make build` from the root directory
4. Run `make test` to run the tests

## Development

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Formatting

```bash
make fmt
```

### Linting

Install [golangci-lint](https://golangci-lint.run/welcome/install/) and run:

```bash
make lint
```

## Project Structure

```
├── cmd/bumpx/         # Main CLI entry point
├── internal/
│   ├── config/        # Configuration loading and validation
│   ├── core/          # Core bump logic
│   ├── files/         # File handlers (cargo, npm, go.mod, etc.)
│   ├── git/           # Git integration
│   ├── hooks/         # Pre/post bump hooks
│   ├── output/        # Output formatting (text, JSON)
│   └── version/       # Version parsing (semver, calver)
├── archive/           # Archived Python version for reference
└── docs/              # Documentation
```

## How to release bumpx

Execute the following commands:

```bash
git checkout main
git pull
make test
bumpx bump minor  # or major/patch
git push origin main --tags
```