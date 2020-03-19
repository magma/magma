/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/feg/gateway/services/session_proxy/credit_control/gy"
)

const (
	// values obtained from example Gx traces
	defaultTAI  uint16 = 2461
	defaultECGI uint32 = 134477315
)

var (
	imsi              string
	sid               string
	ueIP              string
	spgwIP            string
	ratingGroupString string
	usedCredit        uint64
	commands          string
	wait              bool
	help              bool
	msisdn            string
	apn               string
	plmn              string
)

type cliConfig struct {
	serverCfg    *diameter.DiameterServerConfig
	gyClient     *gy.GyClient
	imsi         string
	sessionID    string
	ueIP         string
	spgwIP       string
	ratingGroups []uint32
	usedCredit   uint64
	msisdn       string
	apn          string
	plmn         string
}

func init() {
	// Enable logging
	flag.Set("v", "10")             // enable the most verbose logging
	flag.Set("logtostderr", "true") // enable printing to console

	// New Flags
	flag.BoolVar(&help, "help", false, "[optional] Display this help message")
	flag.StringVar(&imsi, "imsi", "001010000000031", "imsi")
	flag.StringVar(&sid, "sid", "1234", "session id")
	flag.StringVar(&ueIP, "ue_ip", "192.168.1.1", "UE IPv4 address")
	flag.StringVar(&spgwIP, "spgw_ip", "192.168.128.1", "SPGW IPv4 address")
	flag.StringVar(&ratingGroupString, "rating_groups", "1,2,33", "Rating groups to request credit for (comma separated)")
	flag.Uint64Var(&usedCredit, "used_credit", 10000000, "# of bytes to report as used in CCR-Update and Terminate")
	flag.BoolVar(&wait, "wait", true, "wait for key input in between calls")
	flag.StringVar(&commands, "commands", "", "CCR commands to run keyed by letter, e.g. IT. I = Init, T = Terminate")
	flag.StringVar(&msisdn, "msisdn", "541123525401", "msisdn")
	flag.StringVar(&apn, "apn", "TestMagma", "apn")
	flag.StringVar(&plmn, "plmn", "72207", "PLMN ID")

	// Flag help
	allFlags := []string{"help", "imsi", "sid", "rating_groups", "used_credit", "ue_ip", "spgw_ip",
		"commands", "wait", "addr", "network", "host", "realm", "product", "laddr", "dest_host", "dest_realm",
		"msisdn", "apn", "plmn"}
	flag.Usage = func() {
		fmt.Println("Gx Client CLI for testing Gx Diameter CCR calls.")
		fmt.Println("Usage:\n	gx_client_cli")
		fmt.Println("Flags: ")
		for _, flagName := range allFlags {
			fmt.Printf("	%s: %s\n", flagName, flag.Lookup(flagName).Usage)
		}
	}
}

// gy_client_cli is a CLI for testing CCR requests against an OCS. The OCS settings can be
// either loaded from preexisting environment variables or specified through command line
// flags.
// Example usage:
//   gy_client_cli --imsi=001010000000001 --sid="1234" --rating_groups="1,4,2" --commands="IUT"
func main() {
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	serverCfg := gy.GetOCSConfiguration()
	fmt.Printf("Server config: %+v\n", serverCfg)

	clientCfg := gy.GetGyClientConfiguration()
	fmt.Printf("Client config: %+v\n", clientCfg)

	gyGobalCfg := gy.GetGyGlobalConfig()
	fmt.Printf("Gy global config: %+v\n", gyGobalCfg)

	config := &cliConfig{
		serverCfg:    serverCfg,
		gyClient:     gy.NewGyClient(clientCfg, serverCfg, handleReAuth, nil, gyGobalCfg),
		imsi:         imsi,
		sessionID:    fmt.Sprintf("%s-%s", imsi, sid),
		ueIP:         ueIP,
		spgwIP:       spgwIP,
		ratingGroups: parseRatingGroups(ratingGroupString),
		usedCredit:   usedCredit,
		msisdn:       msisdn,
		apn:          apn,
		plmn:         plmn,
	}

	if len(commands) == 0 {
		fmt.Printf("No commands specified. To use this CLI, specify CCR requests to send using the `commands` flag.\n")
		fmt.Printf("For example, specifying `--commands=IUT` runs (I)nit, (U)pdate, (T)erminate\n")
		os.Exit(1)
	}

	for i, cmd := range commands {
		requestNum := uint32(i)
		switch cmd {
		case 'I', 'i':
			sendCreditCall(config, credit_control.CRTInit, requestNum)
		case 'U', 'u':
			sendCreditCall(config, credit_control.CRTUpdate, requestNum)
		case 'T', 't':
			sendCreditCall(config, credit_control.CRTTerminate, requestNum)
		}
		if i < len(commands)-1 && wait {
			fmt.Print("Press 'Enter' to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
	}
}

func sendCreditCall(config *cliConfig, requestType credit_control.CreditRequestType, requestNumber uint32) {
	userLoc, err := diameter.EncodeUserLocation(config.plmn, defaultTAI, defaultECGI)
	if err != nil {
		fmt.Printf("Error encoding user location: %s", err)
		os.Exit(1)
	}

	credits := make([]*gy.UsedCredits, 0, len(config.ratingGroups))
	for _, rg := range config.ratingGroups {
		if requestType == credit_control.CRTInit {
			credits = append(credits, &gy.UsedCredits{RatingGroup: rg})
		} else {
			credits = append(credits, &gy.UsedCredits{
				RatingGroup:  rg,
				InputOctets:  0, // make all used credit output for simplicity
				OutputOctets: config.usedCredit,
				TotalOctets:  config.usedCredit,
			})
		}
	}
	imei := []byte{0x53, 0x45, 0x28, 0x60, 0x45, 0x92, 0x75, 0x11} // bogus value
	ccr := &gy.CreditControlRequest{
		SessionID:     config.sessionID,
		Type:          requestType,
		IMSI:          config.imsi,
		RequestNumber: requestNumber,
		UeIPV4:        config.ueIP,
		SpgwIPV4:      config.spgwIP,
		Credits:       credits,
		Msisdn:        []byte(config.msisdn),
		Apn:           config.apn,
		Imei:          string(imei),
		PlmnID:        config.plmn,
		UserLocation:  userLoc,
	}
	fmt.Printf("Sending CCR: %+v\n", ccr)
	done := make(chan interface{}, 1000)

	config.gyClient.SendCreditControlRequest(config.serverCfg, done, ccr)
	answer := gy.GetAnswer(done)
	fmt.Printf("CCA: %+v\n", answer)
	if requestType != credit_control.CRTTerminate {
		for _, credit := range answer.Credits {
			fmt.Printf("-- Received Credit %d: %+v\n", credit.RatingGroup, credit)
		}
	}
}

func handleReAuth(request *gy.ChargingReAuthRequest) *gy.ChargingReAuthAnswer {
	return nil
}

func parseRatingGroups(ratingGroupString string) []uint32 {
	tokens := strings.Split(ratingGroupString, ",")
	ratingGroups := make([]uint32, 0, len(tokens))
	for _, rgStr := range tokens {
		rg, _ := strconv.ParseUint(rgStr, 10, 32)
		ratingGroups = append(ratingGroups, uint32(rg))
	}
	return ratingGroups
}
