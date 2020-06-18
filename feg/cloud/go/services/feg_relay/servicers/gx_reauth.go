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

// PolicyReAuth installs and deactivates rules on the gateway
func (srv *FegToGwRelayServer) PolicyReAuth(
	ctx context.Context,
	req *protos.PolicyReAuthRequest,
) (*protos.PolicyReAuthAnswer, error) {
	if err := validateFegContext(ctx); err != nil {
		return &protos.PolicyReAuthAnswer{Result: protos.ReAuthResult_OTHER_FAILURE}, err
	}
	hwID, err := getHwIDFromIMSI(ctx, req.Imsi)
	if err != nil {
		return &protos.PolicyReAuthAnswer{Result: protos.ReAuthResult_SESSION_NOT_FOUND},
			fmt.Errorf("unable to get HwID from IMSI %v. err: %v", req.Imsi, err)
	}
	conn, ctx, err := gateway_registry.GetGatewayConnection(gateway_registry.GwSessiondService, hwID)
	if err != nil {
		return &protos.PolicyReAuthAnswer{Result: protos.ReAuthResult_OTHER_FAILURE},
			fmt.Errorf("unable to get connection to the gateway ID: %s", hwID)
	}
	return protos.NewSessionProxyResponderClient(conn).PolicyReAuth(ctx, req)
}
