package servicers

import (
	"golang.org/x/net/context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	cellular "magma/lte/cloud/go/services/cellular/config"
	cellular_protos "magma/lte/cloud/go/services/cellular/protos"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config"
)

// EpsAuthConfig stores all the configs needed to run the service.
type EpsAuthConfig struct {
	LteAuthOp   []byte
	LteAuthAmf  []byte
	SubProfiles map[string]*cellular_protos.NetworkEPCConfig_SubscriptionProfile
}

// getNetworkID returns the network id of the gateway sending a gRPC message
func getNetworkID(ctx context.Context) (string, error) {
	gw := protos.GetClientGateway(ctx)
	if gw == nil {
		err := status.Errorf(codes.PermissionDenied, "Missing Gateway Identity")
		return "", err
	}
	if !gw.Registered() {
		err := status.Errorf(codes.PermissionDenied, "Gateway is not registered")
		return "", err
	}
	return gw.GetNetworkId(), nil
}

// getConfig returns the EpsAuthConfig config for a given networkId.
func getConfig(networkID string) (*EpsAuthConfig, error) {
	configs, err := config.GetConfig(networkID, cellular.CellularNetworkType, networkID)
	if err != nil {
		return nil, err
	}
	if configs == nil {
		return nil, status.Error(codes.NotFound, "got nil when looking up config")
	}
	cellularConfig, ok := configs.(*cellular_protos.CellularNetworkConfig)
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "failed to convert config")
	}
	epc := cellularConfig.GetEpc()
	result := &EpsAuthConfig{
		LteAuthOp:   epc.GetLteAuthOp(),
		LteAuthAmf:  epc.GetLteAuthAmf(),
		SubProfiles: epc.GetSubProfiles(),
	}
	return result, nil
}
