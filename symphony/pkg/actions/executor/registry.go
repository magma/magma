// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"fmt"

	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/pkg/errors"
)

// Registry provides a registry for the executor to register Triggers and Actions
type Registry struct {
	actions  map[core.ActionID]core.Action
	triggers map[core.TriggerID]core.Trigger
}

// NewRegistry initializes a new registry
func NewRegistry() *Registry {
	return &Registry{
		actions:  map[core.ActionID]core.Action{},
		triggers: map[core.TriggerID]core.Trigger{},
	}
}

// TriggerForID returns the Trigger for the specified ID
func (r Registry) TriggerForID(id core.TriggerID) (core.Trigger, error) {
	trigger, exists := r.triggers[id]
	if !exists {
		return nil, errors.Errorf("module does not exist: %v", id)
	}
	return trigger, nil
}

// RegisterTrigger is used by actions to register themselves to the framework
func (r Registry) RegisterTrigger(trigger core.Trigger) error {
	id := trigger.ID()
	if _, exists := r.triggers[id]; exists {
		return errors.Errorf("trigger %v already registered", id)
	}
	r.triggers[id] = trigger
	return nil
}

// MustRegisterTrigger panics if RegisterTrigger fails
func (r Registry) MustRegisterTrigger(trigger core.Trigger) {
	err := r.RegisterTrigger(trigger)
	if err != nil {
		panic(fmt.Sprintf("could not register trigger: %v", err))
	}
}

// ActionForID returns the Action for the specified ID
func (r Registry) ActionForID(id core.ActionID) (core.Action, error) {
	action, exists := r.actions[id]
	if !exists {
		return nil, errors.Errorf("module does not exist: %v", id)
	}
	return action, nil
}

// RegisterAction is used by actions to register themselves to the framework
func (r Registry) RegisterAction(action core.Action) error {
	id := action.ID()
	if _, exists := r.actions[id]; exists {
		return errors.Errorf("action %v already registered", id)
	}
	r.actions[id] = action
	return nil
}

// MustRegisterAction panics if RegisterAction fails
func (r Registry) MustRegisterAction(action core.Action) {
	err := r.RegisterAction(action)
	if err != nil {
		panic(fmt.Sprintf("could not register action: %v", err))
	}
}

// Triggers returns a slice of all registered triggers
func (r Registry) Triggers() []core.Trigger {
	v := make([]core.Trigger, 0, len(r.triggers))
	for _, value := range r.triggers {
		v = append(v, value)
	}
	return v
}
