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

package server

import (
	"context"
	"errors"
	"fbc/cwf/radius/config"
	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/modules/protos"
	"fbc/cwf/radius/monitoring"
	"fbc/cwf/radius/session"
	"fmt"
	"math/rand"
	"net"

	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"
	"fbc/lib/go/radius/rfc2866"

	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// GRPCListener listens to gRpc
type GRPCListener struct {
	Listener
	GrpcServer *grpc.Server
	Server     *Server
	Port       int
	ready      chan bool
}

// GRPCListenerExtraConfig extra config for GRPC listener
type GRPCListenerExtraConfig struct {
	Port int `json:"port"`
}

// NewGRPCListener ...
func NewGRPCListener() *GRPCListener {
	return &GRPCListener{
		ready: make(chan bool),
	}
}

// Init override
func (l *GRPCListener) Init(
	server *Server,
	serverConfig config.ServerConfig,
	listenerConfig config.ListenerConfig,
	_ monitoring.ListenerCounters,
) error {
	if server == nil {
		return errors.New("cannot initialize GRPC listener with null server")
	}

	// Parse configuration
	var cfg GRPCListenerExtraConfig
	err := mapstructure.Decode(listenerConfig.Extra, &cfg)
	if err != nil {
		return err
	}

	l.Server = server
	l.Port = cfg.Port
	return nil
}

// ListenAndServe override
func (l *GRPCListener) ListenAndServe() error {
	// Start listenning
	listenAddress := fmt.Sprintf(":%d", l.Port)
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		l.ready <- false
		return errors.New("grpc listener: failed to open tcp connection" + listenAddress)
	}

	// Start serving
	l.GrpcServer = grpc.NewServer()
	protos.RegisterAuthorizationServer(l.GrpcServer, &authorizationServer{Listener: l})
	go func() {
		l.GrpcServer.Serve(lis)
	}()

	// Signal listener is ready
	go func() {
		l.ready <- true
	}()
	return nil
}

// GetHandleRequest override
func (l *GRPCListener) GetHandleRequest() modules.Middleware {
	return l.HandleRequest
}

// Shutdown override
func (l *GRPCListener) Shutdown(ctx context.Context) error {
	return nil
}

// Ready override
func (l *GRPCListener) Ready() chan bool {
	return l.ready
}

// SetConfig override
func (l *GRPCListener) SetConfig(c config.ListenerConfig) {
	l.Config = c
}

type authorizationServer struct {
	Listener *GRPCListener
}

func (s *authorizationServer) Change(ctx context.Context, request *protos.ChangeRequest) (*protos.CoaResponse, error) {
	// Convert to RADIUS request
	req := radius.Request{
		Packet: &radius.Packet{
			Code:   radius.CodeDisconnectRequest,
			Secret: []byte(s.Listener.Server.config.Secret),
		},
	}

	// Handle RADIUS request
	return s.handleCoaRequest(request.Ctx, &req)
}

func (s *authorizationServer) Disconnect(ctx context.Context, request *protos.DisconnectRequest) (*protos.CoaResponse, error) {
	// Convert to RADIUS request
	req := radius.Request{
		Packet: &radius.Packet{
			Code:   radius.CodeDisconnectRequest,
			Secret: []byte(s.Listener.Server.config.Secret),
		},
	}

	// Handle RADIUS request
	return s.handleCoaRequest(request.Ctx, &req)
}

func (s *authorizationServer) handleCoaRequest(ctx *protos.Context, request *radius.Request) (*protos.CoaResponse, error) {
	if ctx == nil {
		return nil, errors.New("cannot handle a request without context")
	}

	if request == nil {
		return nil, errors.New("cannot handle a nil request")
	}

	// Get session ID from the request, if exists, and setup correlation ID
	srv := s.Listener.Server
	var correlationField = zap.Uint32("correlation", rand.Uint32())
	requestContext := modules.RequestContext{
		RequestID: correlationField.Integer,
		Logger:    srv.logger.With(correlationField),
		SessionID: ctx.SessionId,
		SessionStorage: session.NewSessionStorage(
			srv.multiSessionStorage,
			ctx.SessionId,
		),
	}

	// Load state, read CoA identifier and persist the state again
	state, err := requestContext.SessionStorage.Get()
	if err != nil {
		return nil, err
	}

	// Add Acct-Session-Id attribute
	request.Attributes = radius.Attributes{}
	request.Set(rfc2866.AcctSessionID_Type, radius.Attribute(state.AcctSessionID))
	request.Set(rfc2865.CallingStationID_Type, radius.Attribute(ctx.MacAddr))

	// Set Identifier
	request.Identifier = state.NextCoAIdentifier
	state.NextCoAIdentifier = (state.NextCoAIdentifier + 1) % 0xFF

	// Handle
	counter := monitoring.NewOperation("handle_grpc").Start()
	res, err := s.Listener.HandleRequest(&requestContext, request)
	if err != nil {
		requestContext.Logger.Error("failed to handle request", zap.Error(err))
		counter.Failure("grpc_handle_error")
		return nil, err
	}
	if res == nil {
		requestContext.Logger.Error("got nil response")
		counter.Failure("grpc_nil_response")
		return nil, err
	}
	counter.Success()

	// Persist state
	err = srv.multiSessionStorage.Set(ctx.SessionId, *state)
	if err != nil {
		return nil, err
	}

	// Convert response to CoA response
	return &protos.CoaResponse{
		CoaResponseType: convertCoaCode(res.Code),
		Ctx:             ctx,
	}, nil
}

func convertCoaCode(code radius.Code) protos.CoaResponseCoaResponseTypeEnum {
	if code == radius.CodeCoAACK || code == radius.CodeDisconnectACK {
		return protos.CoaResponse_ACK
	}
	return protos.CoaResponse_NAK
}
