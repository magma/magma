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

type Server struct {
	Server     *echo.Echo
	ListenAddr string
}

type ServerConfig struct {
	ApiRoot      url.URL
	TokenUrl     string
	ClientId     string
	ClientSecret string
}

type ClientConfig struct {
	LocalAddr     string // ip:port or :port
	NotifyApiRoot string
}

func NewSbiServer(listenAddr string) *Server {
	return &Server{
		Server:     echo.New(),
		ListenAddr: listenAddr,
	}
}

func (s *Server) Start() error {
	errChan := make(chan error)
	addr, err := net.ResolveTCPAddr("tcp", s.ListenAddr)
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

func (s *Server) Shutdown(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	s.Server.Shutdown(ctx)
}

// waitForServer waits for the Echo server to be launched.
func (s *Server) waitForServer(errChan <-chan error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			addr := s.Server.ListenerAddr()
			if addr != nil && strings.Contains(addr.String(), ":") {
				// server started
				s.ListenAddr = addr.String()
				return nil
			}
		case err := <-errChan:
			if err == http.ErrServerClosed {
				// not actually an error. Happens when the server is gracefully shutdown
				return nil
			}
			return err
		}
	}
}

func (s ServerConfig) BuildServerString() string {
	return fmt.Sprintf("%s://%s", s.ApiRoot.Scheme, s.ApiRoot.Host)
}

// BuildHttpClient returns an HTTP client to use with any sbi client
func (s ServerConfig) BuildHttpClient() *http.Client {
	tokenConfig := clientcredentials.Config{
		ClientID:     s.ClientId,
		ClientSecret: s.ClientSecret,
		TokenURL:     s.TokenUrl,
	}
	tokenCtxt := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{})
	return tokenConfig.Client(tokenCtxt)
}
