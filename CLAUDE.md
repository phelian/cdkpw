# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

cdkpw (CDK Profile Wrapper) is a Go-based CLI tool that wraps the AWS CDK CLI to automatically inject AWS SSO profile parameters based on stack name patterns.

## Essential Commands

```bash
# Build the binary
just build

# Run tests with coverage
just test

# Install to system
just install

# Run linting (requires golangci-lint)
golangci-lint run
```

## Architecture

The codebase follows a simple wrapper pattern:

1. **Entry Point** (`cmd/cdkpw/main.go`): Initializes the wrapper and executes the command chain
2. **Argument Parser** (`cmd/cdkpw/args.go`): Parses CDK command-line arguments to extract stack names
3. **Configuration** (`cmd/cdkpw/config.go`): Loads YAML configuration and matches stack patterns to AWS profiles
4. **Execution**: Injects the matched `--profile` parameter and passes through to the actual CDK binary

Key design decisions:
- Stack name extraction handles both direct names and wildcard patterns
- Profile matching uses simple substring matching from the configuration
- All non-matching commands are passed through transparently
- Verbose logging levels (0=silent, 1=info, 2=debug) for troubleshooting

## Configuration

The tool expects a YAML configuration file at `~/.cdk/.cdkpw.yml` (or path specified in `CDKPW_CONFIG`):

```yaml
profiles:
  - match: Prod
    profile: prod_admin
  - match: Dev
    profile: dev_admin
cdkLocation: ${CDK_BIN}  # Optional, defaults to 'cdk' in PATH
verbose: 0
```

## Testing Approach

- Unit tests use the `testify` package for assertions
- Tests focus on argument parsing logic and configuration matching
- Test files are colocated with source files (*_test.go pattern)

## Development Notes

- Uses Go 1.23.8 with modules
- Linting configured via `.golangci.yml` with strict rules including varnamelen and revive
- Nix flake available for consistent development environment
- Vendored dependencies committed to repository