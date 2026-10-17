# AGENTS.md

Guidance for AI coding agents working in this repository.

## Project Overview

`github.com/relychan/goutils` is a Go utility library providing common functions for other Go packages within the RelyChan ecosystem. Licensed under Apache 2.0.

- **Go version**: 1.27 (see `go.mod`)
- **Module**: `github.com/relychan/goutils`
- **External dependencies**: only `github.com/pb33f/go-yaml`. UUID support uses the standard library `uuid` package — do not add `github.com/google/uuid`.

## Repository Structure

```
goutils/
├── all_or_strings.go      # AllOrListString, Wildcard, AllOrListWildcardString (JSON/YAML: "*" or []string)
├── data.go                # IsZero, IsZeroPtr, IsZeroer, IsNil, UnwrapPointerFromReflectValue
├── doc.go                 # Package doc
├── duration.go            # Duration type with extended units (y/w/d/h/m/s/ms), ParseDuration
├── equal.go               # DeepEqual, Equaler, Equal*Ptr / EqualMap* / EqualComparable* helpers
├── error.go               # Sentinel errors, CatchWarnErrorFunc, HTTPErrorWithExtensions
├── file.go                # ReadJSONOrYAMLFile, ReadMultiFromJSONOrYAMLFile, FileReaderFromPath, DownloadFile* options
├── format.go              # ToString, FormatScalar, FormatScalarReflection, NullStr
├── http.go                # Doer interface, ExtractHeaders, CloseResponse (with body drain)
├── json.go                # LoadMultiJSONDocumentStream
├── map.go                 # GetSortedKeys, ToAnyMap
├── mapstructure.go        # Decode*/As* helpers converting any/reflect values to strings, numbers, booleans, slices
├── network.go             # ParseSubnet, CG-NAT subnet
├── regexp.go              # RegexpMatcher (text/JSON marshalable regexp)
├── slice.go               # Generic slice utilities (Map, EqualSlice, ToAnySlice, ToNumberSlice, etc.)
├── slug.go                # Slug type and validation
├── stringer.go            # Character/string helpers (IsDigit, ParseIntInRange, HasStringPrefixFold, ...)
├── time.go                # Time type wrapping time.Time
├── uri.go                 # URL/URI parsing, ValidateURLWithOptions, AppendURL, ValidateURLPath
├── validation.go          # Format validators (UUID, email, hostname, date/time, JSON pointer, IP), ValidateIP, ErrCode* constants
├── yaml.go                # YAML tag constants, YAML map helpers, LoadMultiYAMLDocumentStream
├── docs/error.schema.json # JSON schema of the HTTP error format
├── httpheader/
│   └── header.go          # HTTP header name constants and content-type constants
├── httperror/
│   └── error.go           # RFC 9457 HTTPError, ValidationError, constructor helpers, MaxPostCloseReadBytes
└── testdata/              # Fixtures for file reading tests
```

## Common Commands

CI (`.github/workflows/test.yaml`) runs these; run them before finishing a change:

```bash
gofmt -d -s .        # must produce no output
go vet ./...
golangci-lint run    # config in .golangci.yml (v2, nearly all linters enabled)
go test ./...

# Run a single test
go test -run TestFunctionName ./...
```

## Key Patterns and Conventions

### General
- Every `.go` file starts with the Apache 2.0 license header (`Copyright 2026 RelyChan Pte. Ltd`); copy it from an existing file.
- All exported identifiers have doc comments.
- The linter set is strict. Fix lint issues rather than suppressing them; where a `//nolint:<linter>` is unavoidable (e.g. `cyclop`, `funlen` on large switch-based functions), name the specific linter, as existing code does.

### Error Handling
- Sentinel errors use the `Err` prefix (e.g. `ErrBlockedIP`, `ErrMalformedJSON`). General ones live in `error.go`; domain-specific ones sit next to their code (e.g. duration errors in `duration.go`).
- Validators return `*httperror.ValidationError` with a `Code` from the `ErrCode*` constants in `validation.go` (e.g. `ErrCodeInvalidURL`), plus `Detail` and, where useful, `Hint`.
- Use `CatchWarnErrorFunc` / `CatchWarnContextErrorFunc` for deferred close calls where errors are non-critical.
- RFC 9457 types live in `httperror/`: `HTTPError`, `ValidationError` (fields `Detail`, `Hint`, `Pointer`, `Parameter`, `Header`, `Code`), and constructors such as `NewBadRequestError`, `NewNotFoundError`.
- `HTTPErrorWithExtensions` (root package) embeds `httperror.HTTPError` and adds arbitrary extension fields; it has custom JSON marshal/unmarshal, including `message` → `detail` aliasing.
- `CloseResponse` drains the response body (up to `httperror.MaxPostCloseReadBytes`, 256 KiB) before closing so the connection can be reused. Bodies with a known length above that limit are closed without draining.

