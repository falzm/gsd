package gsd

import (
	"container/list"
	"context"
	"sync"
	"time"
)

// PlanOpt represents a Plan creation option.
type PlanOpt func(*Plan) error

// PlanOptContinueOnError instructs the plan to continue its execution when
// one or multiple steps execution fail, whereas by default it stops at the
// first error encountered.
func PlanOptContinueOnError() PlanOpt {
	return func(p *Plan) error {
		p.continueOnError = true
		return nil
	}
}

// PlanOptLimitDuration instructs the plan to time out if its execution
// exceeds the duration d.
func PlanOptLimitDuration(d time.Duration) PlanOpt {
	return func(p *Plan) error {
		p.maxDuration = d
		return nil
	}
}

// Plan represents a plan instance.
type Plan struct {
	state           *State
	steps           *list.List
	continueOnError bool
	maxDuration     time.Duration
}

// NewPlan returns a new plan.
func NewPlan(opts ...PlanOpt) (*Plan, error) {
	plan := Plan{
		state: &State{sync.Map{}},
		steps: list.New(),
	}

	for _, opt := range opts {
		if err := opt(&plan); err != nil {
			return nil, err
		}
	}

	return &plan, nil
}

// AddStep adds a new step to the plan.
func (p *Plan) AddStep(step Step) *Plan {
	p.steps.PushBack(step)

	return p
}

// AddPause injects a pause of duration d after the latest step added to the
// plan. Note: pauses are ignored during the cleanup phase.
func (p *Plan) AddPause(d time.Duration) *Plan {
	return p.AddStep(&pauseStep{d: d})
}

// Execute executes the plan's steps sequentially until completion, or
// stops and returns a non-nil error if a step failed (unless the
// PlanOptContinueOnError option has been specified during plan creation).
func (p *Plan) Execute(ctx context.Context) error {
	var cancel context.CancelFunc

	if p.maxDuration > 0 {
		ctx, cancel = context.WithTimeout(ctx, p.maxDuration)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}

	errCh := make(chan error)
	go func(cancelFunc context.CancelFunc) {
		var (
			lastOK *list.Element
			err    error
		)

		defer cancelFunc()

		for s := p.steps.Front(); s != nil; s = s.Next() {
			step := s.Value.(Step)

			for attempt := 0; attempt <= step.Retries(); attempt++ {
				if err = ctx.Err(); err != nil {
					errCh <- err
					return
				}

				if err = step.PreExec(ctx, p.state); err != nil {
					if !p.continueOnError {
						goto stop
					}
					if step.Retries() > attempt {
						continue
					}
				}

				if err = step.Exec(ctx, p.state); err != nil {
					if !p.continueOnError {
						goto stop
					}
					if step.Retries() > attempt {
						continue
					}
				}

				if err = step.PostExec(ctx, p.state); err != nil {
					if !p.continueOnError {
						goto stop
					}
					if step.Retries() > attempt {
						continue
					}
				}
			}

			// Save last successful step as starting point of the cleanup phase.
			lastOK = s
		}
	stop:

		for step := lastOK; step != nil; step = step.Prev() {
			if ctx.Err() != nil {
				errCh <- ctx.Err()
				return
			}

			// Skip pause steps during cleanup phase.
			if _, ok := step.Value.(*pauseStep); ok {
				continue
			}

			step.Value.(Step).Cleanup(ctx, p.state)
		}

		errCh <- err
	}(cancel)

	select {
	// Plan finished execution.
	case err := <-errCh:
		cancel()
		close(errCh)
		return err

	// Parent context aborted.
	case <-ctx.Done():
		switch ctx.Err() {
		case context.Canceled:
			return ErrCancelled

		case context.DeadlineExceeded:
			return ErrTimeout
		}
	}

	return nil
}

// State returns the plan's current state shared between steps.
func (p *Plan) State() *State {
	return p.state
}
