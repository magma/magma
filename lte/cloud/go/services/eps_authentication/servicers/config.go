package servicers

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	cellular "magma/lte/cloud/go/services/cellular/config"
	cellular_protos "magma/lte/cloud/go/services/cellular/protos"
	"magma/orc8r/cloud/go/services/config"
)

// EpsAuthConfig stores all the configs needed to run the service.
type EpsAuthConfig struct {
	LteAuthOp   []byte
	LteAuthAmf  []byte
	SubProfiles map[string]*cellular_protos.NetworkEPCConfig_SubscriptionProfile
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
