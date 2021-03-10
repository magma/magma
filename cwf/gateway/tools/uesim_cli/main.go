/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2869"
	"magma/cwf/cloud/go/protos"
	"magma/cwf/gateway/registry"
	"magma/cwf/gateway/services/uesim"
	"magma/feg/gateway/services/eap"
	"magma/lte/cloud/go/crypto"
	"magma/orc8r/cloud/go/tools/commands"
	"magma/orc8r/lib/go/service/config"

	"github.com/golang/glog"
)

const (
	DefaultOp = "\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"

	DefaultMaxBreakSecs         = 5
	DefaultTrafficGenLengthSecs = 60
	DefaultRadiusSecret         = "123456"
	DefaultApn                  = "test"
	DefaultMsisdn               = "5100001234"
	DefaultRatType              = 6
)

var (
	cmdRegistry          = new(commands.Map)
	trafficLength uint64 = DefaultTrafficGenLengthSecs
	maxBreak      uint64 = DefaultMaxBreakSecs
)

func init() {
	// Enable logging
	flag.Set("v", "10")             // enable the most verbose logging
	flag.Set("logtostderr", "true") // enable printing to console
	authCmd := cmdRegistry.Add("auth", "Send Authenticate Request to UE Simulator", handleAuthCmd)
	authFlags := authCmd.Flags()
	authFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, // std Usage() & PrintDefaults() use Stderr
			"\tUsage: %s [OPTIONS] %s [%s OPTIONS] <IMSI>\n", os.Args[0], authCmd.Name(), authCmd.Name())
		authFlags.PrintDefaults()
	}

	disconnectCmd := cmdRegistry.Add("disconnect", "Send Disconnect Request to UE Simulator", handleDisconnectCmd)
	disconnectFlags := disconnectCmd.Flags()
	disconnectFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, // std Usage() & PrintDefaults() use Stderr
			"\tUsage: %s [OPTIONS] %s [%s OPTIONS] <IMSI>\n", os.Args[0], disconnectCmd.Name(), disconnectCmd.Name())
		authFlags.PrintDefaults()
	}

	addUeCmd := cmdRegistry.Add("add_ue", "Add UE to UE Simulator", handleAddUeCmd)
	addUeFlags := addUeCmd.Flags()
	addUeFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, // std Usage() & PrintDefaults() use Stderr
			"\tUsage: %s [OPTIONS] %s [%s OPTIONS] <IMSI>\n", os.Args[0], addUeCmd.Name(), addUeCmd.Name())
		authFlags.PrintDefaults()
	}

	trafficGenCmd := cmdRegistry.Add("gen_traffic", "Generate control load by continously triggering authentication", handleTrafficGenCmd)
	trafficGenFlags := trafficGenCmd.Flags()
	trafficGenFlags.Uint64Var(&trafficLength, "length", DefaultTrafficGenLengthSecs, "Amount of time (in seconds) to run traffic generation")
	trafficGenFlags.Uint64Var(&maxBreak, "max_break", DefaultMaxBreakSecs, "Max amount of time between auth requests for each subscriber")
	trafficGenFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, // std Usage() & PrintDefaults() use Stderr
			"\tUsage: %s [OPTIONS] %s [%s OPTIONS]\n", os.Args[0], trafficGenCmd.Name(), trafficGenCmd.Name())
		authFlags.PrintDefaults()
	}
}

func handleAuthCmd(cmd *commands.Command, args []string) int {
	f := cmd.Flags()
	if f.NArg() < 1 {
		fmt.Printf("IMSI argument must be provided\n\n")
		return 1
	}
	if f.NArg() > 1 {
		fmt.Printf("Please provide only an IMSI argument\n")
		return 1
	}
	imsi := strings.TrimSpace(f.Arg(0))
	req := &protos.AuthenticateRequest{
		Imsi:            imsi,
		CalledStationID: "76-02-DE-AD-BE-FF",
	}
	res, err := uesim.Authenticate(req)
	if err != nil || res == nil {
		fmt.Printf("Authenticate Error: %s\n", err)
		return 2
	}
	fmt.Printf("Successfully authenticated: %s\n", imsi)
	return 0
}

func handleAddUeCmd(cmd *commands.Command, args []string) int {

	f := cmd.Flags()
	if f.NArg() < 1 {
		fmt.Printf("IMSI argument must be provided\n\n")
		return 1
	}
	if f.NArg() > 1 {
		fmt.Printf("Please provide only an IMSI argument\n")
		return 1
	}
	imsi := strings.TrimSpace(f.Arg(0))

	ues, err := getConfiguredSubscribers(imsi)
	if err != nil {
		fmt.Printf("Get UE Error: %s\n", err)
		return 2
	}
	err = addUes(ues)
	if err != nil {
		fmt.Printf("%s\n", err)
		return 1
	}
	fmt.Printf("Successfully added: %s\n", imsi)
	return 0
}

