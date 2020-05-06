//
// Copyright (c) Facebook, Inc. and its affiliates.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

// +build !link_local_service

// Package aka implements EAP-AKA provider
package provider

import (
	"errors"
	"fmt"

	"magma/feg/gateway/services/eap/providers"
	"magma/feg/gateway/services/eap/providers/aka/servicers"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/aaa/protos"
	eapp "magma/feg/gateway/services/eap/protos"
	_ "magma/feg/gateway/services/eap/providers/aka/servicers/handlers"
)

// Wrapper to provide a wrapper for GRPC Client to extend it with Cleanup
// functionality
type akaClient struct {
	eapp.EapServiceClient
	cc *grpc.ClientConn
}

func (cl *akaClient) Cleanup() {
	if cl != nil && cl.cc != nil {
		cl.cc.Close()
	}
}

// getAKAClient is a utility function to get a RPC connection to the EAP service
func getAKAClient() (*akaClient, error) {
	conn, err := registry.GetConnection(registry.EAP_AKA)
	if err != nil {
		errMsg := fmt.Sprintf("EAP client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &akaClient{
		eapp.NewEapServiceClient(conn),
		conn,
	}, err
}

// Handle handles passed EAP-AKA payload & returns corresponding result
// this Handle implementation is using GRPC based AKA provider service
func (*providerImpl) Handle(msg *protos.Eap) (*protos.Eap, error) {
	if msg == nil {
		return nil, errors.New("Invalid EAP AKA Message")
	}
	cli, err := getAKAClient()
	if err != nil {
		return nil, err
	}
	return cli.Handle(context.Background(), msg)
}

func NewServiced(_ *servicers.EapAkaSrv) providers.Method {
	return New()
}
