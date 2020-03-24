/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package service implements the core of bootstrapper
package service

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"magma/gateway/service_registry"
	"magma/orc8r/lib/go/definitions"
	"magma/orc8r/lib/go/protos"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/http2"
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

// SyncRpcClient opens a bidirectional connection with the cloud
type SyncRpcClient struct {
	// service registry to coonect to the cloud
	serviceRegistry service_registry.GatewayRegistry

	// responseTimeout in seconds
	cfg Config

	// requests which have been terminated
	terminatedReqs map[uint32]bool

	// terminatedReqsMux
	terminatedReqsMux sync.RWMutex

	// requests which are still being processed
	outstandingReqs map[uint32]context.CancelFunc

	outstandingReqsMux sync.RWMutex

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
		terminatedReqs:  make(map[uint32]bool),
		respCh:          make(chan *protos.SyncRPCResponse),
		broker:          newbrokerImpl(cfg),
	}
	return client
}

// Run starts SyncRPC worker loop and blocks forever
func (c *SyncRpcClient) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for currentBackoffInterval := MinRetryInterval; ; currentBackoffInterval = (currentBackoffInterval + RetryBackoffIncrement) % MaxRetryInterval {

		conn, err := c.serviceRegistry.GetCloudConnection(definitions.DispatcherServiceName)
		if err != nil {
			// TODO
			// Continue for retryable grpc errors
			// Add a delay/jitter for retrying cloud connection for non retryable grpc errors
			log.Printf("[SyncRpc] error creating cloud connection: %v", err)
			time.Sleep(currentBackoffInterval)
			continue
		} else {
			currentBackoffInterval = MinRetryInterval // reset backoff interval
			log.Printf("[SyncRpc] successfully connected to cloud '%s' service", definitions.DispatcherServiceName)
		}

		// this should simply wait here for requests and process responses
		// in case we see any error we will retry connecting to the dispatcher
		c.runSyncRpcClient(ctx, protos.NewSyncRPCServiceClient(conn))
		conn.Close()
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
		log.Printf("[SyncRPC] Failed establishing SyncRpc stream; error: %v", err)
		return err
	}

	// close connection if error
	errChan := make(chan error)
	go func() {
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
					log.Printf("[SyncRpc] send to dispatcher failed: %v", err)
					return err
				}
				timer.Reset(c.cfg.SyncRpcHeartbeatInterval)
			}

		case err := <-errChan:
			// error recd handling processStream
			return err

		case <-timer.C:
			err := stream.Send(&protos.SyncRPCResponse{HeartBeat: true})
			if err != nil {
				log.Printf("[SyncRpc] heartbeat to dispatcher failed")
				return err
			}
			timer.Reset(c.cfg.SyncRpcHeartbeatInterval)

		case <-ctx.Done():
			log.Printf("[SyncRPC] Stopping SyncRpcClient")
			return nil
		}
	}
}

func (c *SyncRpcClient) updateTerminatedReqs(reqID uint32) {
	c.terminatedReqsMux.Lock()
	defer c.terminatedReqsMux.Unlock()
	c.terminatedReqs[reqID] = true
}

func (c *SyncRpcClient) isReqTerminated(reqID uint32) bool {
	c.terminatedReqsMux.RLock()
	defer c.terminatedReqsMux.RUnlock()
	_, ok := c.terminatedReqs[reqID]
	return ok
}

func (c *SyncRpcClient) removeTerminatedReqs(reqID uint32) {
	c.terminatedReqsMux.Lock()
	defer c.terminatedReqsMux.Unlock()
	delete(c.terminatedReqs, reqID)
}

func (c *SyncRpcClient) updateOutstandingdReqs(reqID uint32, cancelFn context.CancelFunc) {
	c.outstandingReqsMux.Lock()
	defer c.outstandingReqsMux.Unlock()
	c.outstandingReqs[reqID] = cancelFn
}

func (c *SyncRpcClient) getOutstandingReqCancelFn(reqID uint32) context.CancelFunc {
	c.outstandingReqsMux.RLock()
	defer c.outstandingReqsMux.RUnlock()
	cancelFn, ok := c.outstandingReqs[reqID]
	if !ok {
		return nil
	}
	return cancelFn
}

func (c *SyncRpcClient) removeOutstandingReqs(reqID uint32) {
	c.outstandingReqsMux.Lock()
	defer c.outstandingReqsMux.Unlock()
	delete(c.outstandingReqs, reqID)
}

// handleSyncRpcRequest forwards the incoming request to the appropriate destination
func (c *SyncRpcClient) handleSyncRpcRequest(inCtx context.Context, req *protos.SyncRPCRequest) {
	if req == nil {
		log.Printf("[SyncRpc] error empty request received")
		return
	}

	if req.HeartBeat {
		return
	}

	// get the cancellation fn of the request if present
	cancelFn := c.getOutstandingReqCancelFn(req.ReqId)
	if req.ConnClosed {
		c.updateTerminatedReqs(req.ReqId)
		if cancelFn != nil {
			cancelFn()
		}
		return
	}

	// return early if request is outstanding
	if cancelFn != nil {
		c.respCh <- buildSyncRpcErrorResponse(req.ReqId, fmt.Sprintf("request ID %d is already being handled", req.ReqId))
		return
	}

	ctx, cancelFn := context.WithCancel(inCtx)
	c.removeTerminatedReqs(req.ReqId)
	c.updateOutstandingdReqs(req.ReqId, cancelFn)
	gatewayReq := req.GetReqBody()
	serviceAddr, err := c.serviceRegistry.GetServiceAddress(gatewayReq.Authority)
	if err != nil {
		log.Printf("[SyncRpc] error getting service address: %v", err)
		return
	}

	go func() {
		c.broker.send(ctx, serviceAddr, req, c.respCh)
		c.removeOutstandingReqs(req.ReqId)
	}()
}

