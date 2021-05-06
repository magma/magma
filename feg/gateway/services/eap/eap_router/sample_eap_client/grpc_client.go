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

// Package main implements sample eap_router service client
package main

import (
	"context"
	"flag"
	"log"
	"reflect"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc"

	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
)

var (
	initResp = []byte("\x02\x00\x00\x38\x01\x30\x30\x30\x31\x30\x31\x30\x30\x30\x30\x30" +
		"\x30\x30\x30\x35\x35\x40\x77\x6c\x61\x6e\x2e\x6d\x6e\x63\x30\x30\x31\x2e\x6d\x63" +
		"\x63\x30\x30\x31\x2e\x33\x67\x70\x70\x6e\x65\x74\x77\x6f\x72\x6b\x2e\x6f\x72\x67")
	expectedInitReq = []byte{
		eap.RequestCode,
		1,
		0, 12, // EAP Len
		23,
		5,
		0, 0,
		10,
		1,
		0, 0}
)

// To test:
//	from magma/feg/gateway/ run:
//		make run
//      then
//		go run magma/feg/gateway/services/eap/eap_router/sample_eap_client
func main() {
	serverAddr := flag.String("addr", "localhost:9109", "eap_router server address (host:port)")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	glog.Errorf("Dailing EAP Router at: %s", *serverAddr)

	conn, err := grpc.DialContext(ctx, *serverAddr,
		grpc.WithBackoffMaxDelay(10*time.Second), grpc.WithBlock(), grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Client dial error: %v", err)
		return
	}
	client := protos.NewEapRouterClient(conn)
	grpcCtx := context.Background()

	methods, err := client.SupportedMethods(grpcCtx, &protos.Void{})

	if err != nil {
		glog.Fatalf("SuportedMethods error: %v", err)
	}
	glog.Infof("Supported EAP Methods: %v\n", methods.Methods)

	glog.Infof("Sending  EAP: %v\n", initResp)
	res, err := client.HandleIdentity(grpcCtx, &protos.EapIdentity{Payload: initResp, Method: 23})
	if err != nil {
		glog.Fatalf("HandleIdentity error: %v", err)
	}
	if !reflect.DeepEqual(res.GetPayload(), expectedInitReq) {
		glog.Fatalf(
			"Unexpected identity Request received\n\tReceived: %.3v\n\tExpected: %.3v",
			res.GetPayload(), expectedInitReq)
	}
	glog.Infof("Received EAP: %v\n", res.GetPayload())
	conn.Close()
}
