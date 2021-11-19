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
	nonceLength = 12
	timeoutDuration = 30 * time.Minute
	maxRetriesToGenerateNonce = 10
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
	n := request.Gateway.NetworkId
	l := request.Gateway.LogicalId.Id

	tokenInfo, err := crs.store.GetTokenInfoFromLogicalID(n, l)
	if err != nil {
		glog.Infof("Could not get tokenInfo for networkID %v and logicalID %v: %v", n, l, err)
	}

	refresh := request.Refresh || tokenInfo == nil || tokenTimedOut(tokenInfo)
	if refresh {
		// TODO(reginawang3495) add a test with tokenInfo = nil to make sure it works
		tokenInfo, err = crs.generateAndSaveTokenInfo(n, l, tokenInfo.Nonce)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Could not generate and save tokenInfo for networkID %v and logicalID %v: %v", n, l, err)
		}
	}

	return &protos.GetTokenResponse{Timeout: tokenInfo.Timeout, Token: nonceToToken(tokenInfo.Nonce)}, nil
}

func (crs *cloudRegistrationServicer) GetGatewayRegistrationInfo(c context.Context, request *protos.GetGatewayRegistrationInfoRequest) (*protos.GetGatewayRegistrationInfoResponse, error) {
	nonce := nonceFromToken(request.Token)

	tokenInfo, err := crs.store.GetTokenInfoFromNonce(nonce)
	if err != nil {
		return &protos.GetGatewayRegistrationInfoResponse{
			Response: &protos.GetGatewayRegistrationInfoResponse_Error{
				Error: fmt.Sprintf("Could not get token info from nonce %v: %v", nonce, err),
			},
		}, nil
	}

	rootCA, err := getRootCA()
	if err != nil {
		return &protos.GetGatewayRegistrationInfoResponse{
			Response: &protos.GetGatewayRegistrationInfoResponse_Error{
				Error: fmt.Sprintf("%v", err),
			},
		}, nil
	}
	domainName := getDomainName()

	return &protos.GetGatewayRegistrationInfoResponse{
		Response: &protos.GetGatewayRegistrationInfoResponse_GatewayRegistrationInfo{
			GatewayRegistrationInfo: &protos.GatewayRegistrationInfo{
				Gateway: &protos.GatewayInfo{
					NetworkId: tokenInfo.Gateway.NetworkId,
					LogicalId: tokenInfo.Gateway.LogicalId,
				},
				RootCA:     rootCA,
				DomainName: domainName,
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
		Gateway: &protos.GatewayInfo{
			NetworkId: networkID,
			LogicalId: &protos.AccessGatewayID{
				Id: logicalID,
			},
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
	body, err := ioutil.ReadFile("/var/opt/magma/certs/rootCA.pem")
	if err != nil {
		return "", status.Errorf(codes.Internal, "Error reading rootCA.pem file: %v", err)
	}
	return string(body), nil
}

// TODO(#10437)
func getDomainName() string {
	return "localhost"
}
