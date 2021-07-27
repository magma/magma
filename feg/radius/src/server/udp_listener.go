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
	"fbc/cwf/radius/config"
	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/monitoring"
	"fbc/lib/go/radius"
	"fmt"
	"math/rand"

	"fbc/cwf/radius/session"
	"sync/atomic"

	"github.com/mitchellh/mapstructure"
	"github.com/patrickmn/go-cache"
	"go.opencensus.io/tag"
	"go.uber.org/zap"
)

// UDPListener listens to Radius udp packets
type UDPListener struct {
	Listener
	Server *radius.PacketServer
	ready  chan bool
}

// UDPListenerExtraConfig extra config for UDP listener
type UDPListenerExtraConfig struct {
	Port int `json:"port"`
}

// NewUDPListener ...
func NewUDPListener() *UDPListener {
	return &UDPListener{
		ready: make(chan bool),
	}
}

// Init override
func (l *UDPListener) Init(
	server *Server,
	serverConfig config.ServerConfig,
	listenerConfig config.ListenerConfig,
	ctrs monitoring.ListenerCounters,
) error {
	// Parse configuration
	var cfg UDPListenerExtraConfig
	err := mapstructure.Decode(listenerConfig.Extra, &cfg)
	if err != nil {
		return err
	}

	// Create packet server
	l.Server = &radius.PacketServer{
		Handler: radius.HandlerFunc(
			generatePacketHandler(l, server, ctrs),
		),
		SecretSource: radius.StaticSecretSource([]byte(serverConfig.Secret)),
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Ready:        make(chan bool),
	}
	return nil
}

// ListenAndServe override
func (l *UDPListener) ListenAndServe() error {
	serverError := make(chan error, 1)
	go func() {
		err := l.Server.ListenAndServe()
		serverError <- err
	}()

	// Wait to see if initialization was successful
	select {
	case _ = <-l.Server.Ready:
		l.ready <- true
		return nil
	case err := <-serverError:
		l.ready <- false
		return err // might be nil if no error
	}
}

// GetHandleRequest override
func (l *UDPListener) GetHandleRequest() modules.Middleware {
	return l.HandleRequest
}

// Shutdown override
func (l *UDPListener) Shutdown(ctx context.Context) error {
	return l.Server.Shutdown(ctx)
}

// Ready override
func (l *UDPListener) Ready() chan bool {
	return l.ready
}

// SetConfig override
func (l *UDPListener) SetConfig(c config.ListenerConfig) {
	l.Config = c
}

// generatePacketHandler A generic handler method to incoming RADIUS packets
func generatePacketHandler(
	l ListenerInterface,
	server *Server,
	ctrs monitoring.ListenerCounters,
) func(radius.ResponseWriter, *radius.Request) {
	server.logger.Debug(
		"Registering handler for listener",
		zap.String("listener", l.GetConfig().Name),
	)
	return func(w radius.ResponseWriter, r *radius.Request) {
		// Make sure no duplicate packet
		dedupOperation := server.counters.DedupPacket.Start(
			tag.Upsert(monitoring.ListenerTag, l.GetConfig().Name),
		)
		requestKey := fmt.Sprintf("%s_%d", r.RemoteAddr, r.Identifier)

		if _, found := server.dedupSet.Get(requestKey); found {
			server.logger.Warn(
				"Duplicate packet was receieved and dropped",
				zap.Stringer("source_ip", r.RemoteAddr),
				zap.Int("identifier", int(r.Identifier)),
			)
			atomic.AddUint32(l.GetDupDropped(), 1)
			dedupOperation.Failure("duplicate_packet_dropped")
			return
		}
		server.dedupSet.Set(requestKey, "-", cache.DefaultExpiration)
		dedupOperation.Success()

		// Get session ID from the request, if exists, and setup correlation ID
		var correlationField = zap.Uint32("correlation", rand.Uint32())
		sessionID := server.GetSessionID(r)
		generatedSessionID := server.GenSessionID(r)

		// Create request context
		requestContext := modules.RequestContext{
			RequestID:      correlationField.Integer,
			Logger:         server.logger.With(correlationField),
			SessionID:      sessionID,
			SessionStorage: session.NewSessionStorageExt(server.multiSessionStorage, sessionID, generatedSessionID),
		}

		// Execute filters
		filterProcessCounter := monitoring.NewOperation("filter_process").Start()
		for _, filter := range server.filters {
			err := filter.Code.Process(&requestContext, l.GetConfig().Name, r)
			if err != nil {
				server.logger.Error("Failed to process reqeust by filter", zap.Error(err), correlationField)
				filterProcessCounter.Failure(
					"filter_failed",
					tag.Upsert(monitoring.FilterTag, filter.Name),
				)
				return
			}
		}
		filterProcessCounter.Success()

		// Execute modules
		listenerHandleCounter := ctrs.StartRequest(r.Code)
		response, err := l.GetHandleRequest()(&requestContext, r)
		if err != nil {
			server.logger.Error("Failed to handle reqeust by listener", zap.Error(err), correlationField)
			listenerHandleCounter.Failure("handle_failed")
			return
		}
		if response == nil {
			server.logger.Error("Got nil response from handler. Response will not be sent", correlationField)
			listenerHandleCounter.Failure("nil_response")
			return
		}
		listenerHandleCounter.GotResponse(response.Code)

		if response == nil {
			server.logger.Warn(
				"Request failed to be handled, as no response returned",
				correlationField,
			)
			return
		}

		// Build response
		server.logger.Warn(
			"Request successfully handled",
			correlationField,
		)
		radiusResponse := r.Response(response.Code)
		for key, values := range response.Attributes {
			for _, value := range values {
				radiusResponse.Add(key, value)
			}
		}
		w.Write(radiusResponse)
	}
}
