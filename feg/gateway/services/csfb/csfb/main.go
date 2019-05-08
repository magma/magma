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
	"github.com/ishidawataru/sctp"
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

	vlrSCTPAddr := getVLRSCTPAddr()
	localSCTPAddr := getSGsInterfaceAddr()
	vlrConn, err := servicers.NewSCTPClientConnection(vlrSCTPAddr, localSCTPAddr)
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
			vlrIP := vlrSCTPAddr.IPAddrs[0].String()
			vlrPort := vlrSCTPAddr.Port
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

func getAddr(address string, defaultIP string, defaultPort int) (string, int) {
	if len(address) == 0 {
		glog.V(2).Infof("Environment variable for address is not found or empty.")
		return defaultIP, defaultPort
	}
	addr := strings.Split(address, ":")
	if len(addr) != 2 {
		glog.Errorf("Address should be in the format: 0.0.0.0:1234.")
		return defaultIP, defaultPort
	}
	port, err := strconv.Atoi(addr[1])
	if err != nil {
		glog.Errorf("Failed to get port number: %s.", err)
		return defaultIP, defaultPort
	}
	return addr[0], port
}

func getSGsInterfaceAddr() *sctp.SCTPAddr {
	localAddr := os.Getenv(servicers.LocalAddrEnv)
	glog.V(2).Info("Getting local SGs interface adddress.")
	ip, port := getAddr(localAddr, "", 0)
	if len(ip) == 0 {
		glog.V(2).Infof("The local SGs interface address is not specified.")
	}
	glog.V(2).Infof("Using %s:%d as the local SGs interface address. ", ip, port)
	return servicers.ConstructSCTPAddr(ip, port)
}

func getVLRSCTPAddr() *sctp.SCTPAddr {
	vlrAddr := os.Getenv(servicers.VLRAddrEnv)
	glog.V(2).Info("Getting VLR adddress.")
	ip, port := getAddr(vlrAddr, servicers.DefaultVLRIPAddress, servicers.DefaultVLRPort)
	glog.V(2).Infof("Using %s:%d as the VLR address. ", ip, port)
	return servicers.ConstructSCTPAddr(ip, port)
}
