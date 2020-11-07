package gsd

import "sync"

// State represents a key/value map provided to the plan's steps for sharing state between steps.
// It is a wrapper around a sync.Map, so it is safe to be used concurrently.
type State struct {
	sync.Map
}
