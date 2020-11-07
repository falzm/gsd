# GSD – Get Sh*t Done

The `gsd` package offers a simple, no-frills way of performing arbitrary
actions executed step-by-step according to a *plan*. It is inspired from
[Packer's `multistep`][packer-multistep] package.


## Usage

```go
// hello.go
package main

import (
    "context"
    "fmt"

    "github.com/falzm/gsd"
)

func main() {
    plan, err := gsd.NewPlan()
    if err != nil {
        // Handle error
    }

    err = plan.
        AddStep(&gsd.GenericStep{
            PreExecFunc: func(ctx context.Context, state *gsd.State) error {
                state.Store("who", "world")
                return nil
            },
            ExecFunc: func(ctx context.Context, state *gsd.State) error {
                who := state.Get("who").(string)
                fmt.Printf("Hello, %s!\n", who)
                return nil
            },
            PostExecFunc: func(ctx context.Context, state *gsd.State) error {
                state.Store("who", "universe")
                return nil
            },
        }).
        AddStep(&gsd.GenericStep{
            ExecFunc: func(ctx context.Context, state *gsd.State) error {
                who := state.Get("who").(string)
                fmt.Printf("Hello, %s!\n", who)
                return nil
            },
            CleanupFunc: func(ctx context.Context, state *gsd.State) {
                state.Delete("who")
                if _, ok := state.Load("who"); !ok {
                        fmt.Println("Goodbye, everybody!")
                }
            },
        }).
        Execute(context.Background())
        if err != nil {
                // Handle error
        }
}
```

```console
$ go run hello.go
Hello, world!
Hello, universe!
Goodbye, everybody!
```

In addition of the `GenericStep` structure, the `Step` interface allows you
to implement steps on your own structures:

```go
// custom.go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/falzm/gsd"
)

type hello struct {
	who string
}

func (h *hello) say() {
	fmt.Printf("Hello, %s!\n", h.who)
}

func (h *hello) PreExec(_ context.Context, _ *gsd.State) error {
    h.who = os.Getenv("USER")
	return nil
}

func (h *hello) Exec(_ context.Context, _ *gsd.State) error {
	h.say()
	return nil
}

func (h *hello) PostExec(_ context.Context, _ *gsd.State) error {
	return nil
}

func (h *hello) Cleanup(_ context.Context, _ *gsd.State) {}

func (h *hello) Retries() int {
	return 0
}

func main() {
	plan, err := gsd.NewPlan()
	if err != nil {
		// Handle error
	}

	err = plan.
		AddStep(new(hello)).
		Execute(context.Background())
	if err != nil {
		// Handle error
	}
}
```

```console
$ go run custom.go
Hello, marc!
```


[packer-multistep]: https://pkg.go.dev/github.com/hashicorp/packer/helper/multistep
