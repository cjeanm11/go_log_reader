package errors

import "errors"

var (
	ErrInsufficientData   error = errors.New("insufficient data")
	ErrOffsetOutOfRange   error = errors.New("offset out of range")
	ErrInvalidSegmentSize error = errors.New("invalid segment size")
	ErrInvalidData        error = errors.New("invalid data")
)
