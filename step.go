package gsd

import "context"

// Step represents the interface to implement a plan step.
// Note: for a step to be considered as successful, all *ExecFunc must
// individually be executed successfully (i.e. returned a nil error).
type Step interface {
	// PreExecFunc is a hook executed before the main step function.
	PreExecFunc(context.Context, *State) error

	// ExecFunc is the task to perform when this step of the plan is reached.
	ExecFunc(context.Context, *State) error

	// PostExecFunc is a hook executed after the main step function, unless
	// the ExecFunc function returned an error.
	PostExecFunc(context.Context, *State) error

	// CleanupFunc is a hook executed during the cleanup phase of the plan
	// execution, which consists in running the plan backwards and execute
	// each step's CleanupFunc sequentially (i.e. from the last recently
	// executed step up to the first).
	// Note: a step's cleanup hook is only executed if all of its *ExecFunc
	// have been executed successfully.
	CleanupFunc(context.Context, *State)
}