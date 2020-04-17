/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package integration

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"time"

	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2869"
	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/cwf/gateway/registry"
	"magma/cwf/gateway/services/uesim"
	"magma/feg/gateway/services/eap"
	"magma/lte/cloud/go/crypto"
	lteprotos "magma/lte/cloud/go/protos"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// todo make Op configurable, or export it in the UESimServer.
const (
	Op              = "\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"
	Secret          = "123456"
	MockHSSRemote   = "HSS_REMOTE"
	MockPCRFRemote  = "PCRF_REMOTE"
	MockOCSRemote   = "OCS_REMOTE"
	MockPCRFRemote2 = "PCRF_REMOTE2"
	MockOCSRemote2  = "OCS_REMOTE2"
	PipelinedRemote = "pipelined.local"
	RedisRemote     = "REDIS"
	CwagIP          = "192.168.70.101"
	OCSPort         = 9201
	PCRFPort        = 9202
	OCSPort2        = 9205
	PCRFPort2       = 9206
	HSSPort         = 9204
	PipelinedPort   = 8443
	RedisPort       = 6380

	defaultMSISDN = "5100001234"
)

type TestRunner struct {
	t     *testing.T
	imsis map[string]bool
}

// imsi -> ruleID -> record
type RecordByIMSI map[string]map[string]*lteprotos.RuleRecord

// NewTestRunner initializes a new TestRunner by making a UESim client and
// and setting the next IMSI.
func NewTestRunner(t *testing.T) *TestRunner {
	fmt.Println("************************* TestRunner setup")
	testRunner := &TestRunner{t: t}

	testRunner.imsis = make(map[string]bool)
	fmt.Printf("Adding Mock HSS service at %s:%d\n", CwagIP, HSSPort)
	registry.AddService(MockHSSRemote, CwagIP, HSSPort)
	fmt.Printf("Adding Mock PCRF service at %s:%d\n", CwagIP, PCRFPort)
	registry.AddService(MockPCRFRemote, CwagIP, PCRFPort)
	fmt.Printf("Adding Mock OCS service at %s:%d\n", CwagIP, OCSPort)
	registry.AddService(MockOCSRemote, CwagIP, OCSPort)
	fmt.Printf("Adding Pipelined service at %s:%d\n", CwagIP, PipelinedPort)
	registry.AddService(PipelinedRemote, CwagIP, PipelinedPort)
	fmt.Printf("Adding Redis service at %s:%d\n", CwagIP, RedisPort)
	registry.AddService(RedisRemote, CwagIP, RedisPort)

	return testRunner
}

// NewTestRunnerWithTwoPCRFandOCS does the same as NewTestRunner but it inclides 2 PCRF and 2 OCS
// Used in scenarios that run 2 PCRFs and 2 OCSs
func NewTestRunnerWithTwoPCRFandOCS(t *testing.T) *TestRunner {
	tr :=NewTestRunner(t)
	fmt.Printf("Adding second Mock PCRF service at %s:%d\n", CwagIP, PCRFPort2)
	registry.AddService(MockPCRFRemote2, CwagIP, PCRFPort2)
	fmt.Printf("Adding second OCS service at %s:%d\n", CwagIP, OCSPort2)
	registry.AddService(MockOCSRemote2, CwagIP, OCSPort2)
	return tr
}

// ConfigUEs creates and adds the specified number of UEs and Subscribers
// to the UE Simulator and the HSS.
func (tr *TestRunner) ConfigUEs(numUEs int) ([]*cwfprotos.UEConfig, error) {
	fmt.Printf("************************* Configuring %d UE(s)\n", numUEs)
	ues := make([]*cwfprotos.UEConfig, 0)
	for i := 0; i < numUEs; i++ {
		imsi := ""
		for {
			imsi = getRandomIMSI()
			_, present := tr.imsis[imsi]
			if !present {
				break
			}
		}
		key, opc, err := getRandKeyOpcFromOp([]byte(Op))
		if err != nil {
			return nil, err
		}
		seq := getRandSeq()

		ue := makeUE(imsi, key, opc, seq)
		sub := makeSubscriber(imsi, key, opc, seq+1)

		err = uesim.AddUE(ue)
		if err != nil {
			return nil, errors.Wrap(err, "Error adding UE to UESimServer")
		}
		err = addSubscriberToHSS(sub)
		if err != nil {
			return nil, errors.Wrap(err, "Error adding Subscriber to HSS")
		}
		err = addSubscriberToPCRF(sub.GetSid())
		if err != nil {
			return nil, errors.Wrap(err, "Error adding Subscriber to PCRF")
		}
		err = addSubscriberToOCS(sub.GetSid())
		if err != nil {
			return nil, errors.Wrap(err, "Error adding Subscriber to OCS")
		}

		ues = append(ues, ue)
		fmt.Printf("Added UE to Simulator, HSS, PCRF, and OCS:\n"+
			"\tIMSI: %s\tKey: %x\tOpc: %x\tSeq: %d\n", imsi, key, opc, seq)
		tr.imsis[imsi] = true
	}
	fmt.Println("Successfully configured UE(s)")
	return ues, nil
}

// Authenticate simulates an authentication between the UE with the specified
// IMSI and the HSS, and returns the resulting Radius packet.
func (tr *TestRunner) Authenticate(imsi string) (*radius.Packet, error) {
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
	tr.t.Logf("Finished Authenticating UE. Resulting RADIUS Packet: %d\n", radiusP)
	return radiusP, nil
}

