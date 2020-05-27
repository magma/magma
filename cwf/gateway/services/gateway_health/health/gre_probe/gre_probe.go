/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package gre_probe

// GREProbe defines an interface to begin a probe of GRE endpoints
// and fetch that status at a later point.
type GREProbe interface {
	// Start begins the probe of the GRE endpoint(s).
	Start() error

	// Stop stops the probe of the GRE endpoint(s).
	Stop()

	// GetStatus fetches the status of the GRE probe. The GREProbeStatus
	// returned contains slices of reachable and unreachable endpoint IPs.
	GetStatus() *GREProbeStatus
}

type GREProbeStatus struct {
	Reachable   []string
	Unreachable []string
}

type GREEndpointStatus uint

const (
	EndpointReachable   GREEndpointStatus = 0
	EndpointUnreachable GREEndpointStatus = 1
)

type DummyGREProbe struct{}

func (d *DummyGREProbe) Start() error {
	return nil
}

func (d *DummyGREProbe) GetStatus() *GREProbeStatus {
	return &GREProbeStatus{
		Reachable:   []string{},
		Unreachable: []string{},
	}
}
