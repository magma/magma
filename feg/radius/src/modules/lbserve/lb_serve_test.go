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

package lbserve

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

func TestLBServeFailsWithNoState(t *testing.T) {
	// Arrange
	var sessionID = "sessionID"

	nextInvoked := false
	logger, _ := zap.NewDevelopment()
	mCtx, _ := Init(logger, modules.ModuleConfig{})

	sessionStorage := session.NewSessionStorage(session.NewMultiSessionMemoryStorage(), sessionID)

	// Act
	_, err := Handle(
		mCtx,
		&modules.RequestContext{
			RequestID:      0,
			Logger:         logger,
			SessionStorage: sessionStorage,
		},
		createRadiusRequest(),
		func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
			nextInvoked = true
			return nil, nil
		},
	)

	// Assert
	require.NotNil(t, err)
	require.False(t, nextInvoked)
}

func TestLBServeFailsWithNoUpstreamHost(t *testing.T) {
	// Arrange
	var sessionID = "sessionID"

	nextInvoked := false
	logger, _ := zap.NewDevelopment()
	mCtx, _ := Init(logger, modules.ModuleConfig{})

	sessionStorage := session.NewSessionStorage(session.NewMultiSessionMemoryStorage(), sessionID)
	sessionStorage.Set(session.State{})

	// Act
	_, err := Handle(
		mCtx,
		&modules.RequestContext{
			RequestID:      0,
			Logger:         logger,
			SessionStorage: sessionStorage,
		},
		createRadiusRequest(),
		func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
			nextInvoked = true
			return nil, nil
		},
	)

	// Assert
	require.NotNil(t, err)
	require.False(t, nextInvoked)
}

func TestLBServeProxiesRequestToRadiusAndReturnsResponse(t *testing.T) {
	// Arrange
	var sessionID = "sessionID"

	nextInvoked := false
	logger, _ := zap.NewDevelopment()
	mCtx, _ := Init(logger, modules.ModuleConfig{})

	// Spawn a radius server
	server, port := spawnRadiusServer()
	defer server.Shutdown(context.Background())

	sessionStorage := session.NewSessionStorage(session.NewMultiSessionMemoryStorage(), sessionID)
	sessionStorage.Set(session.State{UpstreamHost: fmt.Sprintf("localhost:%d", port)})

	// Act
	response, err := Handle(
		mCtx,
		&modules.RequestContext{
			RequestID:      0,
			Logger:         logger,
			SessionStorage: sessionStorage,
		},
		createRadiusRequest(),
		func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
			nextInvoked = true
			return nil, nil
		},
	)
	attr, ok := response.Attributes[rfc2865.State_Type]

	// Assert
	require.Nil(t, err)
	require.False(t, nextInvoked)
	require.True(t, ok)
	require.NotNil(t, attr)
	require.Equal(t, "server_returned_value", string(attr[0]))
}

func spawnRadiusServer() (server *radius.PacketServer, port int) {
	secret := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	port = int((rand.Int63() % 0xFFF) << 4)
	server = &radius.PacketServer{
		Handler: radius.HandlerFunc(
			func(w radius.ResponseWriter, r *radius.Request) {
				fmt.Println("Got RADIUS packet")
				resp := r.Response(radius.CodeAccessAccept)
				resp.Add(rfc2865.State_Type, []byte("server_returned_value"))
				w.Write(resp)
			},
		),
		SecretSource: radius.StaticSecretSource(secret),
		Addr:         fmt.Sprintf(":%d", port),
		Ready:        make(chan bool, 1),
	}
	fmt.Print("Starting server... ")
	go func() {
		_ = server.ListenAndServe()
	}()
	<-server.Ready // Wait for server to get ready
	fmt.Println("Server listening")
	return
}

func TestLBServeFailsWithRadiusError(t *testing.T) {
	// Arrange
	var sessionID = "sessionID"

	nextInvoked := false
	logger, _ := zap.NewDevelopment()
	mCtx, _ := Init(logger, modules.ModuleConfig{})

	// Spawn a radius server
	server, port := spawnFailingRadiusServer()
	defer server.Shutdown(context.Background())

	sessionStorage := session.NewSessionStorage(session.NewMultiSessionMemoryStorage(), sessionID)
	sessionStorage.Set(session.State{UpstreamHost: fmt.Sprintf("localhost:%d", port)})

	// Act
	response, err := Handle(
		mCtx,
		&modules.RequestContext{
			RequestID:      0,
			Logger:         logger,
			SessionStorage: sessionStorage,
		},
		createRadiusRequest(),
		func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
			nextInvoked = true
			return nil, nil
		},
	)
	attr, ok := response.Attributes[rfc2865.State_Type]

	// Assert
	require.Nil(t, err)
	require.False(t, nextInvoked)
	require.True(t, ok)
	require.NotNil(t, attr)
	require.Equal(t, "server_returned_value", string(attr[0]))
}

func spawnFailingRadiusServer() (server *radius.PacketServer, port int) {
	secret := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	port = int((rand.Int63() % 0xFFF) << 4)
	server = &radius.PacketServer{
		Handler: radius.HandlerFunc(
			func(w radius.ResponseWriter, r *radius.Request) {
				fmt.Println("Got RADIUS packet")
				resp := r.Response(radius.CodeAccessAccept)
				resp.Add(rfc2865.State_Type, []byte("server_returned_value"))
				w.Write(resp)
			},
		),
		SecretSource: radius.StaticSecretSource(secret),
		Addr:         fmt.Sprintf(":%d", port),
		Ready:        make(chan bool, 1),
	}
	fmt.Print("Starting server... ")
	go func() {
		_ = server.ListenAndServe()
	}()
	<-server.Ready // Wait for server to get ready
	fmt.Println("Server listening")
	return
}

func createRadiusRequest() *radius.Request {
	packet := radius.New(radius.CodeAccessRequest, []byte{0x01, 0x02, 0x03, 0x4, 0x05, 0x06})
	packet.Attributes[rfc2865.CallingStationID_Type] = []radius.Attribute{radius.Attribute("called")}
	packet.Attributes[rfc2865.CalledStationID_Type] = []radius.Attribute{radius.Attribute("calling")}
	req := &radius.Request{}
	req = req.WithContext(context.Background())
	req.Packet = packet
	return req
}
