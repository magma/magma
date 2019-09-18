/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// h2c server serving requests from FeG to AG.
package gw_to_feg_relay

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/plugin/models"
	"magma/feg/cloud/go/services/health"
	"magma/lte/cloud/go/lte"
	ltemodels "magma/lte/cloud/go/plugin/models"
	"magma/orc8r/cloud/go/http2"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/service/middleware/unary"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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
	if hwId, err := getFeGHwIdEnvOverride(); err == nil {
		glog.Infof("Using FEG_HWID env variable: %s for feg_relay", hwId)
	}
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

func (server *GatewayToFeGServer) useDispatcherHandler(
	responseWriter http.ResponseWriter, req *http.Request,
) {
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
			int(codes.PermissionDenied), http.StatusForbidden))
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
	return
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
	http2.LogRequestWithVerbosity(newReq, 4)
	return newReq, nil
}

func getFeGHwIdEnvOverride() (string, error) {
	hwId, ok := os.LookupEnv("FEG_HWID")
	if ok && len(hwId) > 0 {
		return hwId, nil
	}
	return "", fmt.Errorf("Environment variable FEG_HWID is unset")
}

func getFeGHwIdForNetwork(agNwId string) (string, error) {
	fegEnvHwID, err := getFeGHwIdEnvOverride()
	if err == nil {
		return fegEnvHwID, nil
	}
	cfg, err := configurator.GetNetworkConfigsByType(agNwId, lte.CellularNetworkType)
	if err != nil || cfg == nil {
		return "", fmt.Errorf("Unable to retrieve cellular config for AG network: %s", agNwId)
	}
	cellularConfig, ok := cfg.(*ltemodels.NetworkCellularConfigs)
	if !ok {
		return "", fmt.Errorf("Invalid cellular config found for AG network: %s", agNwId)
	}
	fegNetworkID := cellularConfig.FegNetworkID
	if fegNetworkID == "" {
		return "", fmt.Errorf("FegNetworkID is not set in cellular config for network: %s", agNwId)
	}

	fegCfg, err := configurator.GetNetworkConfigsByType(string(fegNetworkID), feg.FegNetworkType)
	if err != nil || fegCfg == nil {
		return "", fmt.Errorf("Unable to retrieve config for FeG network: %s", fegNetworkID)
	}
	fegNetworkConfig, ok := fegCfg.(*models.NetworkFederationConfigs)
	if !ok {
		return "", fmt.Errorf("Invalid feg config found for FeG network: %s", fegNetworkID)
	}

	servedNetworkIDs := fegNetworkConfig.ServedNetworkIds
	for _, network := range servedNetworkIDs {
		if agNwId == network {
			return getActiveFeGForNetwork(string(fegNetworkID))
		}
	}
	return "", fmt.Errorf("Federated Gateway Network: %s is not configured to serve network: %s", fegNetworkID, agNwId)
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
	if snKeys != nil && len(snKeys) == 1 {
		md[ClientCertSnKeyHttpHeader] = snKeys[0]
	} else {
		return nil, fmt.Errorf("no client cert header present to determine gateway identity")
	}
	cnKeys := headers[ClientCertCnKeyHttpHeader]
	if cnKeys != nil && len(cnKeys) == 1 {
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
	glog.V(5).Infof("response header: %v\n response trailer "+
		"before reading body: %v\n response statusCode: %v\n",
		resp.Header, resp.Trailer, resp.StatusCode)
	if resp == nil {
		errMsg := "nil http response"
		return http2.NewHTTPGrpcError(errMsg,
			int(codes.Internal), http.StatusInternalServerError)
	}
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
