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

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type NotifierServer struct {
	Server      *echo.Echo
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

type BaseClientWithNotifier struct {
	RemoteCfg    RemoteConfig
	NotifyServer *NotifierServer // NotifyServer
}

func NewBaseClientWithNotifyServer(notifierConfig NotifierConfig, remoteConfig RemoteConfig) *BaseClientWithNotifier {
	return &BaseClientWithNotifier{
		RemoteCfg:    remoteConfig,
		NotifyServer: NewNotifierServer(notifierConfig),
	}
}

func NewNotifierServer(notifierConfig NotifierConfig) *NotifierServer {
	return &NotifierServer{
		Server:      echo.New(),
		NotifierCfg: notifierConfig,
	}
}

func (s *NotifierServer) Start() error {
	errChan := make(chan error)
	addr, err := net.ResolveTCPAddr("tcp", s.NotifierCfg.LocalAddr)
	if err != nil {
		return err
	}

	go func() {
		err = s.Server.Start(addr.String())
		if err != nil {
			errChan <- err
		}
	}()

	return s.waitForServer(errChan)
}

func (s *NotifierServer) Shutdown(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	s.Server.Shutdown(ctx)
}

// GetListenerAddr returns the actual address running. Use this function instead of checking
// NotifierConfig.LocalAddr
func (s *NotifierServer) GetListenerAddr() (net.Addr, error) {
	addr := s.Server.ListenerAddr()
	if addr != nil && strings.Contains(addr.String(), ":") {
		return addr, nil
	}
	return nil, fmt.Errorf("Can't get listener from NotifierServer")
}

// waitForServer waits for the Echo server to be launched.
func (s *NotifierServer) waitForServer(errChan <-chan error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			addr, err := s.GetListenerAddr()
			if err != nil {
				// server is not started yet. Try again
				continue
			}
			// update configured address
			s.NotifierCfg.LocalAddr = addr.String()
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
	tokenCtxt := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{})
	return tokenConfig.Client(tokenCtxt)
}
