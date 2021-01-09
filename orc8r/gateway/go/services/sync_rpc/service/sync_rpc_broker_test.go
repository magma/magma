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

// package service_test
package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang/glog"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"magma/orc8r/lib/go/protos"
)

func parseHeaders(hdr http.Header) map[string]string {
	x := make(map[string]string)
	for k, v := range hdr {
		x[k] = strings.Join(v, ",")
	}
	return x
}

func getRandPayload(sz int) []byte {
	b := make([]byte, sz)
	rand.Read(b)
	return []byte(base64.URLEncoding.EncodeToString(b))
}

// magmadTestServer mocks a test grpc service running on the gateway
type magmadTestConfig struct {
	gatewayID             string
	logLines              string
	restartServiceseDelay time.Duration
	startServicesErr      error
}

type magmadTestServer struct {
	protos.UnimplementedMagmadServer
	tc magmadTestConfig
}

func (m *magmadTestServer) StartServices(context.Context, *protos.Void) (*protos.Void, error) {
	return nil, m.tc.startServicesErr
}

func (m *magmadTestServer) StopServices(context.Context, *protos.Void) (*protos.Void, error) {
	return nil, nil
}

func (m *magmadTestServer) Reboot(context.Context, *protos.Void) (*protos.Void, error) {
	return nil, nil
}

func (m *magmadTestServer) RestartServices(context.Context, *protos.RestartServicesRequest) (*protos.Void, error) {
	time.Sleep(m.tc.restartServiceseDelay)
	return nil, nil
}

func (m *magmadTestServer) GetConfigs(context.Context, *protos.Void) (*protos.GatewayConfigs, error) {
	return nil, nil
}

func (m *magmadTestServer) SetConfigs(context.Context, *protos.GatewayConfigs) (*protos.Void, error) {
	return nil, nil
}

func (m *magmadTestServer) RunNetworkTests(context.Context, *protos.NetworkTestRequest) (*protos.NetworkTestResponse, error) {
	return nil, nil
}

func (m *magmadTestServer) GenericCommand(context.Context, *protos.GenericCommandParams) (*protos.GenericCommandResponse, error) {
	return nil, nil
}

func (m *magmadTestServer) GetGatewayId(context.Context, *protos.Void) (*protos.GetGatewayIdResponse, error) {
	resp := &protos.GetGatewayIdResponse{
		GatewayId: m.tc.gatewayID,
	}
	return resp, nil
}

func (m *magmadTestServer) TailLogs(req *protos.TailLogsRequest, s protos.Magmad_TailLogsServer) error {
	if req.Service != "TestService" {
		return nil
	}
	r := strings.NewReader(m.tc.logLines)
	b := make([]byte, 1000)
	for {
		len, err := r.Read(b)
		if err != nil {
			break
		}
		s.SendMsg(&protos.LogLine{Line: string(b[:len])})
	}
	return nil
}

// run instance of the test grpc service
func runTestMagmadServer(server *magmadTestServer, grpcPortCh chan string, stopCh chan struct{}) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":0"))
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}

	v := strings.Split(lis.Addr().String(), ":")
	grpcPortCh <- v[len(v)-1]
	grpcServer := grpc.NewServer()
	protos.RegisterMagmadServer(grpcServer, server)
	go func() {
		grpcServer.Serve(lis)
	}()
	<-stopCh
	grpcServer.GracefulStop()
}

// testBroker creates a test Broker and uses the brokerImpl to verify
// it's send method
type testBroker struct {
	cfg      *Config
	reqID    uint32
	grpcAddr string
	h        http.HandlerFunc

	numKeepaliveConnReqs uint32
}

func newTestBroker(cfg *Config, grpcPortCh chan string) *testBroker {
	t := &testBroker{
		cfg:      cfg,
		grpcAddr: fmt.Sprintf("localhost:%s", <-grpcPortCh),
	}
	t.h = t.handler
	return t
}

