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

package proxy

import (
	"context"
	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/session"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"
	"fmt"
	"math/rand"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"
)

func TestProxy(t *testing.T) {
	// Arrange
	var sessionID = "sessionID"
	randomPort := (rand.Int63() % 0xFFF) << 4
	secret := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "failed to get logger")
	mCtx, _ := Init(logger, modules.ModuleConfig{
		"target": fmt.Sprintf("localhost:%d", randomPort),
	})

	// Spawn a radius server
	radiusServer := radius.PacketServer{
		Handler: radius.HandlerFunc(
			func(w radius.ResponseWriter, r *radius.Request) {
				fmt.Println("Got RADIUS packet")
				resp := r.Response(radius.CodeAccessAccept)
				resp.Add(rfc2865.State_Type, []byte("server_returned_value"))
				w.Write(resp)
			},
		),
		SecretSource: radius.StaticSecretSource(secret),
		Addr:         fmt.Sprintf(":%d", randomPort),
		Ready:        make(chan bool, 1),
	}
	fmt.Print("Starting server... ")
	go func() {
		_ = radiusServer.ListenAndServe()
	}()
	defer radiusServer.Shutdown(context.Background())
	listenSuccess := <-radiusServer.Ready // Wait for server to get ready
	if !listenSuccess {
		return
	}
	fmt.Println("Server listenning")

	// Act
	sessionStorage := session.NewSessionStorage(session.NewMultiSessionMemoryStorage(), sessionID)
	res, err := Handle(
		mCtx,
		&modules.RequestContext{
			RequestID:      0,
			Logger:         logger,
			SessionStorage: sessionStorage,
		},
		createRadiusRequest("called", "calling"),
		func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
			require.Fail(t, "Should never be called (proxy module should not call next()")
			return nil, nil
		},
	)

	// Assert
	require.Nil(t, err)
	require.NotNil(t, res)
	require.NotNil(t, res.Attributes)
	attr, ok := res.Attributes[rfc2865.State_Type]
	require.True(t, ok)
	require.NotNil(t, attr)
	require.Equal(t, "server_returned_value", string(attr[0]))
}

func TestInvalidConfig(t *testing.T) {
	// Arrange
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "failed to get logger")
	_, err = Init(logger, modules.ModuleConfig{
		"notneeded": "config",
		// Missing "target" (on purpose, this is what we're testing!)
	})

	// Assert
	require.NotNil(t, err)
	require.Equal(t, "proxy module cannot be initialize with empty Target value", err.Error())
}

func createRadiusRequest(calledStationID string, callingStationID string) *radius.Request {
	packet := radius.New(radius.CodeAccessRequest, []byte{0x01, 0x02, 0x03, 0x4, 0x05, 0x06})
	packet.Attributes[rfc2865.CallingStationID_Type] = []radius.Attribute{radius.Attribute(callingStationID)}
	packet.Attributes[rfc2865.CalledStationID_Type] = []radius.Attribute{radius.Attribute(calledStationID)}
	req := &radius.Request{}
	req = req.WithContext(context.Background())
	req.Packet = packet
	return req
}
