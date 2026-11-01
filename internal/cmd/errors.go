package cmd

import "errors"

var (
	// ErrInvalidURL error.
	ErrInvalidURL = errors.New("cmd: invalid url")

	// ErrInvalidStatusCode error.
	ErrInvalidStatusCode = errors.New("cmd: invalid status code")

	// ErrInvalidHeadlessShell error.
	ErrInvalidHeadlessShell = errors.New("cmd: invalid headless shell")
)
