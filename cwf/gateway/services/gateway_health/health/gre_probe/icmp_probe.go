/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package gre_probe

import (
	"sync"
	"time"

	"magma/cwf/cloud/go/protos/mconfig"

	"github.com/sparrc/go-ping"
)

// ICMPProbe implements the GRE probe interface
// using ICMP over GRE.
type ICMPProbe struct {
	Endpoints []*mconfig.CwfGatewayHealthConfigGrePeer
	Interval  time.Duration
	PktCount  int

	// Maps endpoints IPs to their current GREEndpointStatus
	endpointStatus map[string]GREEndpointStatus
	sync.RWMutex   // R/W lock synchronizing map endpoint status access
}

const (
	defaultPingTimeout = 5 * time.Second
)

// NewICMPProbe create a new ICMPProbe with the provided endpoints and
// probe interval.
func NewICMPProbe(endpoints []*mconfig.CwfGatewayHealthConfigGrePeer, interval uint32, pktCount int) *ICMPProbe {
	return &ICMPProbe{
		Endpoints:      endpoints,
		Interval:       time.Duration(interval) * time.Second,
		PktCount:       pktCount,
		endpointStatus: map[string]GREEndpointStatus{},
	}
}

// Start begins the ICMP probes of the ICMPProbe's endpoints.
func (i *ICMPProbe) Start() error {
	var pingers []*ping.Pinger
	for _, endpoint := range i.Endpoints {
		p, err := ping.NewPinger(endpoint.Ip)
		if err != nil {
			return err
		}
		// Need privileged mode to work with docker
		p.SetPrivileged(true)
		pingers = append(pingers, p)
	}
	startProbe := func() {
		for {
			time.Sleep(i.Interval)
			i.executeProbe(pingers)
		}
	}

	go startProbe()
	return nil
}

// GetStatus returns the current GREEndpointStatus of each endpoint.
func (i *ICMPProbe) GetStatus() *GREProbeStatus {
	var reachable []string
	var unreachable []string
	i.RLock()
	defer i.RUnlock()
	for ip, status := range i.endpointStatus {
		if status == EndpointUnreachable {
			unreachable = append(unreachable, ip)
		} else if status == EndpointReachable {
			reachable = append(reachable, ip)
		}
	}
	return &GREProbeStatus{
		Reachable:   reachable,
		Unreachable: unreachable,
	}
}

func (i *ICMPProbe) executeProbe(pingers []*ping.Pinger) {
	for _, pinger := range pingers {
		pinger.Count = i.PktCount
		pinger.Timeout = defaultPingTimeout
		// reduce frequency of updates by only updating on finish
		pinger.OnFinish = func(stats *ping.Statistics) {
			i.Lock()
			defer i.Unlock()
			if stats.PacketsRecv == 0 {
				i.endpointStatus[stats.Addr] = EndpointUnreachable
			} else {
				i.endpointStatus[stats.Addr] = EndpointReachable
			}
		}
		pinger.Run()
	}
}
