/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 */
// package service_manager defines and implements API for service management
package service_manager

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// DockerController - docker based controller implementation
type DockerController struct{}

// Name returns Docker controller type name
func (DockerController) Name() string {
	return "docker"
}

// Start starts service and returns error if unsuccessful
func (dc DockerController) Start(service string) error {
	return exec.Command(dc.Name(), "start", service).Run()
}

// Stop stops service and returns error if unsuccessful
func (dc DockerController) Stop(service string) error {
	return exec.Command(dc.Name(), "stop", service).Run()
}

// Restart restarts service and returns error if unsuccessful
func (dc DockerController) Restart(service string) error {
	return exec.Command(dc.Name(), "restart", service).Run()
}

// GetState returns the given service state or error if unsuccessful
func (dc DockerController) GetState(service string) (ServiceState, error) {
	out, err := exec.Command(dc.Name(), "inspect", service).Output()
	if err != nil {
		return Error, err
	}
	state, err := parseDockerInspectResult(out)
	if err != nil {
		err = fmt.Errorf("%v for service '%s', raw output: %s", err, service, string(out))
	}
	return state, err
}

func parseDockerInspectResult(out []byte) (ServiceState, error) {
	inspectRes := &[]struct {
		State struct {
			Status string
		}
	}{}
	err := json.Unmarshal(out, inspectRes)
	if err != nil {
		return Error, err
	}
	if len(*inspectRes) == 0 {
		return Error, fmt.Errorf("Empty returned status")
	}
	res := Unknown
	if returnedState, ok := dockerStates[(*inspectRes)[0].State.Status]; ok {
		res = returnedState
	}
	return res, nil
}

var dockerStates = map[string]ServiceState{
	"created":    Inactive,
	"restarting": Activating,
	"running":    Active,
	"paused":     Inactive,
	"removing":   Deactivating,
	"exited":     Inactive,
	"dead":       Failed,
}
