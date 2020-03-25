/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 */
// package service_manager defines and implements API for service management
package service_manager

type ServiceState int

const (
	_ ServiceState = iota
	Active
	Activating
	Deactivating
	Inactive
	Failed
	Unknown
	Error
)

// ServiceController defines service controller API for service manager providers
type ServiceController interface {
	// Name returns the type of the init system used by the GW, it should match magmad.yml "init_system" value
	Name() string
	// Start starts service and returns error if unsuccessful
	Start(service string) error
	// Stop stops service and returns error if unsuccessful
	Stop(service string) error
	// Restart restarts service and returns error if unsuccessful
	Restart(service string) error
	// GetState returns the given service state or error if unsuccessful
	GetState(service string) (ServiceState, error)
}
