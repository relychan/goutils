# CLAUDE.md

## Project Overview

`github.com/relychan/goutils` is a Go utility library providing common functions for other Go packages within the RelyChan ecosystem. Licensed under Apache 2.0.

- **Go version**: 1.26
- **Module**: `github.com/relychan/goutils`
- **Key dependencies**: `github.com/google/uuid`, `go.yaml.in/yaml/v4`

## Repository Structure

```
goutils/
├── all_or_strings.go      # AllOrStrings type for JSON/YAML unmarshaling (string or []string)
├── data.go                # Generic data/pointer helpers, IsNil
├── duration.go            # Custom Duration type with extended unit support (y/w/d/h/m/s/ms)
├── equal.go               # DeepEqual and Equaler interface
├── error.go               # RFC 9457 error types, sentinel errors, CatchWarnErrorFunc
├── file.go                # ReadJSONOrYAMLFile, ReadMultiFromJSONOrYAMLFile, FileReaderFromPath
├── mapstructure.go        # Map-to-struct decoding utilities
├── network.go             # URL parsing/validation, IP/SSRF protection, subnet utilities
├── regexp.go              # Regexp helpers
├── slice.go               # Generic slice utilities (Map, EqualSlice, ToAnySlice, etc.)
├── slug.go                # Slug validation
├── stringer.go            # ToString, ToDebugString, character helpers
├── time.go                # Time utilities
├── yaml.go                # YAML node helpers, tag constants
└── httpheader/
    └── header.go          # HTTP header name constants and content-type constants
```

## Common Commands

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific file/package
go test -run TestFunctionName ./...

# Lint (if golangci-lint is configured)
golangci-lint run
```

## Key Patterns and Conventions

### Error Handling
- Sentinel errors are defined in `error.go` with `Err` prefix (e.g., `ErrInvalidURI`, `ErrBlockedIP`)
- Use `CatchWarnErrorFunc` for deferred close calls where errors are non-critical
- HTTP errors follow RFC 9457 via `RFC9457Error` struct with constructor helpers (`NewBadRequestError`, `NewNotFoundError`, etc.)

### Generics
- Slice utilities use Go generics extensively (`Map[T, M]`, `ToAnySlice[T]`, `PtrToNumberSlice[T1, T2]`)
- File reading uses generics: `ReadJSONOrYAMLFile[T]`, `ReadMultiFromJSONOrYAMLFile[T]`

### File/URL Reading (`file.go`)
- `FileReaderFromPath` supports both local filesystem paths and HTTP/HTTPS URLs
- `ReadJSONOrYAMLFile[T]` decodes a single document; `ReadMultiFromJSONOrYAMLFile[T]` decodes multiple documents
- Supported extensions: `.json`, `.yaml`, `.yml`

### Network/URL Validation (`network.go`)
- `ValidateURL` / `ValidateURLString` support SSRF protection via `ValidateHTTPURLOptions`
- Options include: `AllowedSchemes`, `AllowedHosts`, `BlockedHosts`, `PublicIPOnly`, `AllowedIPRanges`, `BlockedIPRanges`
- `ValidateIP` checks individual IPs; allowed ranges take highest priority to bypass other rules
- CG-NAT subnet (100.64.0.0/10) is blocked when `PublicIPOnly` is set

### Duration (`duration.go`)
- Custom `Duration` type extending `time.Duration` with units: `y`, `w`, `d`, `h`, `m`, `s`, `ms`
- Units must appear in order from largest to smallest (e.g., `1d2h` is valid; `2h1d` is not)
- Implements JSON, text marshaling/unmarshaling and `pflag.Value`

### YAML (`yaml.go`)
- Tag constants prefixed with `YAML` (e.g., `YAMLStrTag`, `YAMLMapTag`)
- Helper functions: `GetStringValueFromYAMLMap`, `GetNodeValueFromYAMLMap`

### HTTP Headers (`httpheader/`)
- Constants for standard and non-standard HTTP header names
- Content-type constants (e.g., `ContentTypeJSON`, `ContentTypeNdJSON`)
- Cloudflare-specific header constants

## Testing
- Test files follow the `_test.go` convention, colocated with source files
- Test data lives in `testdata/` directory (`config.json`, `config.yaml`)
- No external test framework; standard `testing` package is used
