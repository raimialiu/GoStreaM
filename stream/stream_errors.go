package stream

import "errors"

// Common error definitions
var (
	ErrEmptyStream      = errors.New("stream is empty")
	ErrMultipleElements = errors.New("stream contains multiple elements")
	ErrIndexOutOfRange  = errors.New("index out of range")
	ErrStreamClosed     = errors.New("stream is closed")
)
