package registration

import (
	"context"
	"fmt"

	"github.com/go-openapi/strfmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/bootstrapper"
	"magma/orc8r/cloud/go/services/device"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/tenants"
	"magma/orc8r/lib/go/protos"
)

// RegistrationService is public for ease of testing and mocking out functions
type RegistrationService struct {
	GetGatewayDeviceInfo func(ctx context.Context, token string) (*protos.GatewayDeviceInfo, error)
	RegisterDevice       func(deviceInfo protos.GatewayDeviceInfo, hwid *protos.AccessGatewayID, challengeKey *protos.ChallengeKey) error
	GetControlProxy      func(networkID string) (string, error)
}

func NewRegistrationServicer() protos.RegistrationServer {
	return &RegistrationService{
		GetGatewayDeviceInfo: bootstrapper.GetGatewayDeviceInfo,
		RegisterDevice:       RegisterDevice,
		GetControlProxy:      GetControlProxy,
	}
}

func (r *RegistrationService) Register(c context.Context, request *protos.RegisterRequest) (*protos.RegisterResponse, error) {
	nonce, err := NonceFromToken(request.Token)
	if err != nil {
		return nil, err
	}

	deviceInfo, err := r.GetGatewayDeviceInfo(context.Background(), nonce)
	if err != nil {
		clientErr := makeErr(fmt.Sprintf("could not get device info from token %v: %v", request.Token, err))
		return clientErr, nil
	}

	err = r.RegisterDevice(*deviceInfo, request.Hwid, request.ChallengeKey)
	if err != nil {
		clientErr := makeErr(fmt.Sprintf("error registering device: %v", err))
		return clientErr, nil
	}

	controlProxy, err := r.GetControlProxy(deviceInfo.NetworkId)
	if err != nil {
		clientErr := makeErr(fmt.Sprintf("error getting control proxy: %v", err))
		return clientErr, nil
	}

	res := &protos.RegisterResponse{
		Response: &protos.RegisterResponse_ControlProxy{
			ControlProxy: controlProxy,
		},
	}
	return res, nil
}

func RegisterDevice(deviceInfo protos.GatewayDeviceInfo, hwid *protos.AccessGatewayID, challengeKey *protos.ChallengeKey) error {
	challengeKeyBase64 := strfmt.Base64(challengeKey.Key)
	gatewayRecord := &models.GatewayDevice{
		HardwareID: hwid.Id,
		Key:        &models.ChallengeKey{KeyType: challengeKey.KeyType.String(), Key: &challengeKeyBase64},
	}
	err := device.RegisterDevice(context.Background(), deviceInfo.NetworkId, orc8r.AccessGatewayRecordType, hwid.Id, gatewayRecord, serdes.Device)
	return err
}

func GetControlProxy(networkID string) (string, error) {
	// TODO(#10536) Move functionality to get control_proxy from networkID into tenants service
	tenantList, err := tenants.GetAllTenants(context.Background())
	if err != nil {
		return "", err
	}

	var tenantID int64
	isTenantFound := false
	for _, t := range tenantList.GetTenants() {
		for _, n := range t.Tenant.Networks {
			if n == networkID {
				tenantID = t.Id
				isTenantFound = true
				break
			}
		}
	}

	if !isTenantFound {
		return "", status.Errorf(codes.NotFound, "tenantID for current NetworkID %v not found", networkID)
	}

	cp, err := tenants.GetControlProxy(context.Background(), tenantID)
	if err != nil {
		return "", err
	}

	return cp.ControlProxy, nil
}

// makeErr makes a protos.RegisterResponse_Error for protos.RegisterResponse
func makeErr(errString string) *protos.RegisterResponse {
	errRes := &protos.RegisterResponse{
		Response: &protos.RegisterResponse_Error{
			Error: errString,
		},
	}
	return errRes
}
