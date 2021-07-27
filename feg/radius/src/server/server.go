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
	"fbc/cwf/radius/filters"
	"fbc/cwf/radius/loader"
	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/monitoring"
	"fbc/cwf/radius/session"
	"fmt"
	"sync/atomic"
	"time"

	"go.opencensus.io/tag"

	"fbc/lib/go/radius"

	"github.com/patrickmn/go-cache"

	"go.uber.org/zap"
)

type (
	// RequestContext Info about the request and utils for the handler
	RequestContext struct {
		context.Context
		RequestID      int64
		Logger         *zap.Logger
		SessionStorage session.Storage
	}

	// Response the response of a plugin handler
	Response struct {
		Code       radius.Code
		Attributes radius.Attributes
	}

	// Filter represents a server pluggable filter
	Filter struct {
		Name string
		Code filters.Filter
	}

	// Module represents a listener module
	Module struct {
		Name    string
		Context modules.Context
		Code    modules.Module
	}

	// Server encapsultes an instance of RADIUS server
	Server struct {
		ready               chan bool // wait on this to wait for the server to be ready for work
		terminate           chan bool
		listeners           map[string]ListenerInterface
		filters             []*Filter
		config              config.ServerConfig
		logger              *zap.Logger
		multiSessionStorage session.GlobalStorage
		dedupSet            *cache.Cache
		counters            *monitoring.ServerCounters
	}
)

// New a RADIUS server instance as per config
func New(config config.ServerConfig, logger *zap.Logger, loader loader.Loader) (*Server, error) {
	// Init session storage
	multiSessionStorage, err := initSessionStorage(config, logger)
	if err != nil {
		return nil, err
	}

	// Init server object
	server := Server{
		listeners:           make(map[string]ListenerInterface), // Will be populated by "Start" method
		ready:               make(chan bool, 1),
		filters:             make([]*Filter, 0),
		terminate:           make(chan bool, 1), // Internal channel used for termination of listeners
		config:              config,             // The original config for later reference
		logger:              logger,
		multiSessionStorage: multiSessionStorage,
		dedupSet:            cache.New(config.DedupWindow.Duration, time.Minute),
		counters:            monitoring.CreateServerCounters(),
	}

	serverInitCounter := server.counters.Init.Start()
	logger.Info(
		"allocate new server",
		zap.Int("num_listeners", len(config.Listeners)),
		zap.Int("num_filters", len(config.Filters)),
	)

	// Load filters from config
	for _, filterName := range config.Filters {
		server.counters.FilterInit.Start(
			tag.Upsert(monitoring.FilterTag, filterName),
		)
		filter, err := loader.LoadFilter(filterName)
		if err != nil {
			logger.Error("filter failed to load", zap.String("filter_name", filterName), zap.Error(err))
			server.counters.FilterInit.Failure("load_error")
			return nil, err
		}

		err = filter.Init(&config)
		if err != nil {
			logger.Error("filter failed to init", zap.String("filter_name", filterName), zap.Error(err))
			server.counters.FilterInit.Failure("init_error")
			return nil, err
		}
		server.filters = append(server.filters, &Filter{
			Name: filterName,
			Code: filter,
		})
		server.counters.FilterInit.Success()
	}

	// Load listeners from config
	for _, lconfig := range config.Listeners {
		listenerInitCounter := server.counters.ListenerInit.Start(
			tag.Upsert(monitoring.ListenerTag, lconfig.Name),
		)

		// Create the listener
		var listener ListenerInterface
		switch lconfig.Type {
		case "udp":
			listener = NewUDPListener()
		case "grpc":
			listener = NewGRPCListener()
		case "sse":
			listener = NewSSEListener()
		default:
			logger.Error(
				fmt.Sprintf("failed to create listener, listener type '%s'", lconfig.Type),
				zap.String("listener", lconfig.Name),
			)
			break
		}

		// Set configuration for listener
		listener.SetConfig(lconfig)

		// Load modules
		for _, modDesc := range lconfig.Modules {
			moduleInitCounter := server.counters.ModuleInit.Start(
				tag.Upsert(monitoring.ListenerTag, lconfig.Name),
				tag.Upsert(monitoring.ModuleTag, modDesc.Name),
			)

			logger.Info("loading module", zap.String("module_name", modDesc.Name))
			// Load module
			module, err := loader.LoadModule(modDesc.Name)
			if err != nil {
				logger.Error("module failed to load", zap.String("module_name", modDesc.Name), zap.Error(err))
				moduleInitCounter.Failure("load_error")
				return nil, err
			}
			logger.Debug(
				"Module loaded successfully",
				zap.String("module_name", modDesc.Name),
				zap.Int("precedence", len(listener.GetModules())),
			)

			// Init the module
			logger.Debug("Initializing module", zap.String("module_name", modDesc.Name))
			moduleCtx, err := module.Init(logger, modDesc.Config)
			if err != nil {
				logger.Error("module failed to init", zap.String("module_name", modDesc.Name), zap.Error(err))
				moduleInitCounter.Failure("init_error")
				return nil, err
			}

			listener.AppendModule(&Module{
				Code:    module,
				Context: moduleCtx,
				Name:    modDesc.Name,
			})
			moduleInitCounter.Success()
		}

		// Wrap modules in call chain, leveraging the middleware pattern
		// into the listener's HandleRequest method
		handler := func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
			return nil, nil
		}

		for idx := len(listener.GetModules()) - 1; idx >= 0; idx-- {
			module := (listener.GetModules())[idx]
			handler = wrapMiddleware(listener.GetConfig().Name, handler, module)
		}

		// Initialize the listener
		listener.SetHandleRequest(handler)
		err := listener.Init(
			&server,
			config,
			lconfig,
			monitoring.CreateListenerCounters(lconfig.Name),
		)
		if err != nil {
			listenerInitCounter.Failure("unknown")
			return nil, err
		}

		logger.Debug("listener created", zap.String("listener", lconfig.Name))
		server.listeners[lconfig.Name] = listener
		listenerInitCounter.Success()
	}

	// Down we go!
	serverInitCounter.Success()
	return &server, nil
}