func (t *testBroker) handler(w http.ResponseWriter, req *http.Request) {
	payload, _ := ioutil.ReadAll(req.Body)
	ctx := context.Background()
	testSyncRpcReq := &protos.SyncRPCRequest{
		ReqId: t.reqID,
		ReqBody: &protos.GatewayRequest{
			Authority: req.Host,
			Path:      req.URL.Path,
			Payload:   payload,
			Headers:   parseHeaders(req.Header),
		},
	}

	syncRpcRespCh := make(chan *protos.SyncRPCResponse)
	go func() {
		// This verifies the Broker request mechanism
		p := newbrokerImpl(t.cfg)
		p.send(ctx, t.grpcAddr, testSyncRpcReq, syncRpcRespCh)
	}()

	// send this back to the client
	for {
		resp := <-syncRpcRespCh
		gatewayResp := resp.RespBody
		for k, v := range gatewayResp.Headers {
			w.Header().Add(k, v)
		}
		if gatewayResp.KeepConnActive {
			// ignore and move on
			t.numKeepaliveConnReqs++
			continue
		}

		if len(gatewayResp.Err) == 0 {
			status, err := strconv.Atoi(gatewayResp.Status)
			if err != nil {
				panic(fmt.Errorf("error converting status to int %v", err))
			}

			w.Header().Set("Trailer", "Grpc-Status, Grpc-Message")
			w.WriteHeader(status)
			if gatewayResp.Payload != nil {
				w.Write(gatewayResp.Payload)
			}
			w.(http.Flusher).Flush()

			if _, ok := gatewayResp.Headers["Grpc-Status"]; ok {
				break
			}
		} else {
			w.Header().Set("content-type", "application/grpc")
			w.Header().Set("trailer", "Grpc-Status")
			w.Header().Add("trailer", "Grpc-Message")
			w.Header().Set("grpc-status", "13")
			w.Header().Set("grpc-message", gatewayResp.Err)
			w.WriteHeader(200)
			return
		}
	}
	t.reqID++
	return
}

func runTestBroker(Broker *testBroker, BrokerPortCh chan string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":0"))
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}

	v := strings.Split(lis.Addr().String(), ":")
	BrokerPortCh <- v[len(v)-1]
	srv := &http2.Server{}
	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Printf("l.accept err: %v\n", err)
		}
		go srv.ServeConn(conn, &http2.ServeConnOpts{
			Handler: Broker.h,
		})
	}

}

// Following test creates a testBroker and test grpc service and we verify if
// a grpc client's unary and streaming request works as expected
func TestBrokerSanity(t *testing.T) {
	serverStopCh := make(chan struct{})
	grpcPortCh := make(chan string)
	tc := magmadTestConfig{
		gatewayID:             "Test_Gateway",
		logLines:              string(getRandPayload(10 * 1000)),
		restartServiceseDelay: 30 * time.Second,
		startServicesErr:      errors.New("failed starting test service"),
	}

	// start a test gateway service
	gwServer := &magmadTestServer{tc: tc}
	go runTestMagmadServer(gwServer, grpcPortCh, serverStopCh)

	// start a Broker which accepts grpc client reqs and
	// proxies it to the test gateway service
	cfg := &Config{
		GatewayKeepaliveInterval: time.Second,
		GatewayResponseTimeout:   2 * time.Second,
	}
	Broker := newTestBroker(cfg, grpcPortCh)
	BrokerPortCh := make(chan string)
	go runTestBroker(Broker, BrokerPortCh)

	//create a grpc client and connect to the Broker
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", <-BrokerPortCh),
		grpc.WithInsecure())
	if err != nil {
		t.Fatal("Failed creating a test client")
		return
	}
	defer conn.Close()
	client := protos.NewMagmadClient(conn)

	// TC1 - test unary request with the Broker
	v, _ := client.GetGatewayId(context.Background(), &protos.Void{})
	assert.Equal(t, tc.gatewayID, v.GatewayId)

	// TC2 - test streaming request with the Broker
	v2Stream, err := client.TailLogs(context.Background(), &protos.TailLogsRequest{Service: "TestService"})
	if err != nil {
		t.Fatalf("err %v grpc request failed", err)
		return
	}
	var actualLogs []byte
	for {
		resp, err := v2Stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatalf("error when handling stream %v", err)
			break
		}
		actualLogs = append(actualLogs, []byte(resp.Line)...)
	}
	assert.Equal(t, tc.logLines, string(actualLogs))

	// TC3 - invoke a method on grpc client which returns an error
	_, err = client.StartServices(context.Background(), &protos.Void{})
	sts, _ := status.FromError(err)
	assert.Equal(t, tc.startServicesErr.Error(), sts.Message())

	// // TC4 - test response timeout
	_, err = client.RestartServices(context.Background(), &protos.RestartServicesRequest{})
	sts, _ = status.FromError(err)
	assert.GreaterOrEqual(t, Broker.numKeepaliveConnReqs, uint32(1))
	assert.NotEmpty(t, sts.Message())

	// TC5 - stop magmad service and check if the client recvs error
	serverStopCh <- struct{}{}

	v, err = client.GetGatewayId(context.Background(), &protos.Void{})
	sts, _ = status.FromError(err)
	assert.Contains(t, sts.Message(), "connection refused")
}
