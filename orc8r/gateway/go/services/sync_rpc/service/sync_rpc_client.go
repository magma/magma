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

// package service implements the core of bootstrapper
package service

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/http2"

	"magma/gateway/service_registry"
	"magma/orc8r/lib/go/definitions"
	_ "magma/orc8r/lib/go/initflag"
	"magma/orc8r/lib/go/protos"
)

const (
	// grpc message is delivered as a length prefixed message within the HTTP2 DATA
	// frames, with first 5 bytes for compression and msg length
	// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md
	GRPC_MSGLEN_SZ  = 5
	GRPC_LEN_OFFSET = 1

	MinRetryInterval      = time.Millisecond * 500
	MaxRetryInterval      = time.Second * 30
	RetryBackoffIncrement = MinRetryInterval
	SuccessStartInterval  = time.Second * 5

	DefaultSyncRpcHeartbeatInterval = time.Second * 30
	DefaultGatewayKeepaliveInterval = time.Second * 10
	DefaultGatewayResponseTimeout   = time.Second * 120
)

type Config struct {
	// Sync rpc client sends across heartbeat responses with empty reqID
	// in case no messages were sent within this interval(in seconds)
	SyncRpcHeartbeatInterval time.Duration `yaml:"sync_rpc_heartbeat_interval"`

	// Sync rpc client sends across gateway response with KeepConnActive
	// flag set in case it doesn't receive any responses from gateway service
	// within this interval(in seconds)
	GatewayKeepaliveInterval time.Duration `yaml:"gateway_keepalive_interval"`

	// Sync rpc client sends across a gateway timed out error response to
	// dispatcher in case it doesn't receive any responses from gateway service
	// within timeout period(in seconds)
	GatewayResponseTimeout time.Duration `yaml:"gateway_response_timeout"`
}

// Request - outstanding request
type Request struct {
	CancelFunc context.CancelFunc
	terminated bool
}

// SyncRpcClient opens a bidirectional connection with the cloud
type SyncRpcClient struct {
	sync.RWMutex
	// service registry to coonect to the cloud
	serviceRegistry service_registry.GatewayRegistry

	// responseTimeout in seconds
	cfg Config

	// requests which are still being processed
	outstandingReqs map[uint32]*Request

	// channel receiving broker responses
	respCh chan *protos.SyncRPCResponse

	// broker
	broker broker

	// underlying SyncRPC grpc client
	client protos.SyncRPCServiceClient
}

// NewClient returns a new SyncRPC client using given service registry
// If given registry is nil - use shared GW Service Registry config values
func NewClient(reg service_registry.GatewayRegistry) *SyncRpcClient {
	// create a new grpc connection to the dispatcher in the cloud??
	if reg == nil {
		reg = service_registry.Get()
	}
	cfg := &Config{
		SyncRpcHeartbeatInterval: DefaultSyncRpcHeartbeatInterval,
		GatewayKeepaliveInterval: DefaultGatewayKeepaliveInterval,
		GatewayResponseTimeout:   DefaultGatewayResponseTimeout,
	}
	client := &SyncRpcClient{
		serviceRegistry: reg,
		cfg:             *cfg, // copy configs
		outstandingReqs: make(map[uint32]*Request),
		respCh:          make(chan *protos.SyncRPCResponse),
		broker:          newbrokerImpl(cfg),
	}
	return client
}

