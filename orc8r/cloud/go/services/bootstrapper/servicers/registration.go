package servicers

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/bootstrapper"
	"magma/orc8r/cloud/go/services/bootstrapper/storage"
	"magma/orc8r/cloud/go/services/device"
	models2 "magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/tenants"
	"magma/orc8r/lib/go/protos"

	"github.com/go-openapi/strfmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type registrationServicer struct {
	store storage.Store
}

func NewRegistrationServer(store storage.Store) (protos.RegistrationServer, error) {
	if store == nil {
		return nil, fmt.Errorf("Storage store is nil")
	}
	return &registrationServicer{store}, nil
}

func (rs *registrationServicer) Register(c context.Context, request *protos.RegisterRequest) (*protos.RegisterResponse, error) {
	networkID := request.Token
	nonce := nonceFromToken(request.Token)

	tokenInfo, err := rs.store.GetTokenInfoFromNonce(networkID, nonce)
	if err != nil {
		return formatRegisterResponseError(
			fmt.Sprintf("Could not get token info from nonce for networkID %v and nonce %v: %v", networkID, nonce, err),
		), nil
	}

	if tokenInfo == nil {
		return formatRegisterResponseError(fmt.Sprintf("Could not find token info from nonce %v", nonce)), nil
	}

	if tokenTimedOut(tokenInfo) {
		return formatRegisterResponseError("Token has timed out. Please get another one."), nil
	}

	err = registerDevice()
	if err != nil {
		return formatRegisterResponseError(fmt.Sprintf("Error registering device: %v", err)), nil
	}

	controlProxy, err := getControlProxy()
	if err != nil {
		return formatRegisterResponseError(fmt.Sprintf("Error getting control proxy: %v", err)), nil
	}

	return &protos.RegisterResponse{
		Response: &protos.RegisterResponse_ControlProxy{
			ControlProxy: controlProxy,
		},
	}, nil
}

// TODO finish
func registerDevice(ti protos.TokenInfo, hwid string, challengeKey protos.ChallengeKey) error {
	ckKey := strfmt.Base64(challengeKey.Key)
	gatewayRecord := &models2.GatewayDevice{HardwareID: hwid,
		Key: &models2.ChallengeKey{KeyType: challengeKey.KeyType.String(),
			Key: &ckKey}}
	err := device.RegisterDevice(context.Background(), ti.Gateway.NetworkId, orc8r.AccessGatewayRecordType, ti.Gateway.LogicalId.Id, gatewayRecord, serdes.Device)
	return err
}

// TODO finish
func getControlProxy(networkID string) (string, error) {
	tenants, err := tenants.GetAllTenants(context.Background())
	if err != nil {
		return "", err
	}
	var tID int64
	tIDFound := false
	for _, t := range tenants.GetTenants() {
		for _, n := range t.Tenant.Networks {
			if n == networkID {
				tID = t.Id
				tIDFound = true
				break
			}
		}
	}

	if tIDFound == false {
		return "", status.Errorf(codes.NotFound, "TenantID for current NetworkID %d not found", networkID)
	}

	return tenants.GetControlProxy(tID)}

func nonceToToken(nonce string) string {
	return bootstrapper.TokenPrepend + nonce
}

func nonceFromToken(token string) string {
	return token[len(bootstrapper.TokenPrepend):]
}

func tokenTimedOut(info *protos.TokenInfo) bool {
	return time.Now().Before(time.Unix(0, int64(info.Timeout.Nanos)))
}

func formatRegisterResponseError(errString string) *protos.RegisterResponse {
	return &protos.RegisterResponse{
		Response: &protos.RegisterResponse_Error{
			Error: errString,
		},
	}
}
