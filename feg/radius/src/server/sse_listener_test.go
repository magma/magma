package server

import (
	"encoding/json"
	"errors"
	"fbc/cwf/radius/config"
	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/monitoring"
	"fbc/cwf/radius/session"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2866"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/donovanhide/eventsource"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// CoAServer a CoA server combines SSE and Http for CoA Req streaming + Resp sending
type CoAServer struct {
	HTTPListener net.Listener
	SSEServer    *eventsource.Server
}

func (s *CoAServer) Close() {
	s.HTTPListener.Close()
	s.SSEServer.Close()
}

// CoAEvent a signle SSE of CoA
type CoAEvent SSEEvent

func (t CoAEvent) Id() string    { return (string)(t.Identifier) }
func (t CoAEvent) Event() string { return "coa" }
func (t CoAEvent) Data() string {
	res, err := json.Marshal(t)
	if err != nil {
		return "{}"
	}
	return string(res)
}

func TestCoARequestResponse(t *testing.T) {
	// Arrange
	port := 8080
	srv, err := createCoaServer(port)
	require.Nil(t, err)

	testError := make(chan error, 1)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testError <- nil
	})
	go http.ListenAndServe(fmt.Sprintf(":%d", port+1), nil)
	sseListener := NewSSEListener()
	logger, err := zap.NewDevelopment()
	require.Nil(t, err)
	sseListener.Init(
		&Server{logger: logger, multiSessionStorage: session.NewMultiSessionMemoryStorage()},
		config.ServerConfig{},
		config.ListenerConfig{
			Extra: map[string]interface{}{
				"EventStreamURL": fmt.Sprintf("http://127.0.0.1:%d/coa", port),
				"ResponseURL":    fmt.Sprintf("http://127.0.0.1:%d/coa_response", port+1),
			},
		},
		monitoring.CreateListenerCounters("test_listener"),
	)
	sseListener.SetHandleRequest(
		func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
			require.Equal(t, "session_id", string(r.Get(rfc2866.AcctSessionID_Type)))
			return &modules.Response{
				Code:       radius.CodeAccessAccept,
				Attributes: radius.Attributes{},
			}, nil
		},
	)

	go sseListener.ListenAndServe()
	<-sseListener.Ready()

	timeoutTimer := time.NewTimer(time.Second)
	go func() {
		<-timeoutTimer.C
		testError <- errors.New("test timedout")
	}()

	// Act
	srv.SSEServer.Publish(
		[]string{"coa"},
		&CoAEvent{
			Code:       43,
			Identifier: 123,
			AVPs: map[string][]interface{}{
				"Acct-Session-Id":    []interface{}{"session_id"},
				"User-Name":          []interface{}{"username"},
				"Calling-Station-Id": []interface{}{"00:11:22:33:44:55"},
				"XWF-Authorize-Traffic-Classes": []interface{}{
					map[string]interface{}{
						"XWF-Authorize-Class-Name": "xwf",
						"XWF-Authorize-Bytes-Left": 0,
					},
					map[string]interface{}{
						"XWF-Authorize-Class-Name": "fbs",
						"XWF-Authorize-Bytes-Left": 0,
					},
					map[string]interface{}{
						"XWF-Authorize-Class-Name": "internet",
						"XWF-Authorize-Bytes-Left": 9999,
					},
				},
			},
			ProxyState: "proxy-state",
		},
	)

	// Wait for completion & teardown
	err = <-testError
	srv.Close()

	// Assert
	require.Nil(t, err)
}

func createCoaServer(port int) (*CoAServer, error) {
	sseSrv := eventsource.NewServer()
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("Failed to listen on port %d", port)
	}
	http.HandleFunc("/coa", sseSrv.Handler("coa"))
	go http.Serve(l, nil)
	return &CoAServer{
		HTTPListener: l,
		SSEServer:    sseSrv,
	}, nil
}
