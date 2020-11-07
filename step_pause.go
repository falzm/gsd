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

func (s *pauseStep) PreExecFunc(_ context.Context, _ *State) error {
	return nil
}

func (s *pauseStep) ExecFunc(_ context.Context, _ *State) error {
	time.Sleep(s.d)
	return nil
}

func (s *pauseStep) PostExecFunc(_ context.Context, _ *State) error {
	return nil
}

func (s *pauseStep) CleanupFunc(_ context.Context, _ *State) {}
