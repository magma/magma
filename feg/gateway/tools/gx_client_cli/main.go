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

	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/feg/gateway/services/session_proxy/credit_control/gx"
)

const (
	// values obtained from example Gx traces
	defaultTAI  uint16 = 2461
	defaultECGI uint32 = 134477315
)

var (
	imsi     string
	sid      string
	ueIP     string
	commands string
	wait     bool
	help     bool
	msisdn   string
	apn      string
	plmn     string
	serverid int
)

type cliConfig struct {
	serverCfg *diameter.DiameterServerConfig
	gxClient  *gx.GxClient
	imsi      string
	sessionID string
	ueIP      string
	msisdn    string
	apn       string
	plmn      string
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
	flag.BoolVar(&wait, "wait", true, "wait for key input in between calls")
	flag.StringVar(&commands, "commands", "", "CCR commands to run keyed by letter, e.g. IT. I = Init, T = Terminate")
	flag.StringVar(&msisdn, "msisdn", "541123525401", "msisdn")
	flag.StringVar(&apn, "apn", "TestMagma", "apn")
	flag.StringVar(&plmn, "plmn", "72207", "PLMN ID")
	flag.IntVar(&serverid, "serverid", 0, "Index of one of the configured servers")

	// Flag help
	allFlags := []string{"help", "imsi", "sid", "ue_ip", "commands", "wait", "addr", "network", "host",
		"realm", "product", "laddr", "dest_host", "dest_realm", "msisdn", "apn", "plmn", "serverid"}
	flag.Usage = func() {
		fmt.Println("Gx Client CLI for testing Gx Diameter CCR calls.")
		fmt.Println("Usage:\n	gx_client_cli")
		fmt.Println("Flags: ")
		for _, flagName := range allFlags {
			fmt.Printf("	%s: %s\n", flagName, flag.Lookup(flagName).Usage)
		}
	}
}

// gx_client_cli is a CLI for testing CCR requests against a PCRF. The PCRF settings can be
// either loaded from preexisting environment variables or specified through command line
// flags.
// Example usage:
//   gx_client_cli --imsi=001010000000001 --sid="1234" --commands="I"
func main() {
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	// if no serviceid flag, serviceid will be 0. Config on pos 0 will be used
	serverCfg := gx.GetPCRFConfiguration()[serverid]
	fmt.Printf("Server config: %+v\n", serverCfg)

	clientCfg := gx.GetGxClientConfiguration()[serverid]
	fmt.Printf("Client config: %+v\n", clientCfg)

	globalCfg := gx.GetGxGlobalConfig()

	config := &cliConfig{
		serverCfg: serverCfg,
		gxClient:  gx.NewGxClient(clientCfg, serverCfg, handleReAuth, nil, globalCfg),
		imsi:      imsi,
		sessionID: fmt.Sprintf("%s-%s", imsi, sid),
		ueIP:      ueIP,
		msisdn:    msisdn,
		apn:       apn,
		plmn:      plmn,
	}

	if len(commands) == 0 {
		fmt.Printf("No commands specified. To use this CLI, specify CCR requests to send using the `commands` flag.\n")
		fmt.Printf("For example, specifying `--commands=IT` runs (I)nit, (T)erminate\n")
		os.Exit(1)
	}

	for i, cmd := range commands {
		requestNum := uint32(i)
		switch cmd {
		case 'I', 'i':
			sendInitCall(config, requestNum)
		case 'T', 't':
			sendTerminateCall(config, requestNum)
		}
		if i < len(commands)-1 && wait {
			fmt.Print("Press 'Enter' to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
	}
}

func sendInitCall(config *cliConfig, requestNumber uint32) {
	userLoc, err := diameter.EncodeUserLocation(config.plmn, defaultTAI, defaultECGI)
	if err != nil {
		fmt.Printf("Error encoding user location: %s", err)
		os.Exit(1)
	}
	ccrInit := &gx.CreditControlRequest{
		SessionID:     config.sessionID,
		Type:          credit_control.CRTInit,
		IMSI:          config.imsi,
		RequestNumber: requestNumber,
		IPAddr:        config.ueIP,
		Msisdn:        []byte(config.msisdn),
		Apn:           config.apn,
		PlmnID:        config.plmn,
		UserLocation:  userLoc,
	}

	fmt.Printf("Sending CCR-Init: %+v\n", ccrInit)
	done := make(chan interface{}, 1000)

	config.gxClient.SendCreditControlRequest(config.serverCfg, done, ccrInit)
	answer := gx.GetAnswer(done)
	fmt.Printf("CCA-Init: %+v\n", answer)
	for _, rule := range answer.RuleInstallAVP {
		for _, ruleID := range rule.RuleNames {
			fmt.Printf(" -- Static Rule: %s\n", ruleID)
		}
		for _, baseName := range rule.RuleBaseNames {
			fmt.Printf(" -- Base Name: %s\n", baseName)
		}
		for _, def := range rule.RuleDefinitions {
			fmt.Printf(" -- Rule Definition: %+v\n", def)
		}
	}
}

func sendTerminateCall(config *cliConfig, requestNumber uint32) {
	ccrTerminate := &gx.CreditControlRequest{
		SessionID:     config.sessionID,
		Type:          credit_control.CRTTerminate,
		IMSI:          config.imsi,
		RequestNumber: requestNumber,
		IPAddr:        config.ueIP,
	}
	done := make(chan interface{}, 1000)
	fmt.Printf("Sending CCR-Terminate: %+v\n", ccrTerminate)
	config.gxClient.SendCreditControlRequest(config.serverCfg, done, ccrTerminate)
	answer := gx.GetAnswer(done)
	fmt.Printf("CCA-Terminate: %+v\n", answer)
}

func handleReAuth(request *gx.PolicyReAuthRequest) *gx.PolicyReAuthAnswer {
	return nil
}
