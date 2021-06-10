/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
	serverid          int
)

type cliConfig struct {
	serverCfg    *diameter.DiameterServerConfig
	gyClient     *gy.GyClient
	imsi         string
	sessionID    string
	ueIP         string
	spgwIP       string
	ratingGroups []uint32
	rg2ServiceId map[uint32]uint32
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
	flag.StringVar(&ratingGroupString, "rating_groups", "1,2:200,33", "Rating groups to request credit for (comma separated)")
	flag.Uint64Var(&usedCredit, "used_credit", 10000000, "# of bytes to report as used in CCR-Update and Terminate")
	flag.BoolVar(&wait, "wait", true, "wait for key input in between calls")
	flag.StringVar(&commands, "commands", "", "CCR commands to run keyed by letter, e.g. IT. I = Init, T = Terminate")
	flag.StringVar(&msisdn, "msisdn", "541123525401", "msisdn")
	flag.StringVar(&apn, "apn", "TestMagma", "apn")
	flag.StringVar(&plmn, "plmn", "72207", "PLMN ID")
	flag.IntVar(&serverid, "serverid", 0, "Index of one of the configured servers")

	// Flag help
	allFlags := []string{"help", "imsi", "sid", "rating_groups", "used_credit", "ue_ip", "spgw_ip",
		"commands", "wait", "addr", "network", "host", "realm", "product", "laddr", "dest_host", "dest_realm",
		"msisdn", "apn", "plmn", "serverid"}
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

	// if no serviceid flag, serviceid will be 0. Config on pos 0 will be used
	serverCfg := gy.GetOCSConfiguration()[serverid]
	fmt.Printf("Server config: %+v\n", serverCfg)

	clientCfg := gy.GetGyClientConfiguration()[serverid]
	fmt.Printf("Client config: %+v\n", clientCfg)

	gyGobalCfg := gy.GetGyGlobalConfig()
	fmt.Printf("Gy global config: %+v\n", gyGobalCfg)

	rgs, rg2sid := parseRatingGroups(ratingGroupString)

	config := &cliConfig{
		serverCfg:    serverCfg,
		gyClient:     gy.NewGyClient(clientCfg, serverCfg, handleReAuth, nil, gyGobalCfg),
		imsi:         imsi,
		sessionID:    fmt.Sprintf("%s-%s", imsi, sid),
		ueIP:         ueIP,
		spgwIP:       spgwIP,
		ratingGroups: rgs,
		rg2ServiceId: rg2sid,
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
		var serviceId *uint32
		if sid, ok := config.rg2ServiceId[rg]; ok {
			serviceId = &sid
		}
		if requestType == credit_control.CRTInit {
			credits = append(credits, &gy.UsedCredits{RatingGroup: rg, ServiceIdentifier: serviceId})
		} else {
			credits = append(credits, &gy.UsedCredits{
				RatingGroup:       rg,
				ServiceIdentifier: serviceId,
				InputOctets:       0, // make all used credit output for simplicity
				OutputOctets:      config.usedCredit,
				TotalOctets:       config.usedCredit,
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

func parseRatingGroups(ratingGroupString string) (ratingGroups []uint32, ratingGrpToServiceId map[uint32]uint32) {
	tokens := strings.Split(ratingGroupString, ",")
	ratingGroups = make([]uint32, 0, len(tokens))
	ratingGrpToServiceId = map[uint32]uint32{}
	for _, rgStr := range tokens {
		if i := strings.IndexRune(rgStr, ':'); i > 0 {
			k := strings.TrimSpace(rgStr[0:i])
			v := strings.TrimSpace(rgStr[i+1:])
			if len(k) > 0 && len(v) > 0 {
				rg, err := strconv.ParseUint(k, 10, 32)
				if err == nil {
					si, err := strconv.ParseUint(v, 10, 32)
					if err == nil {
						ratingGroups = append(ratingGroups, uint32(rg))
						ratingGrpToServiceId[uint32(rg)] = uint32(si)
						continue
					}

				}
			}
		}
		rg, _ := strconv.ParseUint(rgStr, 10, 32)
		ratingGroups = append(ratingGroups, uint32(rg))
	}
	return ratingGroups, ratingGrpToServiceId
}
