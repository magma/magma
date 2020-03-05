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
	"magma/gateway/cloud_registry"
	"magma/orc8r/lib/go/definitions"
	"magma/orc8r/lib/go/protos"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/http2"
)

const (
	// grpc message is delivered as a length prefixed message within the HTTP2 DATA
	// frames, with first 5 bytes for compression and msg length
	// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md
	GRPC_MSGLEN_SZ  = 5
	GRPC_LEN_OFFSET = 1
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
	// address of the dispatcher
	cloudRegistry cloud_registry.ProxiedCloudRegistry

	// responseTimeout in seconds
	cfg *Config

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

func NewSyncRpcClient(cfg *Config) *SyncRpcClient {
	// create a new grpc connection to the dispatcher in the cloud??
	client := &SyncRpcClient{
		cloudRegistry:  cloud_registry.ProxiedCloudRegistry{},
		cfg:            cfg,
		terminatedReqs: make(map[uint32]bool),
		respCh:         make(chan *protos.SyncRPCResponse),
		broker:         newbrokerImpl(cfg),
	}
	return client
}

func (c *SyncRpcClient) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		conn, err := c.cloudRegistry.GetCloudConnection(definitions.DispatcherServiceName)
		if err != nil {
			// TODO
			// Continue for retryable grpc errors
			// Add a delay/jitter for retrying cloud connection for non retryable grpc errors
			log.Printf("[SyncRpc] error %v creating cloud connection", err)
			continue
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
		log.Printf("[SyncRPC] error: %v, Failed establishing SyncRpc stream", err)
	}

	// close connection if error
	errChan := make(chan error)
	go func() {
		err := c.processStream(ctx, stream)
		errChan <- err
	}()

	timer := time.NewTimer(c.cfg.SyncRpcHeartbeatInterval)
	defer timer.Stop()
	for {
		select {
		case resp := <-c.respCh:
			if !c.isReqTerminated(resp.ReqId) {
				err := stream.Send(resp)
				if err != nil {
					log.Printf("[SyncRpc] send to dispatcher failed %v", err)
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
		log.Printf("[SyncRpc] error empty request recd")
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
	serviceAddr, err := cloud_registry.GetServiceAddress(gatewayReq.Authority)
	if err != nil {
		log.Printf("[SyncRpc] error: %v getting service address", err)
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
			respErr = errors.Wrap(respErr, fmt.Sprintf("ReqID %d failed", req.ReqId))
			log.Printf("[SyncRPC] error %v", respErr)
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
		respErr = errors.Wrap(err, " grpc request failed ")
		return
	}
	if resp.StatusCode != http.StatusOK {
		respErr = errors.Wrap(err, fmt.Sprintf(" http response error status %s  statuscode %d ", resp.Status, resp.StatusCode))
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
				respErr = errors.Wrap(err, "failed reading length prefix of grpc message")
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
				respErr = errors.Wrap(err, fmt.Sprintf("failed reading data of %d size", msgLen))
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

func (p *brokerImpl) send(ctx context.Context, serviceAddr string, req *protos.SyncRPCRequest, respCh chan *protos.SyncRPCResponse) {
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
				respCh <- &protos.SyncRPCResponse{ReqId: req.ReqId, RespBody: &protos.GatewayResponse{KeepConnActive: true}}
				timer.Reset(p.cfg.GatewayKeepaliveInterval)
			}
		case <-ctx.Done():
			return
		}
	}
}
