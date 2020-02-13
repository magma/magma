/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package mocks

import (
	"context"
	"fmt"
	"log"
	"time"

	"magma/orc8r/lib/go/service/config"

	"google.golang.org/grpc"
)

type MockCloudRegistry struct {
	ServerAddr string
}

// Mocked implementation which returns a grpc connection to the `ServerAddr`
// field in the struct.
func (m *MockCloudRegistry) GetCloudConnection(service string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, m.ServerAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Err: %s", err)
		return nil, fmt.Errorf("Address: %s GRPC Dial error: %s", m.ServerAddr, err)
	} else if ctx.Err() != nil {
		log.Printf("Err: %s", ctx.Err())
		return nil, ctx.Err()
	}
	return conn, nil
}

func (*MockCloudRegistry) GetCloudConnectionFromServiceConfig(serviceConfig *config.ConfigMap, service string) (*grpc.ClientConn, error) {
	return nil, nil
}
