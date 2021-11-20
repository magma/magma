package registration

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	rootCAFilePath = "/var/opt/magma/certs/rootCA.pem"

	// length of timeout for the nonce
	timeoutDuration = 30 * time.Minute

	// number of characters that the nonce will have
	// 52^7 chance = 1.0280717e+12
	nonceLength = 7

	// 1.0280717e+12^5 = a lot
	maxRetriesToGenerateNonce = 5
)

type cloudRegistrationServicer struct {
	store Store
}

func NewCloudRegistrationServicer(store Store) (protos.CloudRegistrationServer, error) {
	if store == nil {
		return nil, fmt.Errorf("Storage store is nil")
	}
	return &cloudRegistrationServicer{store}, nil
}

func (crs *cloudRegistrationServicer) GetToken(c context.Context, request *protos.GetTokenRequest) (*protos.GetTokenResponse, error) {
	networkId := request.Gateway.NetworkId
	logicalId := request.Gateway.LogicalId.Id

	tokenInfo, err := crs.store.GetTokenInfoFromLogicalID(networkId, logicalId)
	if err != nil {
		glog.Infof("Could not get tokenInfo for networkID %v and logicalID %v: %v", networkId, logicalId, err)
	}

	refresh := request.Refresh || tokenInfo == nil || isTokenExpired(tokenInfo)
	if refresh {
		// TODO(reginawang3495) add a test with tokenInfo = nil to make sure it works
		tokenInfo, err = crs.generateAndSaveTokenInfo(networkId, logicalId, tokenInfo.Nonce)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Could not generate and save tokenInfo for networkID %v and logicalID %v: %v", networkId, logicalId, err)
		}
	}

	return &protos.GetTokenResponse{Timeout: tokenInfo.Timeout, Token: nonceToToken(tokenInfo.Nonce)}, nil
}

func (crs *cloudRegistrationServicer) GetInfoForGatewayRegistration(c context.Context, request *protos.GetInfoForGatewayRegistrationRequest) (*protos.GetInfoForGatewayRegistrationResponse, error) {
	rootCA, err := getRootCA()
	if err != nil {
		return nil, err
	}
	domainName := getDomainName()

	return &protos.GetInfoForGatewayRegistrationResponse{
		RootCA:               rootCA,
		DomainName:           domainName,
	}, nil
}

func (crs *cloudRegistrationServicer) GetGatewayPreregisterInfo(c context.Context, request *protos.GetGatewayPreregisterInfoRequest) (*protos.GetGatewayPreregisterInfoResponse, error) {
	nonce := nonceFromToken(request.Token)

	tokenInfo, err := crs.store.GetTokenInfoFromNonce(nonce)
	if err != nil {
		return &protos.GetGatewayPreregisterInfoResponse{
			Response: &protos.GetGatewayPreregisterInfoResponse_Error{
				Error: fmt.Sprintf("Could not get token info from nonce %v: %v", nonce, err),
			},
		}, nil
	}

	return &protos.GetGatewayPreregisterInfoResponse{
		Response:             &protos.GetGatewayPreregisterInfoResponse_GatewayPreregisterInfo{
			GatewayPreregisterInfo: &protos.GatewayPreregisterInfo{
				NetworkId:           tokenInfo.GatewayPreregisterInfo.NetworkId,
				LogicalId:            tokenInfo.GatewayPreregisterInfo.LogicalId,
			},
		},
	}, nil
}

func (crs *cloudRegistrationServicer) generateAndSaveTokenInfo(networkID string, logicalID string, oldNonce string) (*protos.TokenInfo, error) {
	nonce, err := crs.generateUniqueNonce(maxRetriesToGenerateNonce, nonceLength)
	if err != nil {
		return nil, err
	}

	t := time.Now().Add(timeoutDuration)
	timeout :=
		&timestamp.Timestamp{
			Seconds: int64(t.Second()),
			Nanos:   int32(t.Unix()),
		}

	tokenInfo := protos.TokenInfo{
		GatewayPreregisterInfo: &protos.  GatewayPreregisterInfo{
			NetworkId: networkID,
			LogicalId: logicalID,
		},
		Nonce:   nonce,
		Timeout: timeout,
	}

	err = crs.store.SetTokenInfo(oldNonce, tokenInfo)
	if err != nil {
		return nil, err
	}

	return &tokenInfo, nil
}

func (crs *cloudRegistrationServicer) generateUniqueNonce(maxRetries int, length int) (string, error) {
	for i := 0; i < maxRetries; i++ {
		nonce := generateSecureNonce(length)
		isUnique, err := crs.store.IsNonceUnique(nonce)
		if err != nil {
			return "", err
		}

		if isUnique {
			return nonce, nil
		}
	}
	return "", status.Errorf(codes.Internal, "Failed to create unique nonce %d times", maxRetries)
}

func getRootCA() (string, error) {
	body, err := ioutil.ReadFile(rootCAFilePath)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// TODO(#10437)
func getDomainName() string {
	return "localhost"
}
