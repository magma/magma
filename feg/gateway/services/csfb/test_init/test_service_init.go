/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_init

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/csfb/servicers"

	"github.com/ishidawataru/sctp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func GetConnToTestFedGWServiceServer(t *testing.T, connectionInterface servicers.ClientConnectionInterface) *grpc.ClientConn {
	srv, err := servicers.NewCsfbServer(connectionInterface)
	assert.NoError(t, err)

	s := grpc.NewServer()
	protos.RegisterCSFBFedGWServiceServer(s, srv)

	lis, err := net.Listen("tcp", "")
	assert.NoError(t, err)

	go func() {
		err = s.Serve(lis)
		assert.NoError(t, err)
	}()

	addr := lis.Addr()
	conn, err := grpc.Dial(addr.String(), grpc.WithInsecure())
	assert.NoError(t, err)
	return conn
}

func GetMockVLRListenerAndPort(t *testing.T) (*sctp.SCTPListener, int) {
	//gets the default configuration (servicers.DefaultVLRIPAddress)
	config := servicers.GetCsfbConfig()
	ipStr, portNumber, err := servicers.SplitIP(config.Client.ServerAddress)
	assert.Equal(t, servicers.DefaultVLRIPAddress, ipStr)
	assert.Equal(t, servicers.DefaultVLRPort, portNumber)

	ln, err := sctp.ListenSCTP("sctp", servicers.ConstructSCTPAddr(ipStr, 0))
	assert.NoError(t, err)

	port, err := getListenerPort(ln)
	assert.NoError(t, err)

	return ln, port
}

func getListenerPort(listener *sctp.SCTPListener) (int, error) {
	if listener == nil {
		return -1, errors.New("null listener")
	}

	addr := listener.Addr().String()
	arr := strings.Split(addr, ":")
	if len(arr) != 2 {
		return -1, fmt.Errorf("unparsable format of address: %s", addr)
	}

	port, err := strconv.Atoi(arr[1])
	if err != nil {
		return -1, fmt.Errorf("port is not a valid number: %s", arr[1])
	}
	return port, nil
}
