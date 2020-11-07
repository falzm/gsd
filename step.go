package gsd

import "context"

// Step represents the interface to implement a plan step.
// Note: for a step to be considered as successful, all *Exec must
// individually be executed successfully (i.e. returned a nil error).
type Step interface {
	// PreExec is a hook executed before the main step function.
	PreExec(context.Context, *State) error

	// Exec is the task to perform when this step of the plan is reached,
	// unless the PreExec hook returned an error.
	Exec(context.Context, *State) error

	// PostExec is a hook executed after the main step function, unless
	// the Exec() function returned an error.
	PostExec(context.Context, *State) error

	// Cleanup is a hook executed during the cleanup phase of the plan
	// execution, which consists in running the plan backwards and execute
	// each step's Cleanup sequentially (i.e. from the last recently
	// executed step up to the first).
	// Note: a step's cleanup hook is only executed if all of its *ExecFunc
	// have been executed successfully.
	Cleanup(context.Context, *State)

	// Retries return the number of times a step execution should be retried
	// upon error. Note: all *Exec functions are retried at each
	// subsequent attempt, the implementor is responsible to track the state
	// of previous attempts internally if they don't want certain functions to
	// be retried (e.g. if the Exec function has executed successfully but the
	// PostExec hook failed, the Exec function should not be re-executed). The
	// Cleanup hook is not retried.
	Retries() int
}
