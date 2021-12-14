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
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type SbiServer struct {
	Server     *echo.Echo
	ListenAddr string
}

func NewSbiServer(listenAddr string) *SbiServer {
	return &SbiServer{
		Server:     echo.New(),
		ListenAddr: listenAddr,
	}
}

func (sbiServer *SbiServer) Start() error {
	errChan := make(chan error)
	addr, err := net.ResolveTCPAddr("tcp", sbiServer.ListenAddr)
	if err != nil {
		return err
	}

	go func() {
		err = sbiServer.Server.Start(addr.String())
		if err != nil {
			errChan <- err
		}
	}()

	return sbiServer.waitForServer(errChan)
}

func (sbiServer *SbiServer) Shutdown(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	sbiServer.Server.Shutdown(ctx)
}

// waitForServer waits for the Echo server to be launched.
func (sbiServer *SbiServer) waitForServer(errChan <-chan error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			addr := sbiServer.Server.ListenerAddr()
			if addr != nil && strings.Contains(addr.String(), ":") {
				// server started
				sbiServer.ListenAddr = addr.String()
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
