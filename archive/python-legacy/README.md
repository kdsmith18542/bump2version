# Python Legacy Version (bump2version)

This directory contains the archived Python version of bump2version, preserved for porting reference.

The Python version is no longer maintained. The project has been rewritten in Go as **bumpx**.

## Original Project

The original Python bump2version project was a fork of [bumpversion](https://github.com/peritus/bumpversion) maintained by [c4urself](https://github.com/c4urself/bump2version).

## Files

- `bumpversion/` - Python source code
- `tests/` - Python test suite
- `setup.py` - Python package setup
- `setup.cfg` - Package configuration
- `tox.ini` - Tox test configuration
- `MANIFEST.in` - Package manifest
- `.pylintrc` - Pylint configuration
- `Dockerfile` - Docker configuration for testing
- `docker-compose.yml` - Docker Compose configuration
- `.dockerignore` - Docker ignore rules

## For New Development

Please use the new Go version **bumpx** located in the root of this repository.

See the main [README](../../README.md) for documentation.
