/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// gw_to_feg_relay is h2c & GRPC server serving requests from AGWs to FeG
package gw_to_feg_relay

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/serdes"
	"magma/feg/cloud/go/services/feg/obsidian/models"
	"magma/feg/cloud/go/services/health"
	"magma/orc8r/cloud/go/http2"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service/middleware/unary"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
	"magma/orc8r/lib/go/protos"
)

const (
	// Client Certificate CN Header
	ClientCertCnKeyHttpHeader = "X-Magma-Client-Cert-Cn"
	// Client Certificate Serial Number Header
	ClientCertSnKeyHttpHeader = "X-Magma-Client-Cert-Serial"
)

// GatewayToFeGServer is a relay from Gateway to Feg. It is a http/2 server
// serving grpc requests. It authenticates the sender of the request,
// then forward request to FeG, and forward response back to caller.
type GatewayToFeGServer struct {
	*http2.H2CServer
	client *http2.H2CClient
}

// NewGatewayToFegServer creates a new GatewayToFegServer
func NewGatewayToFegServer() *GatewayToFeGServer {
	return &GatewayToFeGServer{
		H2CServer: http2.NewH2CServer(),
		client:    http2.NewH2CClient(),
	}
}

// Run blocks and runs the server on addr.
func (server *GatewayToFeGServer) Run(addr string) {
	server.H2CServer.Run(addr, server.useDispatcherHandler)
}

// Serve blocks and serves on listener. This is intended to be used
// by testing.
func (server *GatewayToFeGServer) Serve(listener net.Listener) {
	server.H2CServer.Serve(listener, server.useDispatcherHandler)
}

func (server *GatewayToFeGServer) useDispatcherHandler(responseWriter http.ResponseWriter, req *http.Request) {
	// check calling gateway's identity through certifier
	gw, err := getGatewayIdentity(req.Header)
	if err != nil || gw == nil {
		glog.Errorf(err.Error())
		http2.WriteErrResponse(responseWriter, http2.NewHTTPGrpcError(
			"Missing gateway identity",
			int(codes.PermissionDenied), http.StatusForbidden))
		return
	}
	if !gw.Registered() {
		http2.WriteErrResponse(responseWriter, http2.NewHTTPGrpcError(
			"Gateway is not registered",
			int(codes.PermissionDenied), http.StatusPreconditionFailed))
		return
	}
	http2.LogRequestWithVerbosity(req, 4)
	// get destination feg hwId
	fegHwId, err := getFeGHwIdForNetwork(gw.NetworkId)
	if err != nil {
		glog.Errorf(err.Error())
		http2.WriteErrResponse(responseWriter,
			http2.NewHTTPGrpcError(err.Error(), int(codes.NotFound),
				http.StatusBadRequest))
		return
	}
	// get dispatcher's http server addr
	addr, addrErr := getDispatcherHttpServerAddr(fegHwId)
	if addrErr != nil {
		glog.Errorf(addrErr.Error())
		http2.WriteErrResponse(responseWriter,
			http2.NewHTTPGrpcError(addrErr.Error(), int(codes.Unavailable),
				http.StatusBadRequest))
		return
	}
	glog.V(4).Infof("Forwarding request to FeG %s via %s", fegHwId, addr)
	// create request to dispatcher http server
	newReq, newReqErr := createNewRequest(req, addr, fegHwId)
	if newReqErr != nil {
		glog.Errorf(newReqErr.Error())
		http2.WriteErrResponse(responseWriter, newReqErr)
		return
	}
	// forward request to dispatcher http server
	resp, relayErr := server.client.Do(newReq)
	if relayErr != nil {
		glog.Error(relayErr.Error())
		http2.WriteErrResponse(responseWriter,
			http2.NewHTTPGrpcError(relayErr.Error(), int(codes.Unavailable),
				http.StatusBadRequest))
		return
	}
	// process response
	respErr := processResponse(responseWriter, resp)
	if respErr != nil {
		glog.Errorf(respErr.Error())
		http2.WriteErrResponse(responseWriter, respErr)
	}
}

func createNewRequest(req *http.Request, addr, hwId string) (*http.Request, *http2.HTTPGrpcError) {
	// clone request so that we can modify it freely
	newReq := new(http.Request)
	*newReq = *req
	// request url cannot be set by the client
	newReq.RequestURI = ""
	newReq.URL.Scheme = "http"
	newReq.URL.Host = addr
	newReq.Header.Set(gateway_registry.GatewayIdHeaderKey, hwId)
	// use service name as authority
	if auth := strings.Split(req.Host, "-"); len(auth) > 0 {
		newReq.Host = auth[0]
	}
	glog.V(4).Info("Cloned request from GW:")
	http2.LogRequestWithVerbosity(newReq, 4)
	return newReq, nil
}