func (tr *TestRunner) AuthenticateAndAssertSuccess(imsi string) {
	radiusP, err := tr.Authenticate(imsi)
	assert.NoError(tr.t, err)

	eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
	assert.NotNil(tr.t, eapMessage, fmt.Sprintf("EAP Message from authentication is nil"))
	assert.True(tr.t, reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode), fmt.Sprintf("UE Authentication did not return success"))
}

func (tr *TestRunner) AuthenticateAndAssertFail(imsi string) {
	radiusP, err := tr.Authenticate(imsi)
	assert.NoError(tr.t, err)

	eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
	assert.NotNil(tr.t, eapMessage)
	assert.True(tr.t, reflect.DeepEqual(int(eapMessage[0]), eap.FailureCode))
}

// Authenticate simulates an authentication between the UE with the specified
// IMSI and the HSS, and returns the resulting Radius packet.
func (tr *TestRunner) Disconnect(imsi string) (*radius.Packet, error) {
	fmt.Printf("************************* Sending a disconnect request UE with IMSI: %s\n", imsi)
	res, err := uesim.Disconnect(&cwfprotos.DisconnectRequest{Imsi: imsi})
	if err != nil {
		return &radius.Packet{}, err
	}
	encoded := res.GetRadiusPacket()
	radiusP, err := radius.Parse(encoded, []byte(Secret))
	if err != nil {
		err = errors.Wrap(err, "Error while parsing encoded Radius packet")
		fmt.Println(err)
		return &radius.Packet{}, err
	}
	tr.t.Logf("Finished Discconnecting UE. Resulting RADIUS Packet: %d\n", radiusP)
	return radiusP, nil
}

// GenULTraffic simulates the UE sending traffic through the CWAG to the Internet
// by running an iperf3 client on the UE simulator and an iperf3 server on the
// Magma traffic server.
func (tr *TestRunner) GenULTraffic(req *cwfprotos.GenTrafficRequest) (*cwfprotos.GenTrafficResponse, error) {
	fmt.Printf("************************* Generating Traffic for UE with Req: %v\n", req)
	return uesim.GenTraffic(req)
}

// Remove subscribers, rules, flows, and monitors to clean up the state for
// consecutive test runs
func (tr *TestRunner) CleanUp() error {
	for imsi, _ := range tr.imsis {
		err := deleteSubscribersFromHSS(imsi)
		if err != nil {
			return err
		}
	}
	err := clearSubscribersFromPCRF()
	if err != nil {
		return err
	}
	err = clearSubscribersFromOCS()
	if err != nil {
		return err
	}

	return nil
}

// GetPolicyUsage is a wrapper around pipelined's GetPolicyUsage and returns
// the policy usage keyed by subscriber ID
func (tr *TestRunner) GetPolicyUsage() (RecordByIMSI, error) {
	recordsBySubID := RecordByIMSI{}
	table, err := getPolicyUsage()
	if err != nil {
		return recordsBySubID, err
	}
	for _, record := range table.Records {
		fmt.Printf("Record %v\n", record)
		_, exists := recordsBySubID[record.Sid]
		if !exists {
			recordsBySubID[record.Sid] = map[string]*lteprotos.RuleRecord{}
		}
		recordsBySubID[record.Sid][record.RuleId] = record
	}
	return recordsBySubID, nil
}

func (tr *TestRunner) WaitForEnforcementStatsToSync() {
	// TODO load this value from pipelined.yml
	enforcementPollPeriod := 1 * time.Second
	time.Sleep(3 * enforcementPollPeriod)
}

func (tr *TestRunner) WaitForPoliciesToSync() {
	// TODO load this value from sessiond.yml (rule_update_interval_sec)
	ruleUpdatePeriod := 1 * time.Second
	time.Sleep(2 * ruleUpdatePeriod)
}

func (tr *TestRunner) WaitForReAuthToProcess() {
	// Todo figure out the best way to figure out when RAR is processed
	time.Sleep(3 * time.Second)
}

// getRandomIMSI makes a random 15-digit IMSI that is not added to the UESim or HSS.
func getRandomIMSI() string {
	imsi := ""
	for len(imsi) < 15 {
		imsi += strconv.Itoa(rand.Intn(10))
	}
	return imsi
}

// RandKeyOpc makes a random 16-byte key and calculates the Opc based off the Op.
func getRandKeyOpcFromOp(op []byte) (key, opc []byte, err error) {
	key = make([]byte, 16)
	rand.Read(key)

	tempOpc, err := crypto.GenerateOpc(key, op)
	if err != nil {
		return nil, nil, err
	}
	opc = tempOpc[:]
	return
}

// getRandSeq makes a random 43-bit Seq.
func getRandSeq() uint64 {
	return rand.Uint64() >> 21
}

// makeUE creates a new UE using the given values.
func makeUE(imsi string, key []byte, opc []byte, seq uint64) *cwfprotos.UEConfig {
	return &cwfprotos.UEConfig{
		Imsi:    imsi,
		AuthKey: key,
		AuthOpc: opc,
		Seq:     seq,
	}
}

func prependIMSIPrefix(imsi string) string {
	return "IMSI" + imsi
}

// MakeSubcriber creates a new Subscriber using the given values.
func makeSubscriber(imsi string, key []byte, opc []byte, seq uint64) *lteprotos.SubscriberData {
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
			ApnConfig:           []*lteprotos.APNConfiguration{&lteprotos.APNConfiguration{}},
		},
	}
}