func initSessionStorage(config config.ServerConfig, logger *zap.Logger) (session.GlobalStorage, error) {
	var multiSessionStorage session.GlobalStorage
	if config.SessionStorage == nil || config.SessionStorage.StorageType == "memory" {
		multiSessionStorage = session.NewMultiSessionMemoryStorage()
		logger.Info("created memory config")
	} else if config.SessionStorage.StorageType == "redis" {
		redis := config.SessionStorage.Redis
		multiSessionStorage = session.NewMultiSessionRedisStorage(redis.Addr, redis.Password, redis.DB)
		logger.Info("created redis config", zap.String("addr", redis.Addr))
	} else {
		logger.Error("missing session storage config")
		err := fmt.Errorf("missing session storage config")
		return nil, err
	}

	return multiSessionStorage, nil
}

func wrapMiddleware(listenerName string, next modules.Middleware, module Module) modules.Middleware {
	return func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
		// Start counter
		counter := monitoring.NewOperation(
			module.Name,
			tag.Upsert(monitoring.ListenerTag, listenerName),
		).Start()

		// Handle
		res, err := module.Code.Handle(module.Context, c, r, next)

		// Complete counter operation
		if err != nil {
			counter.Failure("handle_error")
		} else {
			counter.Success()
		}
		return res, err
	}
}

// Start listening and parsing incoming requests
func (s Server) Start() {
	var err error
	s.logger.Debug("starting server", zap.Int("num_listeners", len(s.listeners)), zap.Int("num_filters", len(s.filters)))
	for _, listener := range s.listeners {
		// Create logger
		logger := s.logger.With(zap.String("listener", listener.GetConfig().Name))
		logger.Debug("Starting listener")

		// Initialize logger
		go func(listener ListenerInterface) {
			logger.Debug("listener go-routine starts...")
			err = listener.ListenAndServe()
			if err != nil {
				logger.Error("starting listener failed", zap.Error(err))
			} else {
				logger.Info("listener initialized successfully")
			}
		}(listener)
	}

	for _, listener := range s.listeners {
		logger := s.logger.With(zap.String("listener", listener.GetConfig().Name))
		// wait for listener to initialize
		logger.Info("waiting for listener to be ready...")
		listenerReady := <-listener.Ready()
		if !listenerReady {
			s.logger.Error("some listeners failed to initialize")
			s.ready <- false
			return
		}
		logger.Info("listener is ready")
	}

	// Server is ready!
	s.logger.Info("all listeners ready, server is up and running")
	s.ready <- true

	// Wait for termination
	<-s.terminate
	s.logger.Info("server was terminated")
}

// StartAndWait start the RADIUS server & block until all listeners are ready
func (s Server) StartAndWait() bool {
	go s.Start()
	// wait for server to be complete initialization & read status
	isReady := <-s.ready
	return isReady
}

// GetDroppedCount gets the total count of packets dropped due to duplicate,
// as depicted in rfc2865 section 3 (Identifier)
func (s Server) GetDroppedCount() uint32 {
	var total uint32
	for _, l := range s.listeners {
		total += atomic.LoadUint32(l.GetDupDropped())
	}
	return total
}

// Stop the radius server
func (s Server) Stop() {
	for name, listener := range s.listeners {
		s.logger.Debug("Shutting down listener", zap.String("listener", name))
		if err := listener.Shutdown(context.Background()); err != nil {
			s.logger.Error(
				"Error shutting down listener",
				zap.String("listener", name),
				zap.Error(err),
			)
		}
	}

	// Signal termination
	s.logger.Debug("All listeners are now down, terminating server")
	s.terminate <- true
}

// getSessionStateAPI returns a per-session accessor to session state
func (s Server) getSessionStateAPI(sessionID string) session.Storage {
	return session.NewSessionStorage(s.multiSessionStorage, sessionID)
}

// getSessionState returns the per-session state maintained by the RADIUS server
// this is meant for test code which needs a way to peek into the server's session state to verify test assertions.
func (s Server) getSessionState(sessionID string) (*session.State, error) {
	sessionStg := s.getSessionStateAPI(sessionID)
	return sessionStg.Get()
}