// Run starts SyncRPC worker loop and block forever
func (c *SyncRpcClient) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for currentBackoffInterval := MinRetryInterval; ; currentBackoffInterval = (currentBackoffInterval + RetryBackoffIncrement) % MaxRetryInterval {
		conn, err := c.serviceRegistry.GetCloudConnection(definitions.DispatcherServiceName)
		if err != nil {
			// TODO
			// Continue for retryable grpc errors
			// Add a delay/jitter for retrying cloud connection for non retryable grpc errors
			glog.Errorf("[SyncRpc] error creating cloud connection: %v", err)
			time.Sleep(currentBackoffInterval)
			continue
		} else {
			glog.Infof("[SyncRpc] successfully connected to cloud '%s' service", definitions.DispatcherServiceName)
		}

		// this should simply wait here for requests and process responses
		// in case we see any error we will retry connecting to the dispatcher
		resetBackoffTime := time.Now().Add(SuccessStartInterval) // make sure, run lasts at least SuccessStartInterval
		c.runSyncRpcClient(ctx, protos.NewSyncRPCServiceClient(conn))
		if time.Now().After(resetBackoffTime) {
			currentBackoffInterval = MinRetryInterval // reset backoff interval
		}
		conn.Close()
		time.Sleep(currentBackoffInterval)
		// exit loop if context is done
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

// run forever
func (c *SyncRpcClient) runSyncRpcClient(ctx context.Context, client protos.SyncRPCServiceClient) error {
	c.client = client
	stream, err := c.client.EstablishSyncRPCStream(ctx)
	if err != nil {
		glog.Errorf("[SyncRPC] Failed establishing SyncRpc stream; error: %v", err)
		return err
	}
	// close connection if error
	errChan := make(chan error)
	go func() {
		stream.Send(&protos.SyncRPCResponse{HeartBeat: true}) // send first heartbeat to establish orc8r queue
		errChan <- c.processStream(ctx, stream)
	}()

	timer := time.NewTimer(c.cfg.SyncRpcHeartbeatInterval)
	defer timer.Stop()
	for {
		select {
		case resp := <-c.respCh:
			if !c.isReqTerminated(resp.ReqId) {
				err := stream.Send(resp)
				if err != nil {
					glog.Errorf("[SyncRpc] send to dispatcher failed: %v", err)
					return err
				}
				timer.Reset(c.cfg.SyncRpcHeartbeatInterval)
				glog.V(3).Infof("[SyncRpc] sent resp: %s", resp)
			} else {
				glog.Errorf("[SyncRpc] request canceled for Id: %d", resp.ReqId)
			}
		case err := <-errChan:
			// error recd handling processStream
			return err
		case <-timer.C:
			err := stream.Send(&protos.SyncRPCResponse{HeartBeat: true})
			if err != nil {
				glog.Error("[SyncRpc] heartbeat to dispatcher failed")
				return err
			}
			glog.V(2).Info("[SyncRpc] heartbeat success")
			timer.Reset(c.cfg.SyncRpcHeartbeatInterval)
		case <-ctx.Done():
			glog.Info("[SyncRPC] Stopping SyncRpcClient")
			return nil
		}
	}
}

func (c *SyncRpcClient) updateTerminatedReqs(reqID uint32) context.CancelFunc {
	c.Lock()
	defer c.Unlock()
	if r, ok := c.outstandingReqs[reqID]; ok {
		r.terminated = true
		return r.CancelFunc
	}
	return nil
}

func (c *SyncRpcClient) isReqTerminated(reqID uint32) bool {
	c.RLock()
	defer c.RUnlock()
	if r, ok := c.outstandingReqs[reqID]; ok {
		return r.terminated
	}
	return false
}

func (c *SyncRpcClient) updateOutstandingdReqs(reqID uint32, cancelFn context.CancelFunc) {
	c.Lock()
	c.outstandingReqs[reqID] = &Request{CancelFunc: cancelFn}
	c.Unlock()
}

func (c *SyncRpcClient) getOutstandingReqCancelFn(reqID uint32) context.CancelFunc {
	c.RLock()
	defer c.RUnlock()
	if r, ok := c.outstandingReqs[reqID]; ok {
		return r.CancelFunc
	}
	return nil
}

func (c *SyncRpcClient) removeOutstandingReq(reqID uint32) {
	c.Lock()
	delete(c.outstandingReqs, reqID)
	c.Unlock()
}

// handleSyncRpcRequest forwards the incoming request to the appropriate destination
func (c *SyncRpcClient) handleSyncRpcRequest(inCtx context.Context, req *protos.SyncRPCRequest) {
	if req == nil {
		glog.Error("[SyncRpc] error empty request received")
		return
	}
	if req.HeartBeat {
		glog.V(3).Info("[SyncRpc] received heartbeat")
		return
	}
	gatewayReq := req.GetReqBody()

	if req.GetConnClosed() {
		glog.V(1).Infof("[SyncRpc] connection closed handling ReqId: %d", req.ReqId)
		cancelFn := c.updateTerminatedReqs(req.ReqId)
		if cancelFn != nil {
			cancelFn()
		} else {
			glog.V(1).Infof("[SyncRpc] closed ReqId: %d is not found", req.ReqId)
		}
		return
	}
	if glog.V(2) {
		glog.Infof("[SyncRpc] request ID %d from GW %s, for service: %s, path: %s",
			req.ReqId, gatewayReq.GetGwId(), gatewayReq.GetAuthority(), gatewayReq.GetPath())
	}
	// get the cancellation fn of the request if present
	// return early if request is outstanding
	if cancelFn := c.getOutstandingReqCancelFn(req.ReqId); cancelFn != nil {
		glog.Warningf("[SyncRpc] duplicate request to %s, %s with Id: %d",
			gatewayReq.GetAuthority(), gatewayReq.GetPath(), req.ReqId)
		c.respCh <- buildSyncRpcErrorResponse(
			req.ReqId, fmt.Sprintf("request ID %d is already being handled", req.ReqId))
		return
	}
	serviceAddr, err := c.serviceRegistry.GetServiceAddress(gatewayReq.GetAuthority())
	if err != nil {
		glog.Errorf("[SyncRpc] error getting service address: %v", err)
		return
	}
	ctx, cancelFn := context.WithCancel(inCtx)
	c.updateOutstandingdReqs(req.ReqId, cancelFn)

	go func() {
		c.broker.send(ctx, serviceAddr, req, c.respCh)
		c.removeOutstandingReq(req.ReqId)
	}()
}

// processStream handles the incoming gateway requests
func (c *SyncRpcClient) processStream(
	ctx context.Context, stream protos.SyncRPCService_EstablishSyncRPCStreamClient) error {

	for {
		req, err := stream.Recv()
		if err != nil {
			glog.Errorf("[SyncRPC] error: %v, failed handling sync request", err)
			return err
		}

		c.handleSyncRpcRequest(ctx, req)
		select {
		case <-ctx.Done():
			glog.Info("[SyncRPC] exiting processing stream")
			return nil
		default:
			break
		}
	}
}

// broker handles the responsibility of translating the incoming SyncGrpcRequest
// into a http2 request to the Grpc service running on the gateway
type broker interface {
	// Send method sends across the request to the appropriate gateway service
	send(context.Context, string, *protos.SyncRPCRequest, chan *protos.SyncRPCResponse)
}

type brokerImpl struct {
	cfg *Config
}

func newbrokerImpl(cfg *Config) *brokerImpl {
	return &brokerImpl{cfg: cfg}
}

func getRequestHeaders(hdr http.Header, trailer http.Header) map[string]string {
	// following block reads the response
	respHeaders := make(map[string]string, len(hdr)+len(trailer))
	for k, v := range hdr {
		respHeaders[k] = strings.Join(v, ",")
	}

	if len(trailer) > 0 {
		for k, v := range trailer {
			respHeaders[k] = strings.Join(v, ",")
		}
	}
	return respHeaders
}

func buildSyncRpcErrorResponse(reqID uint32, err string) *protos.SyncRPCResponse {
	return &protos.SyncRPCResponse{
		ReqId: reqID,
		RespBody: &protos.GatewayResponse{
			Err: err,
		},
	}
}

func buildHeaders(hdrMap map[string]string) http.Header {
	hdr := http.Header{}
	for k, v := range hdrMap {
		hdr.Add(k, v)
	}
	return hdr
}

func sendInternal(address string, req *protos.SyncRPCRequest, respCh chan *protos.SyncRPCResponse, tout time.Duration) {
	var respErr error
	defer func() {
		// handle and send error
		if respErr != nil {
			respErr = fmt.Errorf("ReqID %d failed: %v", req.ReqId, respErr)
			glog.Errorf("[SyncRPC] %v", respErr)
			respCh <- buildSyncRpcErrorResponse(req.ReqId, respErr.Error())
		}
	}()
	// populate headers
	gatewayReq := req.ReqBody

	// http2 client to connect to the grpc port
	// override DialTLS to create a vannilla tcp connection
	client := &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(netw, addr)
			},
		},
		Timeout: tout,
	}
	brokerReq := &http.Request{
		RequestURI: "",
		Method:     "POST",
		URL: &url.URL{
			Scheme: "http",
			Path:   gatewayReq.Path,
			Host:   address,
		},
		Body:   ioutil.NopCloser(bytes.NewReader(gatewayReq.Payload)),
		Header: buildHeaders(gatewayReq.Headers),
		Host:   gatewayReq.Authority,
	}
	resp, err := client.Do(brokerReq)
	if err != nil {
		respErr = fmt.Errorf("grpc request failed: %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		respErr = fmt.Errorf(
			"http response error, status %s statuscode %d, err: %v", resp.Status, resp.StatusCode, err)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		respErr = fmt.Errorf("failed reading grpc message: %v", err)
	}
	lb := len(b)
	if lb < GRPC_MSGLEN_SZ {
		glog.V(2).Infof("[SyncRPC] resp body for ReqId %d is too short: %d", req.ReqId, lb)
		respCh <- &protos.SyncRPCResponse{
			ReqId: req.ReqId,
			RespBody: &protos.GatewayResponse{
				Status:  strconv.Itoa(resp.StatusCode),
				Headers: getRequestHeaders(resp.Header, resp.Trailer),
			},
		}
		return
	}
	// grpc messages are sent as length prefixed messages, first byte is
	// for indicating compression next 4 bytes is length, followed by payload
	// read the first 5 bytes and using the length read the rest of the message
	msgLen := binary.BigEndian.Uint32(b[GRPC_LEN_OFFSET:])
	respMsg := &protos.SyncRPCResponse{
		ReqId: req.ReqId,
		RespBody: &protos.GatewayResponse{
			Status:  strconv.Itoa(resp.StatusCode),
			Headers: getRequestHeaders(resp.Header, resp.Trailer),
			Payload: b,
		},
	}
	respCh <- respMsg
	glog.V(3).Infof("[SyncRPC] sending resp: %s (payload: %v; msg len: %d; body len: %d)", respMsg, b, msgLen, lb)
}

