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

// Package httpserver is a http/2 h2c server. It's run within the same process
// as SyncRPC grpc servicer.
//
// When a client wants to send a grpc request to some service on
// gateway with hardwareId hwId, it identifies the addr of the SyncRPC grpc
// servicer to which the gateway has a bidirectional stream to, and sends a grpc
// request to the httpServer that's within the same process of that grpc servicer.
//
// This httpServer converts httpRequest to GatewayRequest, send it over to grpc
// servicer using GatewayRPCBroker, waits for a response, and converts the
// GatewayResponse to a HttpResponse and send it back to the client.
package httpserver

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"magma/orc8r/cloud/go/http2"
	"magma/orc8r/cloud/go/services/dispatcher/broker"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
)

const (
	DefaultHttpResponseStatus = 200

	responseTimeoutSecs = 15
	maxCancelAttempts   = 5
)

type SyncRPCHttpServer struct {
	*http2.H2CServer
	broker broker.GatewayRPCBroker
}

func NewSyncRPCHttpServer(broker broker.GatewayRPCBroker) *SyncRPCHttpServer {
	return &SyncRPCHttpServer{http2.NewH2CServer(), broker}
}

func (server *SyncRPCHttpServer) Run(addr string) {
	server.H2CServer.Run(addr, server.rootHandler)
}

func (server *SyncRPCHttpServer) Serve(listener net.Listener) {
	server.H2CServer.Serve(listener, server.rootHandler)
}

func (server *SyncRPCHttpServer) rootHandler(responseWriter http.ResponseWriter, req *http.Request) {
	http2.LogRequestWithVerbosity(req, 4)
	respChan, err := server.sendRequest(req)
	if err != nil {
		glog.Errorf(err.Msg)
		// Also write to client.
		http2.WriteErrResponse(responseWriter, err)
		return
	}

	// Wait for response or timeout.
	for {
		select {
		case gwResponse := <-respChan:
			err := processResponse(responseWriter, gwResponse)
			if err != nil {
				glog.Errorf(err.Msg)
				http2.WriteErrResponse(responseWriter, err)
			}
			if isResponseComplete(responseWriter) {
				return
			}
		case <-time.After(time.Second * responseTimeoutSecs):
			http2.WriteErrResponse(
				responseWriter,
				http2.NewHTTPGrpcError("Request timed out", int(codes.DeadlineExceeded), http.StatusRequestTimeout),
			)
			return
		}
	}
}

// sendRequest sends a SyncRPCRequest to the gateway and creates
// a goroutine to notify the gateway when the context is done.
func (server *SyncRPCHttpServer) sendRequest(req *http.Request) (chan *protos.GatewayResponse, *http2.HTTPGrpcError) {
	gwReq, err := createRequest(req)
	if err != nil {
		return nil, err
	}
	gwRespChannel, sendReqErr := server.broker.SendRequestToGateway(gwReq)
	if sendReqErr != nil {
		errMsg := fmt.Sprintf("err sending request %v to gateway: %v", gwReq, sendReqErr)
		return nil, http2.NewHTTPGrpcError(errMsg, int(codes.Internal), http.StatusInternalServerError)
	}

	// If context is done (connection between client and this HTTP/2 server is closed),
	// notify the proxy client in gateway to stop receiving frames
	go func() {
		<-req.Context().Done()
		// Attempt to cancel the request up to maxCancelAttempts times. This is
		// because CancelGatewayRequest may fail if the request queue is full,
		// in which case the cancel message is not added to the request queue.
		// Therefore, if we fail to enqueue the cancel message, sleep for a couple
		// seconds and retry.
		for cancelAttempts := 0; cancelAttempts < maxCancelAttempts; cancelAttempts++ {
			err := server.broker.CancelGatewayRequest(gwReq.GwId, gwRespChannel.ReqId)
			if err == nil {
				return
			}

			time.Sleep(5 * time.Second)
		}
		glog.Errorf("Could not cancel gateway request after %v attempts", maxCancelAttempts)
	}()

	return gwRespChannel.RespChan, nil
}