func triggerAuthenticationLoop(imsi string, success chan<- int, opError chan<- int, protoError chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	successCount := 0
	opErrorCount := 0
	protocolErrorCount := 0
	secret := getRadiusSecret()
	for start := time.Now(); time.Since(start) < (time.Second * time.Duration(trafficLength)); {
		b := rand.Intn(int(maxBreak))
		time.Sleep(time.Duration(b) * time.Second)
		req := &protos.AuthenticateRequest{
			Imsi: imsi,
		}
		res, err := uesim.Authenticate(req)
		if err != nil {
			opErrorCount = opErrorCount + 1
			continue
		}
		encoded := res.GetRadiusPacket()
		radiusP, err := radius.Parse(encoded, []byte(secret))
		if err != nil {
			protocolErrorCount = protocolErrorCount + 1
			continue
		}
		eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
		if eapMessage == nil {
			protocolErrorCount = protocolErrorCount + 1
			continue
		}
		if !reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode) {
			protocolErrorCount = protocolErrorCount + 1
			continue
		}
		successCount = successCount + 1
	}
	success <- successCount
	opError <- opErrorCount
	protoError <- protocolErrorCount
}

func handleTrafficGenCmd(cmd *commands.Command, args []string) int {
	fmt.Print("***** Running UE Sim Traffic Generator *****")
	ues, err := getConfiguredSubscribers()
	if err != nil {
		fmt.Printf("Getting configured subscribers failed: %s\n", err)
		return 1
	}
	err = addUes(ues)
	if err != nil {
		fmt.Printf("%s\n", err)
		return 1
	}

	success := make(chan int, len(ues))
	opError := make(chan int, len(ues))
	protoError := make(chan int, len(ues))
	var wg sync.WaitGroup
	for _, ue := range ues {
		fmt.Printf("***** Running EAP-AKA authentication loop for subscriber: %s\n", ue.GetImsi())
		wg.Add(1)
		go triggerAuthenticationLoop(ue.GetImsi(), success, opError, protoError, &wg)
	}
	totalSuccess := 0
	totalOpError := 0
	totalProtoError := 0
	go func() {
		for i := range success {
			totalSuccess = totalSuccess + i
		}
	}()
	go func() {
		for f := range opError {
			totalOpError = totalOpError + f
		}
	}()
	go func() {
		for v := range protoError {
			totalProtoError = totalProtoError + v
		}
	}()
	wg.Wait()
	close(success)
	close(opError)
	close(protoError)

	fmt.Printf("***** Final Results: *****\n\tSuccess: %d\n\tOperational Errors: %d\n\tProtocol Errors: %d\n",
		totalSuccess,
		totalOpError,
		totalProtoError,
	)
	if totalOpError != 0 || totalProtoError != 0 {
		return 1
	}
	return 0
}

// getConfiguredSubscribers accepts optional imsi as a parameter and
// returns array of protos.UEConfig. If no parameter is provided, then all the
// UEs configured in the uesim.yml are included in the protos.UEConfig.

func getConfiguredSubscribers(params ...string) ([]*protos.UEConfig, error) {
	var imsiArg string
	uecfg, err := config.GetServiceConfig("", registry.UeSim)
	if err != nil {
		return nil, err
	}
	subscribers, ok := uecfg.RawMap["subscribers"]
	if !ok {
		return nil, fmt.Errorf("could not find 'subscribers' in config file")
	}
	rawMap, ok := subscribers.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("unable to convert %T to map %v", subscribers, rawMap)
	}
	var ues []*protos.UEConfig
	if len(params) > 0 {
		// Return UEConfig for a single IMSI
		imsiArg = params[0]
	}

	// Return UEConfig for all the IMSIs configured in uesim.yml
	for k, v := range rawMap {
		imsi, ok := k.(string)
		if !ok {
			continue
		}
		if imsiArg == "" || imsi == imsiArg {
			subscriberCfg, ok := v.(map[interface{}]interface{})
			if !ok {
				continue
			}
			configMap := &config.ConfigMap{RawMap: subscriberCfg}

			ue, err := createUeConfig(imsi, 0, configMap)
			if err != nil {
				glog.Error(err)
				return nil, err
			}
			ues = append(ues, ue)
		}
	}

	return ues, nil
}