### Generics
- Prefer generics over `any` for new helpers (`Map[T, M]`, `ToAnySlice[T]`, `PtrToNumberSlice[T1, T2]`, `ReadJSONOrYAMLFile[T]`).
- Many helpers come in pairs: a value version and a `*Reflection` version taking `reflect.Value` (e.g. `DecodeString` / `DecodeStringReflection`). Keep both in sync when changing behavior.

### Equality and zero checks
- Types can implement `Equaler[T]` (`Equal(T) bool`) and `IsZeroer` (`IsZero() bool`); `DeepEqual`, `EqualSlicePtr`, `IsZero` etc. use them.
- Equality helpers take an `omitZero` flag that treats zero values as equal to absent ones.

### File/URL Reading (`file.go`)
- `FileReaderFromPath` supports local filesystem paths and HTTP/HTTPS URLs; download behavior is configured with `DownloadFile*` options (HTTP client, allowed/blocked hosts, include/exclude paths).
- `ReadJSONOrYAMLFile[T]` decodes a single document; `ReadMultiFromJSONOrYAMLFile[T]` decodes multiple documents.
- Supported extensions: `.json`, `.yaml`, `.yml`.

### URL Parsing and SSRF Protection (`uri.go`, `validation.go`, `network.go`)
- Parsing helpers: `ParseURL`, `ParseHTTPURL`, `ParseAbsoluteURL`, `ParseAbsoluteHTTPURL`, `ParseAbsoluteURI`, `ParseRelativeURI`, `ParseFilePathOrURL`, `ParseFilePathOrHTTPURL`, `SplitPathQueryFragment`.
- `ValidateURLWithOptions` / `ParseAndValidateURLWithOptions` take `ValidateHTTPURLOptions`: `AllowedSchemes`, `AllowedHosts`, `BlockedHosts`, `PublicIPOnly`, `AllowedIPRanges`, `BlockedIPRanges`, and an optional `LookupIP` override (use it in tests to avoid real DNS).
- `ValidateIPOrDomain` resolves domains before checking; `ValidateIP` checks a single IP with these rules, in order:
  1. An IP in `AllowedIPRanges` is always allowed.
  2. With `PublicIPOnly`, private, non-global-unicast, link-local multicast and CG-NAT (100.64.0.0/10) IPs are blocked.
  3. An IP in `BlockedIPRanges` is blocked.
  4. If `AllowedIPRanges` is set, every other IP is blocked.
- `ValidateURLPath` rejects absolute paths; `AppendURL` appends a path (with optional `?query` and `#fragment`) to a `*url.URL` in place.

### Duration (`duration.go`)
- `Duration` extends `time.Duration` with units `y`, `w`, `d`, `h`, `m`, `s`, `ms` (a year is 365 days).
- Units must appear from largest to smallest (`1d2h` is valid; `2h1d` is not).
- Implements JSON and text marshaling, and `pflag.Value` (`Set`, `Type`).

### YAML (`yaml.go`)
- Tag constants use the `YAML` prefix (e.g. `YAMLStrTag`, `YAMLMapTag`).
- Helpers: `GetStringValueFromYAMLMap`, `GetNodeValueFromYAMLMap`, `LoadMultiYAMLDocumentStream`.

### HTTP Headers (`httpheader/`)
- Constants for standard and non-standard HTTP header names, including Cloudflare-specific headers.
- Content-type constants (e.g. `ContentTypeJSON`, `ContentTypeNdJSON`).

## Testing
- Tests are colocated as `_test.go` files and use only the standard `testing` package — no assertion libraries.
- Fixtures live in `testdata/` (`config.json`, `config.yaml`, `config.txt`, `multi_config.json`, `multi_config.yaml`).
- Add tests for new exported functions and for each bug fix, covering edge cases (nil, empty, malformed input).
