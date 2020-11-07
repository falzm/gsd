package gsd

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func testStepFunc(state *State, str string) error {
	if s, ok := state.Load("test"); !ok {
		state.Store("test", str)
	} else {
		state.Store("test", s.(string)+str)
	}

	return nil
}

func TestPlanOptContinueOnError(t *testing.T) {
	plan := &Plan{}

	require.NoError(t, PlanOptContinueOnError()(plan))
	require.True(t, plan.continueOnError)
}

func TestPlanOptLimitDuration(t *testing.T) {
	plan := &Plan{}

	require.NoError(t, PlanOptLimitDuration(time.Second)(plan))
	require.Equal(t, time.Second, plan.maxDuration)
}

func TestNewPlan(t *testing.T) {
	bogusOpt := func() PlanOpt { return func(p *Plan) error { return errors.New("blah") } }
	_, err := NewPlan(bogusOpt())
	require.Error(t, err)

	actual, err := NewPlan(PlanOptContinueOnError())
	expected := &Plan{
		state:           &State{sync.Map{}},
		steps:           list.New(),
		continueOnError: true,
	}

	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestPlan_AddStep(t *testing.T) {
	plan := &Plan{
		state: &State{sync.Map{}},
		steps: list.New(),
	}

	testStep := &pauseStep{}

	plan.AddStep(testStep)

	require.Equal(t, 1, plan.steps.Len())
	require.Equal(t, plan.steps.Front().Value, testStep)
}

func TestPlan_Execute_NoError(t *testing.T) {
	plan, err := NewPlan()
	require.NoError(t, err)

	err = plan.
		AddStep(&GenericStep{
			PreExecFunc:  func(ctx context.Context, state *State) error { return testStepFunc(state, "b") },
			ExecFunc:     func(ctx context.Context, state *State) error { return testStepFunc(state, "o") },
			PostExecFunc: func(ctx context.Context, state *State) error { return testStepFunc(state, "b") },
			CleanupFunc:  func(ctx context.Context, state *State) { _ = testStepFunc(state, "o") },
		}).
		AddStep(&GenericStep{
			PreExecFunc:  func(ctx context.Context, state *State) error { return testStepFunc(state, "k") },
			ExecFunc:     func(ctx context.Context, state *State) error { return testStepFunc(state, "e") },
			PostExecFunc: func(ctx context.Context, state *State) error { return testStepFunc(state, "l") },
			CleanupFunc:  func(ctx context.Context, state *State) { _ = testStepFunc(state, "s") },
		}).
		Execute(context.Background())

	expected := "bobkelso"
	actual, _ := plan.State().Load("test")

	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestPlan_Execute_IntermediateFail(t *testing.T) {
	plan, err := NewPlan()
	require.NoError(t, err)

	err = plan.
		AddStep(&GenericStep{
			PreExecFunc:  func(ctx context.Context, state *State) error { return testStepFunc(state, "b") },
			ExecFunc:     func(ctx context.Context, state *State) error { return testStepFunc(state, "o") },
			PostExecFunc: func(ctx context.Context, state *State) error { return testStepFunc(state, "b") },
			CleanupFunc:  func(ctx context.Context, state *State) { _ = testStepFunc(state, "o") },
		}).
		AddStep(&GenericStep{
			PreExecFunc: func(ctx context.Context, state *State) error { return errors.New("blah") },
		}).
		Execute(context.Background())

	expected := "bobo"
	actual, _ := plan.State().Load("test")

	require.Error(t, err)
	require.Equal(t, expected, actual)
}

func TestPlan_Execute_FailFirst(t *testing.T) {
	plan, err := NewPlan()
	require.NoError(t, err)

	require.Error(t, plan.
		AddStep(&GenericStep{
			PreExecFunc:  func(ctx context.Context, state *State) error { return errors.New("blah") },
			ExecFunc:     func(ctx context.Context, state *State) error { return testStepFunc(state, "b") },
			PostExecFunc: func(ctx context.Context, state *State) error { return testStepFunc(state, "o") },
			CleanupFunc:  func(ctx context.Context, state *State) { _ = testStepFunc(state, "b") },
		}).
		Execute(context.Background()))

	_, ok := plan.State().Load("test")
	require.False(t, ok)
}

func TestPlan_Execute_WithContinueOnError(t *testing.T) {
	plan, err := NewPlan(PlanOptContinueOnError())
	require.NoError(t, err)

	require.NoError(t, plan.
		AddStep(&GenericStep{
			PreExecFunc:  func(ctx context.Context, state *State) error { return errors.New("blah") },
			ExecFunc:     func(ctx context.Context, state *State) error { return testStepFunc(state, "b") },
			PostExecFunc: func(ctx context.Context, state *State) error { return testStepFunc(state, "o") },
			CleanupFunc:  func(ctx context.Context, state *State) { _ = testStepFunc(state, "b") },
		}).
		Execute(context.Background()))

	result, _ := plan.State().Load("test")
	require.Equal(t, "bob", result)
}

func TestPlan_Execute_WithTimeout(t *testing.T) {
	maxDuration := 3 * time.Second

	plan, err := NewPlan(PlanOptLimitDuration(maxDuration))
	require.NoError(t, err)

	err = plan.
		AddStep(&GenericStep{
			ExecFunc: func(ctx context.Context, state *State) error {
				select {
				case <-time.After(maxDuration * 2):
					_ = testStepFunc(state, "done")

				case <-ctx.Done():
					_ = testStepFunc(state, "time-out")
				}

				return nil
			},
		}).
		Execute(context.Background())
	require.EqualError(t, err, ErrTimeout.Error())
}

func TestPlan_Execute_WithRetries(t *testing.T) {
	plan, err := NewPlan()
	require.NoError(t, err)

	testStep := &GenericStep{
		PreExecFunc:  func(ctx context.Context, state *State) error { return testStepFunc(state, "*") },
		ExecFunc:     func(ctx context.Context, state *State) error { return testStepFunc(state, "*") },
		PostExecFunc: func(ctx context.Context, state *State) error { return errors.New("blah") },
	}

	err = plan.
		AddStep(testStep.WithRetries(3)).
		Execute(context.Background())

	actual, _ := plan.State().Load("test")
	require.Error(t, err)
	require.Equal(t, "**", actual)
}

func TestPlan_AddPause(t *testing.T) {
	plan, err := NewPlan()
	require.NoError(t, err)

	pause := time.Second * 3

	require.Neverf(t,
		func() bool {
			_ = plan.AddStep(&GenericStep{}).
				AddPause(pause).
				Execute(context.Background())
			return true
		},
		time.Second,
		100*time.Millisecond,
		"plan execution should not take less than %s", pause)
}

func TestPlan_State(t *testing.T) {
	plan, err := NewPlan()

	require.NoError(t, err)
	require.Equal(t, plan.state, plan.State())
}
