package gsd

import "errors"

// ErrCancelled represents an error reported if a plan is cancelled.
var ErrCancelled = errors.New("plan execution cancelled")

// ErrTimeout represents an error reported if a plan takes too long to execute.
var ErrTimeout = errors.New("plan execution duration limit exceeded")
