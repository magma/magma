/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"context"
	"fmt"

	"magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
)

// ReAuth initiates a credit reauth on the gateway
func (srv *FegToGwRelayServer) ChargingReAuth(
	ctx context.Context,
	req *protos.ChargingReAuthRequest,
) (*protos.ChargingReAuthAnswer, error) {
	if err := validateFegContext(ctx); err != nil {
		return &protos.ChargingReAuthAnswer{Result: protos.ReAuthResult_OTHER_FAILURE}, err
	}
	hwID, err := getHwIDFromIMSI(ctx, req.Sid)
	if err != nil {
		return &protos.ChargingReAuthAnswer{Result: protos.ReAuthResult_SESSION_NOT_FOUND},
			fmt.Errorf("unable to get HwID from IMSI %v. err: %v", req.Sid, err)
	}
	conn, ctx, err := gateway_registry.GetGatewayConnection(gateway_registry.GwSessiondService, hwID)
	if err != nil {
		return &protos.ChargingReAuthAnswer{Result: protos.ReAuthResult_OTHER_FAILURE},
			fmt.Errorf("unable to get connection to the gateway ID: %s", hwID)
	}
	return protos.NewSessionProxyResponderClient(conn).ChargingReAuth(ctx, req)
}
