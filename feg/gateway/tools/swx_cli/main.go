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
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/swx_proxy"
	"magma/feg/gateway/services/swx_proxy/servicers"
	"magma/orc8r/cloud/go/tools/commands"
	orcprotos "magma/orc8r/lib/go/protos"
)

var (
	cmdRegistry  = new(commands.Map)
	config       servicers.SwxProxyConfig
	imsi         string
	numVectors   uint64 = 3
	useRemote    bool
	loopDelaySec int64
	ignoreErrors bool
)

const (
	DefaultProductName     = "magma"
	DefaultNetworkProtocol = "sctp"
)

type swxClient interface {
	Authenticate(
		req *protos.AuthenticationRequest) (*protos.AuthenticationAnswer, error)
	Register(
		req *protos.RegistrationRequest) (*protos.RegistrationAnswer, error)
}

type swxProxyCli struct{}

func (swxProxyCli) Authenticate(
	req *protos.AuthenticationRequest,
) (*protos.AuthenticationAnswer, error) {
	return swx_proxy.Authenticate(req)
}

func (swxProxyCli) Register(
	req *protos.RegistrationRequest,
) (*protos.RegistrationAnswer, error) {
	return swx_proxy.Register(req)
}

type swxBuiltIn struct {
	impl protos.SwxProxyServer
}

func (s swxBuiltIn) Authenticate(
	req *protos.AuthenticationRequest,
) (*protos.AuthenticationAnswer, error) {
	return s.impl.Authenticate(context.Background(), req)
}

func (s swxBuiltIn) Register(
	req *protos.RegistrationRequest,
) (*protos.RegistrationAnswer, error) {
	return s.impl.Register(context.Background(), req)
}

func init() {
	// Enable logging
	flag.Set("v", "10")             // enable the most verbose logging, can be overwritten by 'v' flag
	flag.Set("logtostderr", "true") // enable printing to console, can be overwritten by 'logtostderr' flag

	flag.BoolVar(
		&useRemote,
		"remote_service",
		false,
		"Use remote SWX service (based on the Gateway control proxy configuration)")
	flag.Int64Var(
		&loopDelaySec,
		"loop_delay",
		0,
		"Loop request indefinitely with specified delay between requests in seconds (<= 0 value disables looping)")
	flag.BoolVar(
		&ignoreErrors,
		"ignore_errors",
		false,
		"Ignore errors & continue requests (only valid with non zero loop_delay)")

	config = servicers.SwxProxyConfig{
		ClientCfg: &diameter.DiameterClientConfig{
			ProductName: DefaultProductName,
		},
		ServerCfg: &diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
			Protocol: DefaultNetworkProtocol,
		}},
		VerifyAuthorization: false,
	}
	marCmd := cmdRegistry.Add("MAR", "Send MAR to HSS", handleSwxCmd)
	sarCmd := cmdRegistry.Add("SAR", "Send SAR (with ServerAssignmentType: REGISTER) to HSS", handleSwxCmd)
	interactiveCmd := cmdRegistry.Add("i", "Run CLI in Interactive Mode", runSwxInteractiveMode)
	marFlags := marCmd.Flags()
	sarFlags := sarCmd.Flags()
	iFlags := interactiveCmd.Flags()
	marFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, // std Usage() & PrintDefaults() use Stderr
			"\tUsage: %s [OPTIONS] %s [%s OPTIONS] <IMSI>\n", os.Args[0], marCmd.Name(), marCmd.Name())
		marFlags.PrintDefaults()
	}
	marFlags.StringVar(
		&config.ServerCfg.Addr,
		"hss_addr",
		config.ServerCfg.Addr,
		"HSS address - use to send requests directly to HSS")
	marFlags.StringVar(
		&config.ServerCfg.Protocol,
		"network",
		config.ServerCfg.Protocol,
		"HSS network: tcp/sctp")
	marFlags.StringVar(
		&config.ServerCfg.LocalAddr,
		"local_addr",
		config.ServerCfg.LocalAddr,
		"swx client local address to bind to")
	marFlags.StringVar(&config.ClientCfg.Host, "origin_host", config.ClientCfg.Host, "swx origin host")
	marFlags.StringVar(&config.ClientCfg.Realm, "origin_realm", config.ClientCfg.Realm, "swx origin realm")
	marFlags.StringVar(&config.ServerCfg.DestHost, "dest_host", config.ServerCfg.DestHost, "swx destination host")
	marFlags.StringVar(
		&config.ServerCfg.DestRealm,
		"dest_realm",
		config.ServerCfg.DestRealm,
		"swx destination realm")
	marFlags.Uint64Var(&numVectors, "num_vectors", numVectors, "number of authentication vectors requested")
	marFlags.BoolVar(
		&config.VerifyAuthorization,
		"verify_authorization",
		config.VerifyAuthorization,
		"Ensure that subscriber has NON-3GPP-IP-Access enabled")

	// Use the same flag set for both MAR and SAR
	*sarFlags = *marFlags
	sarFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, // std Usage() & PrintDefaults() use Stderr
			"\tUsage: %s [OPTIONS] %s [%s OPTIONS] <IMSI>\n", os.Args[0], sarCmd.Name(), sarCmd.Name())
		sarFlags.PrintDefaults()
	}
	iFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, // std Usage() & PrintDefaults() use Stderr
			"\tUsage: %s %s\n", os.Args[0], interactiveCmd.Name())
	}
}

