package coanas

import (
	"context"
	"errors"
	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"
	"fmt"
	"net"
	"testing"

	"go.uber.org/zap"

	"fbc/lib/go/radius/rfc2865"

	"github.com/stretchr/testify/require"
)

func TestCoaNas(t *testing.T) {
	// Arrange
	logger, err := zap.NewDevelopment()
	require.Nil(t, err)
	secret := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	mCtx, err := Init(logger, modules.ModuleConfig{
		"port": "4799",
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
		return
	}
	fmt.Println("Server listenning")

	// Act
	require.Nil(t, err)
	res, err := Handle(
		mCtx,
		&modules.RequestContext{
			RequestID:      0,
			Logger:         logger,
			SessionStorage: nil,
		},
		createRadiusRequest("127.0.0.1"),
		func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
			require.Fail(t, "Should never be called (coa nas module should not call next()")
			return nil, nil
		},
	)

	// Assert
	require.Nil(t, err)
	require.NotNil(t, res)
	require.Equal(t, res.Code, radius.CodeDisconnectACK)
}

func TestCoaNasNoResponse(t *testing.T) {
	// Arrange
	secret := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	logger, err := zap.NewDevelopment()
	require.Nil(t, err)
	mCtx, err := Init(logger, modules.ModuleConfig{
		"port": "3799",
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
		return
	}
	fmt.Println("Server listenning")

	// Act
	_, err = Handle(
		mCtx,
		&modules.RequestContext{
			RequestID:      0,
			Logger:         logger,
			SessionStorage: nil,
		},
		createRadiusRequest("127.0.0.1"),
		func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
			return nil, errors.New("error")
		},
	)

	// Assert
	require.NotNil(t, err)
}

func TestCoaNasFieldInvalid(t *testing.T) {
	// Arrange
	secret := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	logger, err := zap.NewDevelopment()
	require.Nil(t, err)
	mCtx, err := Init(logger, modules.ModuleConfig{
		"port": "4799",
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
		return
	}
	fmt.Println("Server listenning")

	// Act
	var s int
	_, err = Handle(
		mCtx,
		&modules.RequestContext{
			RequestID:      0,
			Logger:         logger,
			SessionStorage: nil,
		},
		createRadiusRequest(""),
		func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
			s = 1
			return nil, errors.New("error")
		},
	)

	require.Equal(t, 1, s)
	// Assert
	require.NotNil(t, err)
}

func createRadiusRequest(nasIdentifier string) *radius.Request {
	packet := radius.New(radius.CodeDisconnectRequest, []byte{0x01, 0x02, 0x03, 0x4, 0x05, 0x06})
	rfc2865.NASIPAddress_Add(packet, net.ParseIP(nasIdentifier))
	req := &radius.Request{}
	req = req.WithContext(context.Background())
	req.Packet = packet
	return req
}
