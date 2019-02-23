/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package models_test

import (
	"testing"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/checkind/obsidian/models"

	"github.com/stretchr/testify/assert"
)

func TestGatewayStatus_FromMconfig(t *testing.T) {
	testCases := []struct {
		In  *protos.GatewayStatus
		Out *models.GatewayStatus
	}{
		// Check all fields
		{
			In: &protos.GatewayStatus{
				Time: 42,
				Checkin: &protos.CheckinRequest{
					GatewayId:       "foo",
					MagmaPkgVersion: "bar",
					Status: &protos.ServiceStatus{
						Meta: map[string]string{"baz": "qux"},
					},
					SystemStatus: &protos.SystemStatus{
						Time:         101,
						CpuUser:      10,
						CpuSystem:    11,
						CpuIdle:      12,
						MemTotal:     13,
						MemAvailable: 14,
						MemUsed:      15,
						MemFree:      16,
						UptimeSecs:   17,
					},
					VpnIp:                   "facebook.com",
					KernelVersion:           "42",
					KernelVersionsInstalled: []string{"11"},
				},
			},
			Out: &models.GatewayStatus{
				CheckinTime:             42,
				HardwareID:              "foo",
				KernelVersion:           "42",
				KernelVersionsInstalled: []string{"11"},
				Meta:                    map[string]string{"baz": "qux"},
				SystemStatus: &models.SystemStatus{
					CPUIDLE:      12,
					CPUSystem:    11,
					CPUUser:      10,
					MemAvailable: 14,
					MemFree:      16,
					MemTotal:     13,
					MemUsed:      15,
					Time:         101,
					UptimeSecs:   17,
				},
				Version: "bar",
				VpnIP:   "facebook.com",
			},
		},

		// Nil status from checkin
		{
			In: &protos.GatewayStatus{
				Time: 42,
			},
			Out: &models.GatewayStatus{
				CheckinTime: 42,
				Meta:        map[string]string{},
			},
		},
	}

	for _, tc := range testCases {
		actual := &models.GatewayStatus{}
		err := actual.FromMconfig(tc.In)
		assert.NoError(t, err)
		assert.Equal(t, tc.Out, actual)
	}
}