func (p *brokerImpl) send(
	ctx context.Context, serviceAddr string, req *protos.SyncRPCRequest, respCh chan *protos.SyncRPCResponse) {

	clientRespCh := make(chan *protos.SyncRPCResponse)
	go sendInternal(serviceAddr, req, clientRespCh, p.cfg.GatewayResponseTimeout)

	timer := time.NewTimer(p.cfg.GatewayKeepaliveInterval)
	defer timer.Stop()

	timeoutTime := time.Now().Add(p.cfg.GatewayResponseTimeout)
	for {
		select {
		case resp, ok := <-clientRespCh:
			if !ok {
				glog.Errorf("[SyncRPC] channel closed for ReqId %d", req.ReqId)
			} else {
				respCh <- resp
			}
			return
		case <-timer.C:
			if time.Now().After(timeoutTime) {
				// max request timeout exceeded send error back to the caller
				respCh <- buildSyncRpcErrorResponse(req.ReqId, "grpc request timed out on read")
				return
			}
			// construct SyncRpcResponse and keep connection active
			respCh <- &protos.SyncRPCResponse{
				ReqId:    req.ReqId,
				RespBody: &protos.GatewayResponse{KeepConnActive: true}}
			timer.Reset(p.cfg.GatewayKeepaliveInterval)
			glog.V(2).Infof("[SyncRPC] sending keepalive while on ReqId %d", req.ReqId)
		case <-ctx.Done():
			return
		}
	}
}
