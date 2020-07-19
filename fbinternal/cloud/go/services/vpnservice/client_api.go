package vpnservice

import (
	fbprotos "magma/fbinternal/cloud/go/protos"
	"magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

const ServiceName = "VPNSERVICE"

// Utility function to get a RPC connection to the VPN service
func getVPNServiceClient() (fbprotos.VPNServiceClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := errors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}

	return fbprotos.NewVPNServiceClient(conn), err
}

// Get the VPN CA certificate
func GetVPNCA() (cert []byte, err error) {
	client, err := getVPNServiceClient()
	if err != nil {
		return nil, err
	}

	caMsg, err := client.GetCA(context.Background(), &protos.Void{})
	if err != nil {
		glog.Errorf("Failed to get VPN CA: %s", err)
		return nil, err
	}

	return caMsg.Cert, nil
}

// Return a certificate signed by the VPN CA, with serial number.
// CSR is in ASN.1 DER encoding.
func RequestSignedCert(csr []byte) (sn string, cert []byte, err error) {
	client, err := getVPNServiceClient()
	if err != nil {
		return "", nil, err
	}

	req := &fbprotos.VPNCertRequest{Request: csr}
	certMsg, err := client.RequestCert(context.Background(), req)
	if err != nil {
		glog.Errorf("Failed to retrieve signed VPN cert: %s", err)
		return "", nil, err
	}
	return certMsg.Serial, certMsg.Cert, nil
}
