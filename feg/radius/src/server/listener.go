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
)

// ListenerInterface encapsulates runtime for concreate listeners
type ListenerInterface interface {
	// Common methods
	GetModules() []Module
	SetModules(m []Module)
	AppendModule(m *Module)
	GetConfig() config.ListenerConfig
	SetConfig(config config.ListenerConfig)
	SetHandleRequest(hr modules.Middleware)
	GetHandleRequest() modules.Middleware
	Ready() chan bool
	GetDupDropped() *uint32

	//
	// Listener-specific methods
	//

	// Initialize the listener
	Init(
		server *Server,
		serverConfig config.ServerConfig,
		listenerConfig config.ListenerConfig,
		counters monitoring.ListenerCounters,
	) error

	// Shutdown Blocking call to shutting down a listener
	Shutdown(ctx context.Context) error

	// ListenAndServe() Starts listenning and serving requests.
	// The method MUST return (as opposed to serving in a loop). If a loop is
	// needed, it should be spawned in a separate go routine from within this
	// method. Notice that when listener is ready, the channel returned from
	// Ready() must be sent a `true` value (or `false` upon failure)
	ListenAndServe() error
}

// Listener base implementation of a listener
type Listener struct {
	ListenerInterface
	Config        config.ListenerConfig
	Modules       []Module
	HandleRequest modules.Middleware
	Server        *Server
	dupDropped    uint32
}

// GetModules ...
func (l *Listener) GetModules() []Module {
	return l.Modules
}

// SetModules ...
func (l *Listener) SetModules(m []Module) {
	l.Modules = m
}

// AppendModule ...
func (l *Listener) AppendModule(m *Module) {
	l.Modules = append(l.Modules, *m)
}

// GetConfig ...
func (l *Listener) GetConfig() config.ListenerConfig {
	return l.Config
}

// SetConfig ...
func (l *Listener) SetConfig(c config.ListenerConfig) {
	l.Config = c
}

// GetHandleRequest ...
func (l *Listener) GetHandleRequest() modules.Middleware {
	return l.HandleRequest
}

// SetHandleRequest ...
func (l *Listener) SetHandleRequest(hr modules.Middleware) {
	l.HandleRequest = hr
}

// GetDupDropped override
func (l *Listener) GetDupDropped() *uint32 {
	return &l.dupDropped
}
