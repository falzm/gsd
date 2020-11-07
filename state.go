package gsd

import "sync"

// State represents a key/value map provided to the plan's steps for sharing
// state between steps. It is a wrapper around a sync.Map, so it is safe to be
// used concurrently.
type State struct {
	sync.Map
}

// Get returns the stored value corresponding to the key k, or the nil value
// if the requested key is not found in the state.
func (s *State) Get(k string) interface{} {
	if v, ok := s.Load(k); ok {
		return v
	}

	return nil
}
