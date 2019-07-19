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
	"magma/cwf/gateway/services/uesim"
	"magma/feg/gateway/services/testcore/hss"
	"magma/lte/cloud/go/crypto"
	lteprotos "magma/lte/cloud/go/protos"

	"github.com/pkg/errors"
)

// todo make Op configurable, or export it in the UESimServer.
const (
	Op     = "\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"
	Secret = "123456"
)

type TestRunner struct {
	Imsis map[string]bool
}

// NewTestRunner initializes a new TestRunner by making a UESim client and
// and setting the next IMSI.
func NewTestRunner() *TestRunner {
	fmt.Println("************************* TestRunner setup")
	testRunner := &TestRunner{}

	testRunner.Imsis = make(map[string]bool)

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
		sub := MakeSubscriber(imsi, key, opc, seq)

		err = uesim.AddUE(ue)
		if err != nil {
			return nil, errors.Wrap(err, "Error adding UE to UESimServer")
		}

		err = hss.AddSubscriber(sub)
		if err != nil {
			return nil, errors.Wrap(err, "Error adding Subscriber to HSS")
		}

		ues = append(ues, ue)
		fmt.Printf("Added UE to Simulator and HSS:\n"+
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
	}
}
