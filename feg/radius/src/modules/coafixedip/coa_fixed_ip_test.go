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

package coafixedip

import (
	"context"
	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"
	"fmt"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"
)

func TestCoaFixed(t *testing.T) {
	// Arrange
	secret := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	logger, _ := zap.NewDevelopment()
	mCtx, err := Init(logger, modules.ModuleConfig{
		"target": fmt.Sprintf("127.0.0.1:%d", 4799),
	})
	require.Nil(t, err)

	// Spawn a mock radius server to return response for the coa request
	radiusServer := radius.PacketServer{
		Handler: radius.HandlerFunc(
			func(w radius.ResponseWriter, r *radius.Request) {
				fmt.Println("Got RADIUS packet")
				resp := r.Response(radius.CodeDisconnectACK)
				fmt.Println("Sending RADIUS response")
				w.Write(resp)
			},
		),
		SecretSource: radius.StaticSecretSource(secret),
		Addr:         fmt.Sprintf(":%d", 4799),
		Ready:        make(chan bool, 1),
	}
	fmt.Print("Starting server... ")
	go func() {
		_ = radiusServer.ListenAndServe()
	}()
	defer radiusServer.Shutdown(context.Background())
	listenSuccess := <-radiusServer.Ready // Wait for server to get ready
	if !listenSuccess {
		require.Fail(t, "radiusServer start error")
	}
	fmt.Println("Server listenning")

	// Act
	res, err := Handle(
		mCtx,
		&modules.RequestContext{
			RequestID:      0,
			Logger:         logger,
			SessionStorage: nil,
		},
		createRadiusRequest(),
		func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
			require.Fail(t, "Should never be called (coa fixed module should not call next()")
			return nil, nil
		},
	)

	// Assert
	require.Nil(t, err)
	require.NotNil(t, res)
	require.Equal(t, res.Code, radius.CodeDisconnectACK)
}

func createRadiusRequest() *radius.Request {
	packet := radius.New(radius.CodeDisconnectRequest, []byte{0x01, 0x02, 0x03, 0x4, 0x05, 0x06})
	req := &radius.Request{}
	req = req.WithContext(context.Background())
	req.Packet = packet
	return req
}