func createUeConfig(imsi string, seq_num uint64, configMap *config.ConfigMap) (*protos.UEConfig, error) {
	var op string
	var msisdn, apn string = DefaultMsisdn, DefaultApn
	var rat int = DefaultRatType
	uecfg, err := config.GetServiceConfig("", registry.UeSim)
	if err != nil {
		op = DefaultOp
	}
	op, err = uecfg.GetString("op")
	if err != nil {
		op = DefaultOp
	}

	// If auth_key is incorrect, return error
	authKey, err := configMap.GetString("auth_key")
	if err != nil {
		return nil, fmt.Errorf("Could not add subscriber due to missing auth_key: %s", err)
	}
	authKeyBytes, err := hex.DecodeString(authKey)
	if err != nil {
		return nil, fmt.Errorf("Could not add subscriber due to incorrect auth key format: %s", err)
	}
	opc, err := crypto.GenerateOpc(authKeyBytes, []byte(op))
	if err != nil {
		return nil, fmt.Errorf("could not generate OPc for subscriber: %s: %s", imsi, err)
	}
	msisdn, err = configMap.GetString("msisdn")
	if err != nil {
		glog.Infof("MSISDN not set for Imsi[%s], setting default MSISDN %s ", imsi, DefaultMsisdn)
	}
	apn, err = configMap.GetString("apn")
	if err != nil {
		glog.Infof("APN not set for Imsi[%s], setting default APN %s", imsi, DefaultApn)
	}
	rat, err = configMap.GetInt("rat")
	if err != nil {
		glog.Infof("RAT-Type not set for Imsi[%s], setting default Rat-Type %d", imsi, DefaultRatType)
	}

	glog.Infof("Creating UE with IMSI:[%s] MSISDN[%s] APN[%s] RAT[%d] ", imsi, msisdn, apn, rat)

	return &protos.UEConfig{
		Imsi:    imsi,
		AuthKey: authKeyBytes,
		AuthOpc: opc[:],
		Seq:     seq_num,
		HsslessCfg: &protos.AuthenticateRequestHssLess{
			Msisdn: msisdn,
			Apn:    apn,
			Rat:    uint32(rat),
		},
	}, nil

}

func getRadiusSecret() string {
	uecfg, err := config.GetServiceConfig("", registry.UeSim)
	if err != nil {
		return DefaultRadiusSecret
	}
	radiusSecret, err := uecfg.GetString("radius_secret")
	if err != nil {
		return DefaultRadiusSecret
	}
	return radiusSecret
}

func main() {
	flag.Parse()
	// Init help for all commands
	flag.Usage = func() {
		cmd := os.Args[0]
		fmt.Printf(
			"\nUsage: \033[1m%s command [OPTIONS]\033[0m\n\n",
			filepath.Base(cmd))
		fmt.Println("Commands:")
		cmdRegistry.Usage()
	}
	flag.Parse()
	cmdName := flag.Arg(0)
	if len(flag.Args()) < 1 || cmdName == "" || cmdName == "help" || cmdName == "h" {
		flag.Usage()
		os.Exit(1)
	}
	cmd := cmdRegistry.Get(cmdName)
	if cmd == nil {
		fmt.Println("\nInvalid Command: ", cmdName)
		flag.Usage()
		os.Exit(1)
	}
	args := os.Args[2:]
	cmd.Flags().Parse(args)
	os.Exit(cmd.Handle(args))
}

func addUes(ues []*protos.UEConfig) error {

	for _, ue := range ues {
		err := uesim.AddUE(ue)
		if err != nil {
			return fmt.Errorf("Adding configured subscribers failed: %s", err)
		}
	}
	return nil
}

func handleDisconnectCmd(cmd *commands.Command, args []string) int {
	f := cmd.Flags()
	if f.NArg() < 1 {
		fmt.Printf("IMSI argument must be provided\n\n")
		return 1
	}
	if f.NArg() > 1 {
		fmt.Printf("Please provide only an IMSI argument\n")
		return 1
	}
	imsi := strings.TrimSpace(f.Arg(0))
	req := &protos.DisconnectRequest{
		Imsi:            imsi,
		CalledStationID: "76-02-DE-AD-BE-FF",
	}
	res, err := uesim.Disconnect(req)
	if err != nil || res == nil {
		fmt.Printf("Disconnect Error: %s\n", err)
		return 2
	}
	fmt.Printf("Successfully Disconnected: %s\n", imsi)
	return 0
}
