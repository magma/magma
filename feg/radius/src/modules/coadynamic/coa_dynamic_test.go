package coadynamic

import (
	"context"
	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2866"
	"fmt"
	"net"
	"sync/atomic"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"
)

func TestCoaDynamic(t *testing.T) {
	// Arrange
	secret := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	logger, _ := zap.NewDevelopment()
	ctx, err := Init(logger, modules.ModuleConfig{
		"port": 4799,
	})
	require.Nil(t, err)

	// Spawn a mock radius server to return response for the coa request
	var radiusResponseCounter uint32
	radiusServer := radius.PacketServer{
		Handler: radius.HandlerFunc(
			func(w radius.ResponseWriter, r *radius.Request) {
				atomic.AddUint32(&radiusResponseCounter, 1)
				resp := r.Response(radius.CodeDisconnectACK)
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

	// Sending a coa request - expected to fail
	generateRequest(ctx, radius.CodeDisconnectRequest, t, "session1", false)
	require.Equal(t, uint32(1), atomic.LoadUint32(&radiusResponseCounter))

	// Sending a non coa request
	generateRequest(ctx, radius.CodeAccountingRequest, t, "session2")
	require.Equal(t, uint32(1), atomic.LoadUint32(&radiusResponseCounter))

	// Sending a coa request
	res, err := generateRequest(ctx, radius.CodeDisconnectRequest, t, "session3", false)
	require.Equal(t, uint32(2), atomic.LoadUint32(&radiusResponseCounter))

	// Assert
	require.Nil(t, err)
	require.NotNil(t, res)
	require.Equal(t, res.Code, radius.CodeDisconnectACK)
}

func generateRequest(ctx modules.Context, code radius.Code, t *testing.T, sessionID string, next ...bool) (*modules.Response, error) {
	logger, _ := zap.NewDevelopment()
	nextCalled := false

	// Update tracker with some target endpoint
	tracker := GetRadiusTracker()
	tracker.Set(&radius.Request{
		Packet: &radius.Packet{
			Attributes: radius.Attributes{
				rfc2866.AcctSessionID_Type: []radius.Attribute{radius.Attribute(sessionID)},
			},
		},
		RemoteAddr: IPAddr{"127.0.0.1:1313"},
	})

	// Handle
	res, err := Handle(
		ctx,
		&modules.RequestContext{
			RequestID:      0,
			Logger:         logger,
			SessionStorage: nil,
		},
		createRadiusRequest(code, sessionID),
		func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
			nextCalled = true
			return nil, nil
		},
	)

	// Verify
	nextCalledExpected := true
	if len(next) > 0 {
		nextCalledExpected = next[0]
	}
	require.Equal(t, nextCalledExpected, nextCalled)

	return res, err
}

func createRadiusRequest(code radius.Code, sessionID string) *radius.Request {
	packet := radius.New(code, []byte{0x01, 0x02, 0x03, 0x4, 0x05, 0x06})
	packet.Attributes[rfc2866.AcctSessionID_Type] = []radius.Attribute{radius.Attribute(sessionID)}
	req := &radius.Request{}
	req.RemoteAddr = &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 4799}
	req = req.WithContext(context.Background())
	req.Packet = packet
	return req
}

// IPAddr
type IPAddr struct{ IP string }

func (a IPAddr) Network() string { return "ip" }
func (a IPAddr) String() string  { return a.IP }
