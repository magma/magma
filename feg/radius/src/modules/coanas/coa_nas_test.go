package coanas

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"

	"fbc/cwf/radius/modules"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

func TestCoaNas(t *testing.T) {
	// Arrange
	secret := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	port := 4799
	logger, err := zap.NewDevelopment()
	require.Nil(t, err)
	mCtx, err := Init(logger, modules.ModuleConfig{
		"port": fmt.Sprint(port),
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
		Addr:         fmt.Sprintf(":%d", port),
	}
	fmt.Print("Starting server... ")
	go func() {
		_ = radiusServer.ListenAndServe()
	}()
	defer radiusServer.Shutdown(context.Background())
	err = modules.WaitForRadiusServerToBeReady(secret, fmt.Sprint(port))
	require.Nil(t, err)
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
	port := 3799
	logger, err := zap.NewDevelopment()
	require.Nil(t, err)
	mCtx, err := Init(logger, modules.ModuleConfig{
		"port": fmt.Sprint(port),
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
	}
	fmt.Print("Starting server... ")
	go func() {
		_ = radiusServer.ListenAndServe()
	}()
	defer radiusServer.Shutdown(context.Background())
	err = modules.WaitForRadiusServerToBeReady(secret, fmt.Sprint(port))
	require.Nil(t, err)
	fmt.Println("Server listening")

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
	port := 4799
	logger, err := zap.NewDevelopment()
	require.Nil(t, err)
	mCtx, err := Init(logger, modules.ModuleConfig{
		"port": fmt.Sprint(port),
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
		Addr:         fmt.Sprintf(":%d", port),
	}
	fmt.Print("Starting server... ")
	go func() {
		_ = radiusServer.ListenAndServe()
	}()
	defer radiusServer.Shutdown(context.Background())
	err = modules.WaitForRadiusServerToBeReady(secret, fmt.Sprint(port))
	require.Nil(t, err)
	fmt.Println("Server listening")

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
