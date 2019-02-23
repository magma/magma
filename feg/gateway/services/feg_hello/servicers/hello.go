/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"fmt"

	"magma/feg/cloud/go/protos"
	"magma/orc8r/cloud/go/http2"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type helloServer struct{}

func NewFegHelloServer() *helloServer {
	return &helloServer{}
}

func (srv *helloServer) SayHello(ctx context.Context, req *protos.HelloRequest) (*protos.HelloReply, error) {
	if codes.Code(req.GrpcErrCode) == codes.OK {
		return &protos.HelloReply{Greeting: req.Greeting}, nil
	}
	if req.GrpcErrCode > 16 {
		msg := fmt.Sprintf("requested errorCode %v is out of bound. Valid Range: 0 - 16", req.GrpcErrCode)
		glog.Errorf(msg)
		return &protos.HelloReply{Greeting: req.Greeting}, status.Errorf(codes.OutOfRange, msg)
	}
	return &protos.HelloReply{Greeting: ""},
		status.Errorf(codes.Code(req.GrpcErrCode), http2.PercentEncode(fmt.Sprintf("echo req msg was %v", req.Greeting)))

}
