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
}

func (s *GenericStep) PreExecFunc(ctx context.Context, state *State) error {
	if s.PreExec != nil {
		return s.PreExec(ctx, state)
	}

	return nil
}

func (s *GenericStep) ExecFunc(ctx context.Context, state *State) error {
	if s.Exec != nil {
		return s.Exec(ctx, state)
	}

	return nil
}

func (s *GenericStep) PostExecFunc(ctx context.Context, state *State) error {
	if s.PostExec != nil {
		return s.PostExec(ctx, state)
	}

	return nil
}

func (s *GenericStep) CleanupFunc(ctx context.Context, state *State) {
	if s.Cleanup != nil {
		s.Cleanup(ctx, state)
	}
}
