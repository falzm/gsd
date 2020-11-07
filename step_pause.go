package gsd

import (
	"context"
	"time"
)

// pauseStep is a "virtual" Step implementation used as a convenient
// alternative to implementing a time.Sleep call in a GenericStep struct.
// This also offers the benefit of detecting when a step being executed is
// a pause step, allowing to skip it during the cleanup phase.
type pauseStep struct {
	d time.Duration
}

func (s *pauseStep) PreExec(_ context.Context, _ *State) error {
	return nil
}

func (s *pauseStep) Exec(_ context.Context, _ *State) error {
	time.Sleep(s.d)
	return nil
}

func (s *pauseStep) PostExec(_ context.Context, _ *State) error {
	return nil
}

func (s *pauseStep) Cleanup(_ context.Context, _ *State) {}

func (s *pauseStep) Retries() int {
	return 0
}
