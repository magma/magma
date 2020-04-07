/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package system_health

// SystemHealth defines an interface to fetch system health and enable/disable
// functionality necessary for promotion/demotion from failovers
type SystemHealth interface {
	// GetSystemStats provides system level health stats
	GetSystemStats() (*SystemStats, error)

	// Disable allows the disabling of system level functionality. It is
	// up to implementors to determine specific functionality.
	Disable() error

	// Enable allows the enabling of system level functionality. It is
	// up to implementors to determine specific functionality.
	Enable() error
}

type SystemStats struct {
	CpuUtilPct float64
	MemUtilPct float64
}

type DummySystemStatsProvider struct{}

func (d *DummySystemStatsProvider) GetSystemStats() (*SystemStats, error) {
	return &SystemStats{MemUtilPct: 0.1, CpuUtilPct: 0.1}, nil
}

func (d *DummySystemStatsProvider) Enable() error {
	return nil
}

func (d *DummySystemStatsProvider) Disable() error {
	return nil
}
