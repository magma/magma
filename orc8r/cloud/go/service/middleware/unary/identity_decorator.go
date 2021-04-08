/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// package interceptors implements all cloud service framework unary interceptors
package unary

import (
	"net"
	"time"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/certifier"
	certprotos "magma/orc8r/cloud/go/services/certifier/protos"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/lib/go/metrics"
	"magma/orc8r/lib/go/protos"
	unarylib "magma/orc8r/lib/go/service/middleware/unary"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// SetIdentityFromContext is an identity decorator implements Identity injector
// for all authenticated requests.
// It looks for x-magma-client-cert-cn and x-magma-client-cert-serial HTTP headers
// in the context, verifies validity of the client certificate and injects
// a valid, verified client Identity into RPC context
// SetIdentityFromContext can only modify CTX, it doesn't affect other RPC
// parameters
const (
	// Client Certificate CN Header
	CLIENT_CERT_CN_KEY = "x-magma-client-cert-cn"
	// Client Certificate Serial Number Header
	CLIENT_CERT_SN_KEY = "x-magma-client-cert-serial"
)

const (
	ERROR_MSG_NO_METADATA      = "Missing Required CTX Metadata"
	ERROR_MSG_INVALID_CERT     = "Invalid Client Certificate"
	ERROR_MSG_UNKNOWN_CERT     = "Unknown Client Certificate"
	ERROR_MSG_EXPIRED_CERT     = "Expired Client Certificate"
	ERROR_MSG_MISSING_IDENTITY = "Missing Certificate Identity"
	ERROR_MSG_INVALID_TYPE     = "Invalid Certificate Owner"
	ERROR_MSG_UNKNOWN_CLIENT   = "Unknown Client Address"

	// GW should start bootstrap 20 hours prior to cert expiration, give it 10 hours to try & start counting
	// bootstrap failures after that
	CERT_EXPIRATION_DURATION_THRESHOLD = time.Hour * 10
)

var gwExpiringCert = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "gateway_expiring_cert",
		Help: "Count of GW cloud requests with soon to expire Client Certificate (indicated GW bootstrapper failure",
	},
	[]string{metrics.NetworkLabelName, metrics.GatewayLabelName},
)

func init() {
	prometheus.MustRegister(gwExpiringCert)
}

// SetIdentityFromContext finds Identity associated with caller's Client
// Certificate Serial Number (if present), makes sure that the found Identity
// is of a Gateway & fills in all available Gateway Identity information
// SetIdentityFromContext will bypass the Identity checks for local callers
// (other services on the cloud) and allowlisted RPCs (methods in
// identityDecoratorBypassList)
func SetIdentityFromContext(ctx context.Context, _ interface{}, info *grpc.UnaryServerInfo) (newCtx context.Context, newReq interface{}, resp interface{}, err error) {
	// There are 5 possible outcomes:
	// 1. !ok -> type assertion: mdIncomingKey{} is present, but it's not of MD type
	//    It should never happen & possibly indicates a hacking attempt -> reject
	//    request
	// 2. ctxMetadata == nil -> same as case #1, should never happen -> reject
	//    request
	// 3. ctxMetadata.Len() is 0: potentially possible for internal service 2
	//    service calls -> accept request
	// 4. x-magma-client-cert-serial is not present -> possible for internal
	//    service to service calls -> accept request
	// 5. x-magma-client-cert-serial is present -> external request, continue
	//    verification below

	ctxMetadata, ok := metadata.FromIncomingContext(ctx)
	if !ok || ctxMetadata == nil {
		glog.Info(ERROR_MSG_NO_METADATA)
		// Metadata should always be present for GRPC client calls
		// If we want to enable ANY calls from local clients, we need to
		// change the return statement to:
		//   return newCtx, newReq, resp, ensureLocalPeer(ctx)
		// but, it may present a se3curity risk
		err = status.Error(codes.Unauthenticated, ERROR_MSG_NO_METADATA)
		return newCtx, newReq, resp, err
	}

	// First, try to find the caller's identity
	snlist, snok := ctxMetadata[CLIENT_CERT_SN_KEY]
	if snok { // there is a certificate serial number (CSN) list in CTX
		if len(snlist) != 1 {
			// there can be only one CSN, error out if not
			glog.Infof("Multiple CSNs found in metadata: %+v", ctxMetadata)
			err = status.Error(codes.Unauthenticated, "Multiple CSNs present")
		} else {
			// One CSN is found, find Identity associated with it
			var gwIdentity *protos.Identity
			var certExpTime int64
			// Check if SN is the reserved value used for all inter-orc8r calls
			if snlist[0] == unarylib.ORC8R_CLIENT_CERT_VALUE {
				return newCtx, newReq, resp, nil
			}
			gwIdentity, certExpTime, err = findGatewayIdentity(snlist[0], ctxMetadata)
			if err == nil {
				// If a valid GW Identity is found, add it into CTX for use
				// by the callee
				newCtx = gwIdentity.NewContextWithIdentity(protos.NewContextWithCertExpiration(ctx, certExpTime))
				return newCtx, newReq, resp, err
			}
		}
	} else {
		// No CSNs found, check there is also no Certificate Common Names
		if _, ok = ctxMetadata[CLIENT_CERT_CN_KEY]; ok {
			// CN header is present while SN header is missing - possible
			// security hack, either both or neither of the headers should be
			// set
			glog.Infof("CCN is present without SCN in metadata: %+v", ctxMetadata)
			err = status.Error(codes.Unauthenticated, "Inconsistent Request Signature")
		}
	}

	if info != nil {
		// Check if the call is for a allowlisted method - anything is allowed
		// do this check past possible identity decoration to still allow to add
		// valid identity even to allowlisted requests
		if _, ok := identityDecoratorBypassList[info.FullMethod]; ok {
			// Bypass method (Bootstrapper & Co.), shortcut...
			return newCtx, newReq, resp, nil
		}
	}

	// Only allow local clients if there is no previous errors (snok == false)
	if err == nil {
		// We assume that only external calls forwarded by cloud proxy (or unit
		// tests) will have CSN & CCN headers set. The absence of the headers
		// along with client IP verification will indicate a local service to
		// service or Obsidian to service call
		// For internal calls, no identity verification needed, just make sure
		// it's a local client
		err = ensureLocalPeer(ctx)
		if err != nil {
			var rpc string
			if info != nil {
				rpc = info.FullMethod
			} else {
				rpc = "Undefined"
			}
			glog.Infof("Empty CTX Metadata from non-local %s client: %v", rpc, err)
		}
	}
	return newCtx, newReq, resp, err
}