// processStream handles the incoming gateway requests
func (c *SyncRpcClient) processStream(ctx context.Context,
	stream protos.SyncRPCService_EstablishSyncRPCStreamClient) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			log.Printf("[SyncRPC] error: %v, failed handling sync request", err)
			return err
		}

		c.handleSyncRpcRequest(ctx, req)
		select {
		case <-ctx.Done():
			log.Printf("[SyncRPC] exiting processing stream")
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
	respHeaders := make(map[string]string, len(hdr))
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

func (p *brokerImpl) sendInternal(ctx context.Context, serviceAddr string, req *protos.SyncRPCRequest, respCh chan *protos.SyncRPCResponse) {
	defer close(respCh)
	var respErr error
	defer func() {
		// handle and send error
		if respErr != nil {
			respErr = fmt.Errorf("ReqID %d failed: %v", req.ReqId, respErr)
			log.Printf("[SyncRPC] error: %v", respErr)
			respCh <- buildSyncRpcErrorResponse(req.ReqId, respErr.Error())
		}
	}()

	// populate headers
	gatewayReq := req.ReqBody

	// http2 client to connect to the grpc port
	// override DialTLS to create a vannilla tcp connection
	client := &http.Client{Transport: &http2.Transport{
		AllowHTTP: true,
		DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(netw, addr)
		}}}

	brokerReq := &http.Request{
		RequestURI: "",
		Method:     "POST",
		URL: &url.URL{
			Scheme: "http",
			Path:   gatewayReq.Path,
			Host:   serviceAddr,
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

	defer resp.Body.Close()
	reader := io.Reader(resp.Body)

	// grpc messages are sent as length prefixed messages, first byte is
	// for indicating compression next 4 bytes is length, followed by payload
	// read the first 5 bytes and using the length read the rest of the message
	buf := make([]byte, GRPC_MSGLEN_SZ)
	for {
		_, err := io.ReadFull(reader, buf)
		if err != nil {
			if err == io.EOF {
				respCh <- &protos.SyncRPCResponse{
					ReqId: req.ReqId,
					RespBody: &protos.GatewayResponse{
						Status:  strconv.Itoa(resp.StatusCode),
						Headers: getRequestHeaders(resp.Header, resp.Trailer),
					},
				}
			} else {
				respErr = fmt.Errorf("failed reading length prefix of grpc message: %v", err)
			}
			return
		}
		msgLen := binary.BigEndian.Uint32(buf[GRPC_LEN_OFFSET:])
		msgBuf := make([]byte, msgLen)
		_, err = io.ReadFull(reader, msgBuf)
		if err != nil {
			if err == io.EOF {
				respCh <- &protos.SyncRPCResponse{
					ReqId: req.ReqId,
					RespBody: &protos.GatewayResponse{
						Status:  strconv.Itoa(resp.StatusCode),
						Headers: getRequestHeaders(resp.Header, resp.Trailer),
					},
				}
			} else {
				respErr = fmt.Errorf("failed reading data of %d size: %v", msgLen, err)
			}
			return
		}
		respCh <- &protos.SyncRPCResponse{
			ReqId: req.ReqId,
			RespBody: &protos.GatewayResponse{
				Status:  strconv.Itoa(resp.StatusCode),
				Headers: getRequestHeaders(resp.Header, resp.Trailer),
				Payload: append(buf, msgBuf...),
			},
		}
		// check if context is being terminated
		select {
		case <-ctx.Done():
			return
		default:
			// do nothing
		}
	}
}

func (p *brokerImpl) send(
	ctx context.Context, serviceAddr string, req *protos.SyncRPCRequest, respCh chan *protos.SyncRPCResponse) {

	clientRespCh := make(chan *protos.SyncRPCResponse)
	go p.sendInternal(ctx, serviceAddr, req, clientRespCh)

	timer := time.NewTimer(p.cfg.GatewayKeepaliveInterval)
	defer timer.Stop()

	lastMsgReadTime := time.Now()
	for {
		select {
		case resp, ok := <-clientRespCh:
			if !ok {
				return
			}
			respCh <- resp
			timer.Reset(p.cfg.GatewayKeepaliveInterval)
		case <-timer.C:
			if time.Now().Sub(lastMsgReadTime) > p.cfg.GatewayResponseTimeout {
				// max request timeout exceeded send error back to the caller
				respCh <- buildSyncRpcErrorResponse(req.ReqId, "grpc request timed out on read")
				return
			} else {
				// construct SyncRpcResponse and keep connection active
				respCh <- &protos.SyncRPCResponse{
					ReqId:    req.ReqId,
					RespBody: &protos.GatewayResponse{KeepConnActive: true}}
				timer.Reset(p.cfg.GatewayKeepaliveInterval)
			}
		case <-ctx.Done():
			return
		}
	}
}