func getFeGHwIdForNetwork(agNwID string) (string, error) {
	cfg, err := configurator.LoadNetworkConfig(agNwID, feg.FederatedNetworkType, serdes.Network)
	if err != nil {
		return "", fmt.Errorf("could not load federated network configs for access network %s: %s", agNwID, err)
	}
	federatedConfig, ok := cfg.(*models.FederatedNetworkConfigs)
	if !ok || federatedConfig == nil {
		return "", fmt.Errorf("invalid federated network config found for network: %s", agNwID)
	}
	if federatedConfig.FegNetworkID == nil || *federatedConfig.FegNetworkID == "" {
		return "", fmt.Errorf("FegNetworkID is empty in network config of network: %s", agNwID)
	}
	fegCfg, err := configurator.LoadNetworkConfig(*federatedConfig.FegNetworkID, feg.FegNetworkType, serdes.Network)
	if err != nil || fegCfg == nil {
		return "", fmt.Errorf("unable to retrieve config for federation network: %s", *federatedConfig.FegNetworkID)
	}
	networkFegConfigs, ok := fegCfg.(*models.NetworkFederationConfigs)
	if !ok || networkFegConfigs == nil {
		return "", fmt.Errorf("invalid federation network config found for network: %s", *federatedConfig.FegNetworkID)
	}
	servedNetworkIDs := networkFegConfigs.ServedNetworkIds
	for _, network := range servedNetworkIDs {
		if agNwID == network {
			return getActiveFeGForNetwork(*federatedConfig.FegNetworkID)
		}
	}
	return "", fmt.Errorf("federation network %s is not configured to serve network: %s", *federatedConfig.FegNetworkID, agNwID)
}

func getActiveFeGForNetwork(fegNetworkID string) (string, error) {
	activeGW, err := health.GetActiveGateway(fegNetworkID)
	if err != nil {
		return "", fmt.Errorf("Unable to retrieve active FeG for network: %s; %s", fegNetworkID, err)
	}
	hardwareID, err := configurator.GetPhysicalIDOfEntity(fegNetworkID, orc8r.MagmadGatewayType, activeGW)
	if err != nil {
		return "", fmt.Errorf("Unable to retrieve hardware ID for active feg: %s in network: %s; %s", activeGW, fegNetworkID, err)
	}
	if len(hardwareID) == 0 {
		return "", fmt.Errorf("Unable to retrieve Hardware ID for active feg: %s in network: %s", activeGW, fegNetworkID)
	}
	return hardwareID, nil
}

func getDispatcherHttpServerAddr(hwId string) (string, error) {
	addr, err := gateway_registry.GetServiceAddressForGateway(hwId)
	if err != nil {
		return "", err
	}
	return addr, nil
}

func getGatewayIdentity(headers http.Header) (*protos.Identity_Gateway, error) {
	if headers == nil {
		return nil, fmt.Errorf("no headers present to determine gateway identity")
	}
	md := map[string]string{}
	snKeys := headers[ClientCertSnKeyHttpHeader]
	if len(snKeys) == 1 {
		md[ClientCertSnKeyHttpHeader] = snKeys[0]
	} else {
		return nil, fmt.Errorf("no client cert header present to determine gateway identity")
	}
	cnKeys := headers[ClientCertCnKeyHttpHeader]
	if len(cnKeys) == 1 {
		md[ClientCertCnKeyHttpHeader] = cnKeys[0]
	}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(md))
	newCtx, _, _, err := unary.SetIdentityFromContext(ctx, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting gateway identity. err: %v\n", err)
	}
	gw := protos.GetClientGateway(newCtx)
	return gw, nil
}

func processResponse(w http.ResponseWriter, resp *http.Response) *http2.HTTPGrpcError {
	if resp == nil {
		errMsg := "nil http response"
		return http2.NewHTTPGrpcError(errMsg, int(codes.Internal), http.StatusInternalServerError)
	}
	glog.V(5).Infof(
		"response header: %v\n response trailer before reading body: %v\n response statusCode: %v\n",
		resp.Header, resp.Trailer, resp.StatusCode,
	)
	writeHeadersToResponseWriter(resp.Header, w)
	payload, err := getRespPayload(resp.Body)
	if err != nil || payload == nil {
		return http2.NewHTTPGrpcError(
			fmt.Sprintf("failed to forward response payload. err: %v",
				err),
			int(codes.Internal), http.StatusInternalServerError)
	}
	glog.V(5).Infof("Response trailer after reading body: %v\n", resp.Trailer)
	w.WriteHeader(resp.StatusCode)
	w.Write(payload)
	return nil
}

func getRespPayload(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	return ioutil.ReadAll(body)
}

func writeHeadersToResponseWriter(headers http.Header,
	w http.ResponseWriter) {
	// see how to write trailers: https://golang.org/pkg/net/http/#example_ResponseWriter_trailers
	w.Header().Set("Trailer", "Grpc-Status")
	w.Header().Add("Trailer", "Grpc-Message")
	for k, vals := range headers {
		for _, val := range vals {
			if len(w.Header().Get(k)) == 0 {
				w.Header().Set(k, val)
			} else {
				w.Header().Add(k, val)
			}
		}
	}
}
