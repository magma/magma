/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package system_health

import (
	"fmt"

	"github.com/coreos/go-iptables/iptables"
	"github.com/golang/glog"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

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

const (
	filterTable = "filter"
	inputChain  = "INPUT"
)

// SystemsStats define the metrics this provider will collect.
type SystemStats struct {
	CpuUtilPct float32
	MemUtilPct float32
}

// CWAGSystemHealthProvider defines a system health provider.
type CWAGSystemHealthProvider struct {
	iptables      *iptables.IPTables
	icmpInterface string
}

// NewCWAGSystemHealthProvider creates a new CWAGSystemHealthProvider with
// initialized iptables.
func NewCWAGSystemHealthProvider(eth string) (*CWAGSystemHealthProvider, error) {
	ipt, err := iptables.New()
	if err != nil {
		return nil, err
	}
	return &CWAGSystemHealthProvider{
		iptables:      ipt,
		icmpInterface: eth,
	}, nil
}

// GetSystemStats collects and return the stats defined in SystemStats.
func (c *CWAGSystemHealthProvider) GetSystemStats() (*SystemStats, error) {
	stats := &SystemStats{}
	cpuUtilPctArray, cpuErr := cpu.Percent(0, false)
	if cpuErr == nil && len(cpuUtilPctArray) == 1 {
		stats.CpuUtilPct = float32(cpuUtilPctArray[0]) / 100
	}
	virtualMem, vmErr := mem.VirtualMemory()
	if vmErr == nil {
		stats.MemUtilPct = float32(virtualMem.UsedPercent) / 100
	}
	if cpuErr != nil || vmErr != nil {
		return stats, fmt.Errorf("Error collecting system stats; CPU Result: %v, MEM Result: %v,", cpuErr, vmErr)
	}
	return stats, nil
}

// Enable removes the ICMP DROP rule from iptables for the configured interface.
// If the iptables rule doesn't exist, Enable has no effect.
func (c *CWAGSystemHealthProvider) Enable() error {
	exists, err := c.doesICMPDropExist()
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	icmpDropCmd := c.getICMPDropCmd()
	err = c.iptables.Delete(filterTable, inputChain, icmpDropCmd...)
	if err != nil {
		glog.Errorf("Unable to remove ICMP DROP rule from iptables for %s", c.icmpInterface)
		return err
	}
	glog.Infof("Successfully removed iptables ICMP DROP for %s", c.icmpInterface)
	return nil
}

// Disable adds an ICMP DROP rule from iptables for the configured interface.
// If the iptables rule already exists, Disable has no effect.
func (c *CWAGSystemHealthProvider) Disable() error {
	exists, err := c.doesICMPDropExist()
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	icmpDropCmd := c.getICMPDropCmd()
	err = c.iptables.Append(filterTable, inputChain, icmpDropCmd...)
	if err != nil {
		glog.Errorf("Error adding iptables ICMP DROP rule for %s", c.icmpInterface)
		return err
	}
	glog.Infof("Successfully added iptables ICMP DROP for %s", c.icmpInterface)
	return nil
}

func (c *CWAGSystemHealthProvider) doesICMPDropExist() (bool, error) {
	icmpDropCmd := c.getICMPDropCmd()
	return c.iptables.Exists(filterTable, inputChain, icmpDropCmd...)
}

func (c *CWAGSystemHealthProvider) getICMPDropCmd() []string {
	return []string{"-i", c.icmpInterface, "--proto", "icmp", "-j", "DROP"}
}
