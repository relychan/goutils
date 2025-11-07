package goutils

import (
	"context"
	"errors"
	"log/slog"
)

var (
	// ErrInvalidHTTPScheme represents an invalid http(s) scheme error.
	ErrInvalidHTTPScheme = errors.New("invalid http(s) scheme")
	// ErrInvalidSlug represents an invalid slug error.
	ErrInvalidSlug = errors.New("invalid slug")
)

// CatchWarnErrorFunc catches the closer function and prints error with the WARN level.
func CatchWarnErrorFunc(fn func() error) {
	err := fn()
	if err != nil {
		slog.Warn(err.Error())
	}
}

// CatchWarnContextErrorFunc catches the closer function with context and prints error with the WARN level.
func CatchWarnContextErrorFunc(fn func(ctx context.Context) error) {
	err := fn(context.TODO())
	if err != nil {
		slog.Warn(err.Error())
	}
}
