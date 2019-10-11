/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package integ_tests

import (
	"fmt"
	"math/rand"
	"strconv"

	"fbc/lib/go/radius"
	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/cwf/gateway/registry"
	"magma/cwf/gateway/services/uesim"
	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/testcore/hss"
	"magma/lte/cloud/go/crypto"
	lteprotos "magma/lte/cloud/go/protos"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// todo make Op configurable, or export it in the UESimServer.
const (
	Op             = "\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"
	Secret         = "123456"
	MockHSSRemote  = "HSS_REMOTE"
	MockPCRFRemote = "PCRF_REMOTE"
	MockOCSRemote  = "OCS_REMOTE"
	CwagIP         = "192.168.70.101"
	HSSPort        = 9204
	OCSPort        = 9201
	PCRFPort       = 9202

	defaultMSISDN = "5100001234"
)

type TestRunner struct {
	Imsis map[string]bool
}

// Wrapper for GRPC Client functionality
type hssClient struct {
	fegprotos.HSSConfiguratorClient
	cc *grpc.ClientConn
}

// Wrapper for GRPC Client functionality
type pcrfClient struct {
	fegprotos.MockPCRFClient
	cc *grpc.ClientConn
}

// Wrapper for GRPC Client functionality
type ocsClient struct {
	fegprotos.MockOCSClient
	cc *grpc.ClientConn
}

// getHSSClient is a utility function to getHSSClient a RPC connection to a
// remote HSS service.
func getHSSClient() (*hssClient, error) {
	var conn *grpc.ClientConn
	var err error
	conn, err = registry.GetConnection(MockHSSRemote)
	if err != nil {
		errMsg := fmt.Sprintf("HSS client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &hssClient{
		fegprotos.NewHSSConfiguratorClient(conn),
		conn,
	}, err
}

// getPCRFClient is a utility function to get an RPC connection to a
// remote PCRF service.
func getPCRFClient() (*pcrfClient, error) {
	var conn *grpc.ClientConn
	var err error
	conn, err = registry.GetConnection(MockPCRFRemote)
	if err != nil {
		errMsg := fmt.Sprintf("PCRF client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &pcrfClient{
		fegprotos.NewMockPCRFClient(conn),
		conn,
	}, err
}

// getOCSClient is a utility function to an RPC connection to a
// remote OCS service.
func getOCSClient() (*ocsClient, error) {
	var conn *grpc.ClientConn
	var err error
	conn, err = registry.GetConnection(MockOCSRemote)
	if err != nil {
		errMsg := fmt.Sprintf("PCRF client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &ocsClient{
		fegprotos.NewMockOCSClient(conn),
		conn,
	}, err
}

// addSubscriber tries to add this subscriber to the HSS server.
// This function returns an AlreadyExists error if the subscriber has already
// been added.
// Input: The subscriber data which will be added.
func addSubscriberToHss(sub *lteprotos.SubscriberData) error {
	err := hss.VerifySubscriberData(sub)
	if err != nil {
		errMsg := fmt.Errorf("Invalid AddSubscriberRequest provided: %s", err)
		return errors.New(errMsg.Error())
	}
	cli, err := getHSSClient()
	if err != nil {
		return err
	}
	_, err = cli.AddSubscriber(context.Background(), sub)
	return err
}

// addSubscriber tries to add this subscriber to the PCRF server.
// Input: The subscriber data which will be added.
func addSubscriberToPcrf(sub *lteprotos.SubscriberID) error {
	cli, err := getPCRFClient()
	if err != nil {
		return err
	}
	_, err = cli.CreateAccount(context.Background(), sub)
	return err
}

// addSubscriber tries to add this subscriber to the OCS server.
// Input: The subscriber data which will be added.
func addSubscriberToOcs(sub *lteprotos.SubscriberID) error {
	cli, err := getOCSClient()
	if err != nil {
		return err
	}
	_, err = cli.CreateAccount(context.Background(), sub)
	return err
}

// NewTestRunner initializes a new TestRunner by making a UESim client and
// and setting the next IMSI.
func NewTestRunner() *TestRunner {
	fmt.Println("************************* TestRunner setup")
	testRunner := &TestRunner{}

	testRunner.Imsis = make(map[string]bool)
	fmt.Printf("Adding Mock HSS service at %s:%d\n", CwagIP, HSSPort)
	registry.AddService(MockHSSRemote, CwagIP, HSSPort)
	fmt.Printf("Adding Mock PCRF service at %s:%d\n", CwagIP, PCRFPort)
	registry.AddService(MockPCRFRemote, CwagIP, PCRFPort)
	fmt.Printf("Adding Mock OCS service at %s:%d\n", CwagIP, OCSPort)
	registry.AddService(MockOCSRemote, CwagIP, OCSPort)

	return testRunner
}

// ConfigUEs creates and adds the specified number of UEs and Subscribers
// to the UE Simulator and the HSS.
func (testRunner *TestRunner) ConfigUEs(numUEs int) ([]*cwfprotos.UEConfig, error) {
	fmt.Printf("************************* Configuring %d UE(s)\n", numUEs)
	ues := make([]*cwfprotos.UEConfig, 0)
	for i := 0; i < numUEs; i++ {
		imsi := ""
		for {
			imsi = RandImsi()
			_, present := testRunner.Imsis[imsi]
			if !present {
				break
			}
		}
		key, opc, err := RandKeyOpcFromOp([]byte(Op))
		if err != nil {
			return nil, err
		}
		seq := RandSeq()

		ue := MakeUE(imsi, key, opc, seq)
		sub := MakeSubscriber(imsi, key, opc, seq+1)

		err = uesim.AddUE(ue)
		if err != nil {
			return nil, errors.Wrap(err, "Error adding UE to UESimServer")
		}
		err = addSubscriberToHss(sub)
		if err != nil {
			return nil, errors.Wrap(err, "Error adding Subscriber to HSS")
		}
		err = addSubscriberToPcrf(sub.GetSid())
		if err != nil {
			return nil, errors.Wrap(err, "Error adding Subscriber to PCRF")
		}
		err = addSubscriberToOcs(sub.GetSid())
		if err != nil {
			return nil, errors.Wrap(err, "Error adding Subscriber to OCS")
		}

		ues = append(ues, ue)
		fmt.Printf("Added UE to Simulator, HSS, PCRF, and OCS:\n"+
			"\tIMSI: %s\tKey: %x\tOpc: %x\tSeq: %d\n", imsi, key, opc, seq)
		testRunner.Imsis[imsi] = true
	}
	fmt.Println("Successfully configured UE(s)")
	return ues, nil
}

// Authenticate simulates an authentication between the UE with the specified
// IMSI and the HSS, and returns the resulting Radius packet.
func (testRunner *TestRunner) Authenticate(imsi string) (*radius.Packet, error) {
	fmt.Printf("************************* Authenticating UE with IMSI: %s\n", imsi)
	res, err := uesim.Authenticate(&cwfprotos.AuthenticateRequest{Imsi: imsi})
	if err != nil {
		fmt.Println(err)
		return &radius.Packet{}, err
	}
	encoded := res.GetRadiusPacket()
	radiusP, err := radius.Parse(encoded, []byte(Secret))
	if err != nil {
		err = errors.Wrap(err, "Error while parsing encoded Radius packet")
		fmt.Println(err)
		return &radius.Packet{}, err
	}
	fmt.Printf("Finished Authenticating UE. Resulting RADIUS Packet: %d\n", radiusP)
	return radiusP, nil
}

// GenULTraffic simulates the UE sending traffic through the CWAG to the Internet
// by running an iperf3 client on the UE simulator and an iperf3 server on the
// Magma traffic server.
func (testRunner *TestRunner) GenULTraffic(imsi string) error {
	fmt.Printf("************************* Generating Traffic for UE with IMSI: %s\n", imsi)
	req := &cwfprotos.GenTrafficRequest{
		Imsi: imsi,
	}
	return uesim.GenTraffic(req)
}

// RandImsi makes a random 15-digit IMSI that is not added to the UESim or HSS.
func RandImsi() string {
	imsi := ""
	for len(imsi) < 15 {
		imsi += strconv.Itoa(rand.Intn(10))
	}
	return imsi
}

// RandKeyOpc makes a random 16-byte key and calculates the Opc based off the Op.
func RandKeyOpcFromOp(op []byte) (key, opc []byte, err error) {
	key = make([]byte, 16)
	rand.Read(key)

	tempOpc, err := crypto.GenerateOpc(key, op)
	if err != nil {
		return nil, nil, err
	}
	opc = tempOpc[:]
	return
}

// RandSeq makes a random 43-bit Seq.
func RandSeq() uint64 {
	return rand.Uint64() >> 21
}

// MakeUE creates a new UE using the given values.
func MakeUE(imsi string, key []byte, opc []byte, seq uint64) *cwfprotos.UEConfig {
	return &cwfprotos.UEConfig{
		Imsi:    imsi,
		AuthKey: key,
		AuthOpc: opc,
		Seq:     seq,
	}
}

// MakeSubcriber creates a new Subscriber using the given values.
func MakeSubscriber(imsi string, key []byte, opc []byte, seq uint64) *lteprotos.SubscriberData {
	return &lteprotos.SubscriberData{
		Sid: &lteprotos.SubscriberID{
			Id:   imsi,
			Type: 1,
		},
		Lte: &lteprotos.LTESubscription{
			State:    1,
			AuthAlgo: 0,
			AuthKey:  key,
			AuthOpc:  opc,
		},
		State: &lteprotos.SubscriberState{
			LteAuthNextSeq: seq,
		},
		Non_3Gpp: &lteprotos.Non3GPPUserProfile{
			Msisdn:              defaultMSISDN,
			Non_3GppIpAccess:    lteprotos.Non3GPPUserProfile_NON_3GPP_SUBSCRIPTION_ALLOWED,
			Non_3GppIpAccessApn: lteprotos.Non3GPPUserProfile_NON_3GPP_APNS_ENABLE,
			ApnConfig:           &lteprotos.APNConfiguration{},
		},
	}
}
