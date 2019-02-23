/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"flag"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/csfb"
	"magma/feg/gateway/services/csfb/servicers"
	"magma/feg/gateway/services/csfb/servicers/decode"
	"magma/feg/gateway/services/csfb/servicers/decode/message"
	"magma/orc8r/cloud/go/service"

	"github.com/golang/glog"
)

const MaxVLRConnectAttempts uint = 200

func init() {
	flag.Parse()
}

func main() {
	// Create the service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.CSFB)
	if err != nil {
		glog.Fatalf("Error creating CSFB service: %s", err)
	}

	vlrIP, vlrPort := getVLRSCTPAddr()
	vlrConn, err := servicers.NewSCTPClientConnection(vlrIP, vlrPort)
	if err != nil {
		glog.Fatalf("Failed to create VLR connection: %s", err)
	}

	servicer, err := servicers.NewCsfbServer(vlrConn)
	if err != nil {
		glog.Fatalf("Failed to create CSFB service: %v", err)
	}
	protos.RegisterCSFBFedGWServiceServer(srv.GrpcServer, servicer)

	defer vlrConn.CloseConn() // attempt to close from main thread if GRPC srv errors out

	go func() {
		for retries := uint(0); retries <= MaxVLRConnectAttempts; retries++ {
			err := vlrConn.EstablishConn()
			if err != nil {
				glog.Errorf("Error connecting to VLR Server @ %s:%d; %s; attempt #%d", vlrIP, vlrPort, err, retries)
				time.Sleep(time.Second * time.Duration(retries))
				continue
			}
			var receivedMsg []byte
			for {
				// blocked until a message is received
				receivedMsg, err = vlrConn.Receive()
				if err != nil {
					if err == io.EOF {
						glog.Errorf("Connection to %s:%d is closed by the VLR server", vlrIP, vlrPort)
					} else {
						glog.Errorf("Failed to receive message from %s:%d: %s", vlrIP, vlrPort, err)
					}
					clerr := vlrConn.CloseConn()
					if clerr != nil {
						glog.Errorf("Error closing VLR connection: %s", err)
					}
					break // break out & try to reconnect
				}
				msgType, decodedMsg, err := message.SGsMessageDecoder(receivedMsg)
				if err != nil {
					glog.Errorf("Failed to decode VLR message: %s", err)
					continue
				}
				if msgType == decode.SGsAPResetIndication {
					glog.V(2).Info("Sending Reset Ack to VLR")
					err = servicer.SendResetAck()
					if err != nil {
						glog.Errorf(
							"Failed to send Reset Ack to VLR: %s",
							err,
						)
					}
				}
				_, err = csfb.SendSGsMessageToGateway(msgType, decodedMsg)
				if err != nil {
					glog.Errorf("Failed to send message to gateway: %s", err)
					continue
				}
			}
		}
		glog.Fatalf("Exceeded Maximum VLR Connect Retry Attempts - %d", MaxVLRConnectAttempts)
	}()

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Errorf("Error running service: %s", err)
	}
}

func getVLRSCTPAddr() (string, int) {
	envValue := os.Getenv(servicers.VLRAddrEnv)
	if len(envValue) == 0 {
		glog.Infof("Environment variable for VLR address is not found, "+
			"default address %s:%d is used. ",
			servicers.DefaultVLRIPAddress,
			servicers.DefaultVLRPort,
		)
		return servicers.DefaultVLRIPAddress, servicers.DefaultVLRPort
	}
	addr := strings.Split(envValue, ":")
	if len(addr) != 2 {
		glog.Errorf("VLR address should be in the format: 0.0.0.0:1234, "+
			"default address %s:%d is used. ",
			servicers.DefaultVLRIPAddress,
			servicers.DefaultVLRPort,
		)
		return servicers.DefaultVLRIPAddress, servicers.DefaultVLRPort
	}
	port, err := strconv.Atoi(addr[1])
	if err != nil {
		glog.Infof("Failed to get port number: %s, "+
			"default address %s:%d is used. ",
			err,
			servicers.DefaultVLRIPAddress,
			servicers.DefaultVLRPort,
		)
		return servicers.DefaultVLRIPAddress, servicers.DefaultVLRPort
	}
	glog.Infof("Using %s as the VLR address. ", envValue)
	return addr[0], port
}
