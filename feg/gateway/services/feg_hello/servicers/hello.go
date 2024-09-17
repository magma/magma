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

package servicers

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"magma/feg/cloud/go/protos"
	"magma/orc8r/cloud/go/http2"
)

type helloServer struct{}

func NewFegHelloServer() *helloServer {
	return &helloServer{}
}

func (srv *helloServer) SayHello(ctx context.Context, req *protos.HelloRequest) (*protos.HelloReply, error) {
	glog.Infof("[FeG HELLO] received greeting from AGW: '%s', status: %d (%s)",
		req.GetGreeting(), req.GetGrpcErrCode(), codes.Code(req.GetGrpcErrCode()).String())

	res := getHelloReplyWithTimestamps(req)

	if codes.Code(req.GrpcErrCode) == codes.OK {
		res.Greeting = req.Greeting
		return res, nil
	}
	if req.GrpcErrCode > 16 {
		msg := fmt.Sprintf("requested errorCode %v is out of bound. Valid Range: 0 - 16", req.GrpcErrCode)
		glog.Errorf(msg)
		res.Greeting = req.Greeting
		return res, status.Errorf(codes.OutOfRange, msg)
	}
	res.Greeting = ""
	return res, status.Errorf(
		codes.Code(req.GrpcErrCode), http2.PercentEncode(
			fmt.Sprintf("echo req msg was %v", req.Greeting)))
}

func getHelloReplyWithTimestamps(req *protos.HelloRequest) *protos.HelloReply {
	timestamps := &protos.ResponseTimestamps{
		AgwToFegRelayTimestamp: req.AgwToFegRelayTimestamp,
		FegTimestamp:           timestamppb.Now(),
	}
	return &protos.HelloReply{Timestamps: timestamps}
}