func handleSwxCmd(cmd *commands.Command, args []string) int {
	f := cmd.Flags()
	if f.NArg() < 1 {
		fmt.Print("IMSI argument must be provided\n\n")
		return 1
	}
	if f.NArg() > 1 {
		fmt.Printf("Please provide only an IMSI argument - all other parameters should be provided with flags: %+v\n", f.Args())
		return 1
	}
	imsi = strings.TrimSpace(f.Arg(0))
	err := validateImsi(imsi)
	if err != nil {
		fmt.Print(err.Error())
		cmd.Usage()
		return 1
	}
	return sendSwxRequest(cmd.Name())
}

func sendSwxRequest(requestName string) int {
	requestName = strings.ToLower(requestName)
	var swxCli swxClient
	var addr string
	// Use built-in proxy
	if len(config.ServerCfg.Addr) > 0 {
		swxProxyBuiltIn, err := servicers.NewSwxProxy(&config)
		if err != nil {
			fmt.Print(err.Error())
			return 1
		}
		swxCli = swxBuiltIn{swxProxyBuiltIn}
		addr = config.ServerCfg.Addr
	} else {
		swxCli = swxProxyCli{}
		if useRemote {
			addr = "<REMOTE Address>"
		} else {
			addr, _ = registry.GetServiceAddress(registry.SWX_PROXY)
		}
	}
	if requestName == "sar" {
		return sendSar(addr, swxCli)
	}
	return sendMar(addr, swxCli)
}

func sendSar(addr string, client swxClient) int {
	req := &protos.RegistrationRequest{
		UserName: imsi,
	}
	json, err := orcprotos.MarshalIntern(req)
	if err != nil {
		fmt.Print("Unable to convert request to JSON for printing; Still attempting to send request...")
	} else {
		fmt.Printf("Sending SAR (REGISTER) to %s:\n%s\n%+#v\n\n", addr, json, *req)
	}
	res, err := client.Register(req)
	if err != nil || res == nil {
		fmt.Printf("Register Error: %s\n", err)
		return 2
	}
	fmt.Printf("Successfully registered %s\n", imsi)
	return 0
}

func sendMar(addr string, client swxClient) int {
	req := &protos.AuthenticationRequest{
		UserName:             imsi,
		SipNumAuthVectors:    uint32(numVectors),
		AuthenticationScheme: protos.AuthenticationScheme_EAP_AKA,
		ResyncInfo:           nil,
		RetrieveUserProfile:  true,
	}
	json, err := orcprotos.MarshalIntern(req)
	if err != nil {
		fmt.Printf("Unable to convert request to JSON for printing; Still attempting to send request...")
	} else {
		fmt.Printf("Sending MAR to %s:\n%s\n%+#v\n\n", addr, json, *req)
	}
	res, err := client.Authenticate(req)
	if err != nil || res == nil {
		fmt.Printf("Authenticate Error: %s\n", err)
		return 2
	}
	json, err = orcprotos.MarshalIntern(res)
	if err != nil {
		fmt.Printf("Marshal Error %v for result: %+v", err, *res)
		return 3
	}
	fmt.Printf("Received successful MAA:\n%s\n%+v\n", json, *res)
	return 0
}