// findGatewayIdentity returns 'decorated' Gateway Identity corresponding to the
// given certificate serialNumber and it's certificate expiration time in Unix time seconds
// The Identity is 'decorated' with all information that can be gathered about
// the given GW's Hardware Id, such as network ID & Logical ID. At a minimum -
// the returned Identity should have a valid, verified via Certifier HW ID.
// If the target PRC needs Network and/or logical ID, the service should handle
// their absence for unregistered Gateways and return an error.
// The identity middleware only ensures that GW is who it says it is (HwID)
func findGatewayIdentity(serialNumber string, md metadata.MD) (*protos.Identity, int64, error) {
	certInfo, err := getCertifierIinfo(serialNumber, md)
	if err != nil {
		return nil, 0, err
	}

	id := certInfo.GetId()
	gwIdentity := id.GetGateway()
	if gwIdentity == nil {
		glog.Infof("Identity (%s) of Cert SN %s from metadata %+v is not a Gateway", id.HashString(), serialNumber, md)
		return nil, 0, status.Error(codes.PermissionDenied, ERROR_MSG_INVALID_TYPE)
	}

	// At this point we should have a valid GW Identity with HardwareId, so
	// the Gateway is authenticated

	entity, err := configurator.LoadEntityForPhysicalID(gwIdentity.HardwareId, configurator.EntityLoadCriteria{}, serdes.Entity)
	if err != nil {
		glog.Infof("Unregistered Gateway Id: %s for Cert SN: %s; err: %s; metadata: %+v", gwIdentity.HardwareId, serialNumber, err, md)
	}
	networkID := entity.NetworkID
	logicalID := entity.Key

	expiration, _ := ptypes.Timestamp(certInfo.GetNotAfter())
	expSeconds := expiration.Unix()

	if expiration.Sub(clock.Now()) < CERT_EXPIRATION_DURATION_THRESHOLD {
		gwExpiringCert.WithLabelValues(networkID, logicalID).Inc()
	}

	return identity.NewGateway(gwIdentity.HardwareId, networkID, logicalID), expSeconds, nil
}

// getCertifierIdentity retrieves 'raw' identity associated with the Certificate
// SerialNumber from certifier
func getCertifierIinfo(serialNumber string, md metadata.MD) (*certprotos.CertificateInfo, error) {
	// Call Certifier & get the Identity from it
	// & error out if SN is not found or expired
	certInfo, err := certifier.GetCertificateIdentity(serialNumber)
	if err != nil {
		glog.Infof("Lookup error '%s' for Cert SN: %s, metadata: %+v", err, serialNumber, md)
		return nil, status.Error(codes.PermissionDenied, ERROR_MSG_UNKNOWN_CERT)
	}
	if certInfo == nil {
		glog.Infof("Missing Certificate Info for Cert SN: %s, metadata: %+v", serialNumber, md)
		return nil, status.Error(codes.PermissionDenied, ERROR_MSG_INVALID_CERT)
	}
	// Check if certificate time is not expired/not active yet
	err = certifier.VerifyDateRange(certInfo)
	if err != nil {
		glog.Infof("Certificate Validation Error '%s' for Cert SN: %s, metadata: %+v", err, serialNumber, md)
		return nil, status.Error(codes.PermissionDenied, ERROR_MSG_EXPIRED_CERT)
	}
	if certInfo.Id == nil {
		glog.Infof("Missing Gateway ID for Cert SN: %s, metadata: %+v", serialNumber, md)
		return nil, status.Error(codes.PermissionDenied, ERROR_MSG_MISSING_IDENTITY)
	}

	return certInfo, nil
}

// ensureLocalPeer retrieves & parses caller address and verifies that it's
// local (loopback)
// returns an error if it's missing, invalid or not a local address
func ensureLocalPeer(ctx context.Context) error {
	caller, peerok := peer.FromContext(ctx)
	if !peerok || caller == nil {
		return status.Error(codes.PermissionDenied, ERROR_MSG_UNKNOWN_CLIENT)
	}
	host, _, err := net.SplitHostPort(caller.Addr.String())
	if err != nil {
		host = caller.Addr.String()
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return status.Errorf(codes.PermissionDenied, "Invalid Client Address: %+v", caller.Addr)
	}
	if !ip.IsLoopback() {
		return status.Errorf(codes.PermissionDenied, "Missing Client Certificate from Client %s", ip.String())
	}
	return nil
}
