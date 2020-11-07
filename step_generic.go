package gsd

import "context"

// GenericStep is a generic Step implementation allowing users to provide
// arbitrary pre-exec/exec/post-exec/cleanup functions to be executed during
// the step's evaluation.
type GenericStep struct {
	PreExecFunc  func(context.Context, *State) error
	ExecFunc     func(context.Context, *State) error
	PostExecFunc func(context.Context, *State) error
	CleanupFunc  func(context.Context, *State)

	retries    int
	preExecOK  bool
	execOK     bool
	postExecOK bool
}

func (s *GenericStep) PreExec(ctx context.Context, state *State) error {
	if s.PreExecFunc != nil && !s.preExecOK {
		return s.PreExecFunc(ctx, state)
	}

	s.preExecOK = true

	return nil
}

func (s *GenericStep) Exec(ctx context.Context, state *State) error {
	if s.ExecFunc != nil && !s.execOK {
		return s.ExecFunc(ctx, state)
	}

	s.execOK = true

	return nil
}

func (s *GenericStep) PostExec(ctx context.Context, state *State) error {
	if s.PostExecFunc != nil && !s.postExecOK {
		return s.PostExecFunc(ctx, state)
	}

	s.postExecOK = true

	return nil
}

func (s *GenericStep) Cleanup(ctx context.Context, state *State) {
	if s.CleanupFunc != nil {
		s.CleanupFunc(ctx, state)
	}
}

func (s *GenericStep) Retries() int {
	return s.retries
}

func (s *GenericStep) WithRetries(n int) Step {
	s.retries = n
	return s
}
