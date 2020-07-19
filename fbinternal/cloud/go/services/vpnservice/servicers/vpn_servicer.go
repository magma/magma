package servicers

import (
	"fmt"
	"io/ioutil"
	"time"

	fbprotos "magma/fbinternal/cloud/go/protos"
	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/services/certifier"
	certprotos "magma/orc8r/cloud/go/services/certifier/protos"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/ptypes/duration"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VPNServicer struct {
	taKeyPath string
}

type clientType int32

const (
	// Types of VPN clients
	clientTypeInvalid = -1
	clientTypeGateway = 0
	clientTypeSupport = 1

	// Duration for which certs are valid for support team / gateway
	supportTime = 12 * time.Hour            // Short-lived
	gwTime      = 10 * 365 * 24 * time.Hour // Seems like right now it's 10 years.
)

func NewVPNServicer(taKeyPath string) (srv *VPNServicer) {
	return &VPNServicer{taKeyPath: taKeyPath}
}

// Return the signing certificate
func (srv *VPNServicer) GetCA(context.Context, *protos.Void) (*protos.CACert, error) {
	getCAReq := &certprotos.GetCARequest{CertType: protos.CertType_VPN}
	caFromCertifier, err := certifier.GetCACert(getCAReq)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve VPN CA: %s", err)
	}
	caMsg := &protos.CACert{Cert: caFromCertifier.Cert}
	return caMsg, nil
}

func (srv *VPNServicer) RequestPSK(ctx context.Context, req *protos.Void) (*fbprotos.PSK, error) {
	// Get gateway id from context
	gw := protos.GetClientGateway(ctx)
	if gw == nil {
		return nil, status.Errorf(
			codes.PermissionDenied, "Missing Gateway Identity")
	}
	if !gw.Registered() {
		return nil, status.Errorf(
			codes.PermissionDenied, "Gateway is not registered")
	}
	taKey, err := loadTAKey(srv.taKeyPath)
	if err != nil {
		return nil, fmt.Errorf("err loading PSK: %v", err)
	}
	return &fbprotos.PSK{TaKey: taKey}, nil
}

func loadTAKey(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// Authenticate and return client type from calling context
func getClientType(ctx context.Context) (clientType, error) {
	id := protos.GetClientIdentity(ctx)
	if identity.IsGateway(id) {
		gw := id.GetGateway()
		if gw == nil {
			return clientTypeInvalid, fmt.Errorf("missing gateway identity")
		}
		if !gw.Registered() {
			return clientTypeInvalid, fmt.Errorf("gateway is not registered")
		}
		return clientTypeGateway, nil

	} else if identity.IsOperator(id) {
		operator := id.GetOperator()
		if len(operator) == 0 {
			return clientTypeInvalid, fmt.Errorf("missing operator identity")
		}
		return clientTypeSupport, nil

	} else {
		return clientTypeInvalid, fmt.Errorf("identity must be of gateway or operator")
	}
}

// Authenticate & decide valid time based on identity of requester
func getValidDuration(client clientType) (int64, error) {
	switch client {
	case clientTypeGateway:
		return int64(gwTime.Seconds()), nil
	case clientTypeSupport:
		return int64(supportTime.Seconds()), nil
	default:
		return 0, fmt.Errorf("invalid client type")
	}
}

// Get identity from context and return a signed certificate
func (srv *VPNServicer) RequestCert(
	ctx context.Context,
	req *fbprotos.VPNCertRequest,
) (*fbprotos.VPNCertificate, error) {
	if req == nil {
		return nil, fmt.Errorf("no cert request given (nil argument)")
	}

	clientTyp, err := getClientType(ctx)
	if err != nil {
		return nil, fmt.Errorf("error authenticating client: %s", err)
	}

	validTime, err := getValidDuration(clientTyp)
	if err != nil {
		return nil, fmt.Errorf("error setting cert valid duration: %s", err)
	}

	csr := &protos.CSR{
		Id:        protos.GetClientIdentity(ctx),
		CertType:  protos.CertType_VPN,
		ValidTime: &duration.Duration{Seconds: validTime, Nanos: 0},
		CsrDer:    req.Request,
	}

	signedCert, err := certifier.SignCSR(csr)
	if err != nil {
		return nil, fmt.Errorf("error signing vpn cert request: %s", err)
	}

	signedCertMsg := &fbprotos.VPNCertificate{
		Serial: signedCert.Sn.Sn,
		Cert:   signedCert.CertDer,
	}
	return signedCertMsg, nil
}
