/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/plugin/models"
	"magma/orc8r/cloud/go/services/configurator"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EpsAuthConfig stores all the configs needed to run the service.
type EpsAuthConfig struct {
	LteAuthOp   []byte
	LteAuthAmf  []byte
	SubProfiles map[string]models.NetworkEpcConfigsSubProfilesAnon
}

// getConfig returns the EpsAuthConfig config for a given networkId.
func getConfig(networkID string) (*EpsAuthConfig, error) {
	iCellularConfigs, err := configurator.LoadNetworkConfig(networkID, lte.CellularNetworkType)
	if err != nil {
		return nil, err
	}
	if iCellularConfigs == nil {
		return nil, status.Error(codes.NotFound, "got nil when looking up config")
	}
	cellularConfig, ok := iCellularConfigs.(*models.NetworkCellularConfigs)
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "failed to convert config")
	}
	epc := cellularConfig.Epc
	result := &EpsAuthConfig{
		LteAuthOp:   epc.LteAuthOp,
		LteAuthAmf:  epc.LteAuthAmf,
		SubProfiles: epc.SubProfiles,
	}
	return result, nil
}