// createRequest converts a HTTP request to a GatewayRequest.
func createRequest(req *http.Request) (*protos.GatewayRequest, *http2.HTTPGrpcError) {
	headers := req.Header
	gwIds := headers[gateway_registry.GatewayIdHeaderKey]
	if len(gwIds) == 0 || len(gwIds[0]) == 0 {
		return nil, http2.NewHTTPGrpcError("No Gatewayid provided in metaData", int(codes.InvalidArgument), http.StatusBadRequest)
	}
	gwId := gwIds[0]
	delete(headers, gateway_registry.GatewayIdHeaderKey)
	authority, err := getAuthority(req.Host)
	if err != nil {
		return nil, err
	}
	path, err := getPath(req.URL)
	if err != nil {
		return nil, err
	}
	body, err := getPayload(req.Body)
	if err != nil {
		return nil, err
	}
	gwReq := &protos.GatewayRequest{
		GwId:      gwId,
		Authority: authority,
		Path:      path,
		Headers:   convertHeadersForProto(headers),
		Payload:   body,
	}
	return gwReq, nil
}

func getAuthority(host string) (string, *http2.HTTPGrpcError) {
	if len(host) == 0 {
		return "", http2.NewHTTPGrpcError("No authority provided", int(codes.InvalidArgument), http.StatusBadRequest)
	} else {
		return host, nil
	}
}

func getPath(url *url.URL) (string, *http2.HTTPGrpcError) {
	if url == nil || len(url.Path) == 0 {
		return "", http2.NewHTTPGrpcError("No url path provided", int(codes.InvalidArgument),
			http.StatusBadRequest)
	}
	return url.Path, nil
}

func getPayload(body io.ReadCloser) ([]byte, *http2.HTTPGrpcError) {
	payload, err := ioutil.ReadAll(body)
	defer body.Close()
	if err != nil {
		errMsg := fmt.Sprintf("err reading req body: %v", err)
		return nil, http2.NewHTTPGrpcError(errMsg, int(codes.InvalidArgument), http.StatusBadRequest)
	}
	return payload, nil
}

func convertHeadersForProto(headers http.Header) map[string]string {
	ret := make(map[string]string)
	for k, vals := range headers {
		ret[k] = concatenateHeaders(vals)
	}
	return ret
}

func concatenateHeaders(headers []string) string {
	return strings.Join(headers, ",")
}

// processResponse converts the GatewayResponse to an HTTP response and sends it
// back to the client.
func processResponse(w http.ResponseWriter, gwResp *protos.GatewayResponse) *http2.HTTPGrpcError {
	if gwResp.KeepConnActive {
		return nil
	}
	if gwResp == nil {
		// Remains for backward compatibility, but it shouldn't get forwarded
		// a nil GatewayResponse in new versions.
		return http2.NewHTTPGrpcError("nil GatewayResponse", int(codes.Internal), http.StatusInternalServerError)
	}
	if gwResp.Err != "" {
		return http2.NewHTTPGrpcError(gwResp.Err, int(codes.Internal), http.StatusInternalServerError)
	}
	headers := gwResp.GetHeaders()
	writeHeadersToResponse(headers, w)
	httpStatus, err := getHttpStatusFromGatewayResponse(gwResp.Status)
	w.WriteHeader(httpStatus)
	if gwResp.Payload != nil {
		w.Write(gwResp.Payload)
	}
	w.(http.Flusher).Flush()
	if err != nil {
		// Only log, and do not send to client.
		glog.Errorf("%v\n", err)
	}
	return nil
}

func getHttpStatusFromGatewayResponse(gwRespStatus string) (int, error) {
	httpStatus, err := strconv.Atoi(gwRespStatus)
	if err != nil {
		return DefaultHttpResponseStatus, fmt.Errorf("cannot parse status of gatewayResponse: %v\n", err)
	}
	// invalid status code, defaults to 200
	if statusText := http.StatusText(httpStatus); len(statusText) == 0 {
		return DefaultHttpResponseStatus, fmt.Errorf("Unrecognized httpStatus: %v\n", httpStatus)
	}
	return httpStatus, nil
}

func writeHeadersToResponse(headers map[string]string, w http.ResponseWriter) {
	// see how to write trailers: https://golang.org/pkg/net/http/#example_ResponseWriter_trailers
	w.Header().Set("Trailer", "Grpc-Status, Grpc-Message")
	for k, v := range headers {
		vals := strings.Split(v, ",")
		for _, val := range vals {
			w.Header().Add(k, val)
		}
	}
}

func isResponseComplete(w http.ResponseWriter) bool {
	return len(w.Header().Get("Grpc-Status")) != 0
}
