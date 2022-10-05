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

package httpserver_test

// Everything commented out until we can write some tests that don't depend
// on mobilityd, which is an lte service

//
//func TestSyncRPCHttpServerInvalidArguments(t *testing.T) {
//	addr, broker := test_init.StartTestHttpServer(t)
//	conn := getClientConnection(addr.String(), t, "mobility")
//	defer conn.Close()
//	client := protos.NewMobilityServiceClient(conn)
//	in := new(protos.Void)
//
//	// test with no gatewayid set expects an InvalidArgument error
//	_, err := client.ListAddedIPv4Blocks(context.Background(), in)
//	assert.EqualError(t, err,
//		"rpc error: code = InvalidArgument desc = No Gatewayid provided in metaData")
//	broker.AssertNotCalled(t, "SendRequestToGateway", mock.Anything)
//	broker.AssertExpectations(t)
//}
//
//func TestSyncRPCHttpServerInternalErr(t *testing.T) {
//	addr, broker := test_init.StartTestHttpServer(t)
//	// test with SendRequestToGateway returning err expects an INTERNAL error.
//	err := errors.New("test error")
//	broker.On("SendRequestToGateway", mock.AnythingOfType("*protos.GatewayRequest")).Return(nil, err)
//	conn := getClientConnection(addr.String(), t, "mobility")
//	defer conn.Close()
//	client := protos.NewMobilityServiceClient(conn)
//	in := new(protos.Void)
//	customHeader := metadata.New(map[string]string{"gatewayId": "gatewayHwIdVal"})
//	ctx := metadata.NewOutgoingContext(context.Background(), customHeader)
//	_, err = client.ListAddedIPv4Blocks(ctx, in)
//	expectedErrPrefix := "err sending request gwId:\"gatewayHwIdVal\" authority:\"mobility\" " +
//		"path:\"/magma.MobilityService/ListAddedIPv4Blocks\""
//	assert.Regexp(t,
//		regexp.MustCompile(fmt.Sprintf("rpc error: code = Internal desc = %v.*", expectedErrPrefix)),
//		err.Error())
//}
//
//func TestSyncRPCHttpServerTimeout(t *testing.T) {
//	addr, broker := test_init.StartTestHttpServer(t)
//	respChan := make(chan *protos.GatewayResponse)
//	broker.On("SendRequestToGateway", mock.AnythingOfType("*protos.GatewayRequest")).
//		Return(respChan, nil)
//	conn := getClientConnection(addr.String(), t, "mobility")
//	defer conn.Close()
//	client := protos.NewMobilityServiceClient(conn)
//	in := new(protos.Void)
//	customHeader := metadata.New(map[string]string{"gatewayId": "gatewayHwIdVal"})
//	ctx := metadata.NewOutgoingContext(context.Background(), customHeader)
//	_, err := client.ListAddedIPv4Blocks(ctx, in)
//	// if no response is sent to the respChan, expects a DeadlineExceeded error
//	assert.EqualError(t, err,
//		"rpc error: code = DeadlineExceeded desc = Request timed out")
//}
//
//func TestSyncRPCHttpServer_Success(t *testing.T) {
//	addr, broker := test_init.StartTestHttpServer(t)
//	respChan := make(chan *protos.GatewayResponse, 1)
//	broker.On("SendRequestToGateway", mock.AnythingOfType("*protos.GatewayRequest")).
//		Return(respChan, nil)
//	conn := getClientConnection(addr.String(), t, "mobility")
//	defer conn.Close()
//	client := protos.NewMobilityServiceClient(conn)
//	in := new(protos.Void)
//	customHeader := metadata.New(map[string]string{"gatewayId": "gatewayHwIdVal"})
//	ctx := metadata.NewOutgoingContext(context.Background(), customHeader)
//	gwResp := &protos.GatewayResponse{
//		Status:  "200",
//		Headers: map[string]string{"content-type": "application/grpc", "grpc-status": "0", "grpc-message": ""},
//		Payload: []byte{0, 0, 0, 0, 0},
//	}
//	go func() {
//		respChan <- gwResp
//		close(respChan)
//	}()
//	_, err := client.ListAddedIPv4Blocks(ctx, in)
//	assert.NoError(t, err)
//	broker.AssertExpectations(t)
//
//}
//
//func getClientConnection(addr string, t *testing.T, serviceName string) *grpc.ClientConn {
//	conn, err := grpc.DialContext(context.Background(), addr, grpc.WithInsecure(), grpc.WithAuthority(serviceName))
//	if err != nil {
//		t.Errorf("failed to connect to service %v via SyncRPCHttpServer: %v\n", serviceName, err)
//	}
//	return conn
//}
