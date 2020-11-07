package gsd

import "context"

// GenericStep is a generic Step implementation allowing users to provide
// arbitrary pre-exec/exec/post-exec/cleanup functions to be executed during
// the step's evaluation.
type GenericStep struct {
	PreExec  func(context.Context, *State) error
	Exec     func(context.Context, *State) error
	PostExec func(context.Context, *State) error
	Cleanup  func(context.Context, *State)

	retries        int
	preExecFuncOK  bool
	execFuncOK     bool
	postExecFuncOK bool
}

func (s *GenericStep) PreExecFunc(ctx context.Context, state *State) error {
	if s.PreExec != nil && !s.preExecFuncOK {
		return s.PreExec(ctx, state)
	}

	s.preExecFuncOK = true

	return nil
}

func (s *GenericStep) ExecFunc(ctx context.Context, state *State) error {
	if s.Exec != nil && !s.execFuncOK {
		return s.Exec(ctx, state)
	}

	s.execFuncOK = true

	return nil
}

func (s *GenericStep) PostExecFunc(ctx context.Context, state *State) error {
	if s.PostExec != nil && !s.postExecFuncOK {
		return s.PostExec(ctx, state)
	}

	s.postExecFuncOK = true

	return nil
}

func (s *GenericStep) CleanupFunc(ctx context.Context, state *State) {
	if s.Cleanup != nil {
		s.Cleanup(ctx, state)
	}
}

func (s *GenericStep) Retries() int {
	return s.retries
}

func (s *GenericStep) WithRetries(n int) Step {
	s.retries = n
	return s
}
