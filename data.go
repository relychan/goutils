package goutils

import (
	"github.com/google/uuid"
)

// NewUUIDv7 creates a random UUID version 7. Fallback to v4 if there is error.
func NewUUIDv7() uuid.UUID {
	value, err := uuid.NewV7()
	if err == nil {
		return value
	}

	return uuid.New()
}

// ToPtr converts the a typed value to its pointer.
func ToPtr[T any](value T) *T {
	return &value
}