func sendInteractiveBuiltinRequest(reader *bufio.Reader, requestType string) int {
	err := getInteractiveRequestParameters(reader, requestType)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	config.ServerCfg.DestHost, err = getUserInput("Please provide the destination host: ", reader)
	if err != nil {
		return 1
	}
	config.ServerCfg.DestRealm, err = getUserInput("Please provide the destination realm: ", reader)
	if err != nil {
		return 1
	}
	config.ClientCfg.Host, err = getUserInput("Please provide the origin host: ", reader)
	if err != nil {
		return 1
	}
	config.ClientCfg.Realm, err = getUserInput("Please provide the origin realm: ", reader)
	if err != nil {
		return 1
	}
	config.ServerCfg.Addr, err = getUserInput("Please provide the HSS Address in format <addr>:<portNo> ", reader)
	if err != nil {
		return 1
	}
	config.ServerCfg.LocalAddr, err = getUserInput("Please provide Local Address in format <addr>:<portNo> ", reader)
	if err != nil {
		return 1
	}
	return sendSwxRequest(requestType)
}

func sendInteractiveProxiedRequest(reader *bufio.Reader, requestType string) int {
	err := getInteractiveRequestParameters(reader, requestType)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	return sendSwxRequest(requestType)
}

func runSwxInteractiveMode(cmd *commands.Command, args []string) int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("------------------------\n")
	fmt.Printf("SWx CLI Interactive Mode\n")
	fmt.Printf("------------------------\n\n")
	answer, err := getUserInput("Do you want to send the request through the FeG's swx_proxy service? (Y/[N]): ", reader)
	if err != nil {
		return 1
	}
	answer = strings.ToLower(answer)
	if answer != "y" && answer != "n" {
		fmt.Println("Please provide either 'Y' for proxy service request or 'N' for direct request")
		return 1
	}
	requestType, err := getUserInput("Do you want to send an MAR or SAR Swx message? (MAR/[SAR]): ", reader)
	if err != nil {
		return 1
	}
	requestType = strings.ToLower(requestType)
	if requestType != "sar" && requestType != "mar" {
		fmt.Println("Please provide either 'MAR' or 'SAR' request type")
		return 1
	}
	switch answer {
	case "y":
		return sendInteractiveProxiedRequest(reader, requestType)
	case "n":
		return sendInteractiveBuiltinRequest(reader, requestType)
	}
	fmt.Println("Please provide either 'Y' for proxy service request or 'N' for direct request")
	return 1
}

func main() {
	flag.Parse()
	// Init help for all commands
	flag.Usage = func() {
		cmd := os.Args[0]
		fmt.Printf(
			"\nUsage: \033[1m%s command [OPTIONS]\033[0m\n\n",
			filepath.Base(cmd))
		flag.PrintDefaults()
		fmt.Println("Commands:")
		cmdRegistry.Usage()
	}
	flag.Parse()
	if useRemote {
		os.Setenv("USE_REMOTE_SWX_PROXY", "true")
	} else {
		os.Setenv("USE_REMOTE_SWX_PROXY", "false")
	}
	loopInterval := time.Second * time.Duration(loopDelaySec)
	for {
		exitCode, err := cmdRegistry.HandleCommand()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			flag.Usage()
		}
		if loopInterval <= 0 || (exitCode != 0 && (!ignoreErrors)) {
			os.Exit(exitCode)
		}
		time.Sleep(loopInterval)
	}
}

func getInteractiveRequestParameters(reader *bufio.Reader, requestType string) error {
	imsiValue, err := getUserInput("Please provide an IMSI value: ", reader)
	if err != nil {
		return err
	}
	err = validateImsi(imsiValue)
	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	imsi = imsiValue
	switch requestType {
	case "mar":
		vectorsStr, err := getUserInput("Please provide a requested number of auth vectors (default 3): ", reader)
		if err != nil {
			return err
		}
		if vectorsStr != "" {
			numVectors, err = strconv.ParseUint(vectorsStr, 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
		}
		return nil
	case "sar":
		return nil
	}
	err = fmt.Errorf("Invalid request type %s provided", requestType)
	fmt.Println(err.Error())
	return err
}

func getUserInput(prompt string, reader *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	if reader == nil {
		err := fmt.Errorf("Nil IO reader provided")
		fmt.Print(err.Error())
		return "", err
	}
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Print(err.Error())
		return "", err
	}
	return strings.TrimSuffix(input, "\n"), nil
}

func validateImsi(imsi string) error {
	if len(imsi) < 6 || len(imsi) > 15 {
		return fmt.Errorf("The IMSI specified must be 6 - 15 digits long\n\n")
	}
	_, err := strconv.ParseUint(imsi, 10, 64)
	if err != nil {
		return fmt.Errorf("Invalid IMSI '%s': %v\n\n", imsi, err)
	}
	return nil
}
