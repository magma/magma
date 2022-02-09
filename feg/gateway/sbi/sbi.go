/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sbi

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	// ShutdownTimeout is used as a timeout to exit the Echo Server Shutdown
	ShutdownTimeout = 10 * time.Second
)

// EchoServer is a wrapper of echo which includes some extra functions like StartWithWait
// can be used as a sbi server
type EchoServer struct {
	*echo.Echo
}

// NotifierServer is the server handling notifications on a sbi client
type NotifierServer struct {
	Server      *EchoServer
	NotifierCfg NotifierConfig
}

// RemoteConfig represents the server configuration we are going to connect to
type RemoteConfig struct {
	ApiRoot      url.URL
	TokenUrl     string
	ClientId     string
	ClientSecret string
}

// NotifierConfig represent the configuration for the local NotifierServer
type NotifierConfig struct {
	LocalAddr     string // ip:port or :port (can differ once started, Use GetListenerAddr)
	NotifyApiRoot string
}

// BaseClientWithNotifier is a struct representing a generic sbi client. We can use this generic
// client with any sbi specific clients (like N7 or N40)
type BaseClientWithNotifier struct {
	RemoteCfg    RemoteConfig
	NotifyServer *NotifierServer // NotifyServer
}

func NewEchoServer() *EchoServer {
	e := echo.New()
	if glog.V(2) {
		e.Use(ServerLoggingMiddleware())
	}
	return &EchoServer{e}
}

func NewBaseClientWithNotifyServer(notifierConfig NotifierConfig, remoteConfig RemoteConfig) *BaseClientWithNotifier {
	return &BaseClientWithNotifier{
		RemoteCfg:    remoteConfig,
		NotifyServer: NewNotifierServer(notifierConfig),
	}
}

func NewNotifierServer(notifierConfig NotifierConfig) *NotifierServer {
	return &NotifierServer{
		Server:      NewEchoServer(),
		NotifierCfg: notifierConfig,
	}
}

func (s *EchoServer) StartWithWait(addrs string) error {
	addr, err := net.ResolveTCPAddr("tcp", addrs)
	if err != nil {
		return err
	}
	errChan := make(chan error)
	go func() {
		err := s.Start(addr.String())
		if err != nil {
			errChan <- err
		}
	}()

	return s.waitForServer(errChan)
}

func (s *EchoServer) Shutdown(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	s.Server.Shutdown(ctx)
}

// GetListenerAddr returns the actual address running. Use this function instead of checking
// NotifierConfig.LocalAddr
func (s *EchoServer) GetListenerAddr() (net.Addr, error) {
	addr := s.ListenerAddr()
	if addr != nil && strings.Contains(addr.String(), ":") {
		return addr, nil
	}
	return nil, fmt.Errorf("Can't get listener from NotifierServer")
}

// waitForServer waits for the Echo server to be launched.
func (s *EchoServer) waitForServer(errChan <-chan error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			_, err := s.GetListenerAddr()
			if err != nil {
				// server is not started yet. Try again
				continue
			}
			return nil
		case err := <-errChan:
			if err == http.ErrServerClosed {
				// not actually an error. Happens when the server is gracefully shutdown
				return nil
			}
			return err
		}
	}
}

func (rc RemoteConfig) BuildServerString() string {
	return fmt.Sprintf("%s://%s", rc.ApiRoot.Scheme, rc.ApiRoot.Host)
}

// BuildHttpClient returns an HTTP client to use with any sbi client
func (rc RemoteConfig) BuildHttpClient() *http.Client {
	tokenConfig := clientcredentials.Config{
		ClientID:     rc.ClientId,
		ClientSecret: rc.ClientSecret,
		TokenURL:     rc.TokenUrl,
	}
	var httpClient *http.Client
	if glog.V(2) {
		httpClient = NewLoggingHttpClient()
	} else {
		httpClient = &http.Client{}
	}
	tokenCtxt := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	return tokenConfig.Client(tokenCtxt)
}

func (s *NotifierServer) Start() error {
	err := s.Server.StartWithWait(s.NotifierCfg.LocalAddr)
	if err != nil {
		return fmt.Errorf("NotifierServer could not Start: %s", err)
	}
	addr, err := s.Server.GetListenerAddr()
	if err != nil {
		return fmt.Errorf("NotifierServer could not get address after Start: %s", err)
	}
	// update/rewrite the address (just in case we used port 0)
	s.NotifierCfg.LocalAddr = addr.String()
	return nil
}

func (s *NotifierServer) Stop() {
	s.Server.Shutdown(ShutdownTimeout)
}
