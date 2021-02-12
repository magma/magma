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
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/s8_proxy/servicers"
	"magma/feg/gateway/services/s8_proxy/servicers/mock_pgw"
	"magma/orc8r/cloud/go/tools/commands"

	"github.com/golang/glog"
	"github.com/golang/protobuf/jsonpb"
	protobuf_proto "github.com/golang/protobuf/proto"
)

var (
	cmdRegistry   = new(commands.Map)
	proxyAddr     string
	IMSI          string = "123456789012345"
	useMconfig    bool
	useBuiltinCli bool

	testServer     bool
	testServerAddr string = "127.0.0.1:0"
	pgwServerAddr  string

	withecho  bool   = false
	apn       string = "internet.com"
	uTeid     uint
	localPort string = "2123"
)

func init() {

	// Enable logging
	flag.Set("v", "10")             // enable the most verbose logging, can be overwritten by 'v' flag
	flag.Set("logtostderr", "true") // enable printing to console, can be overwritten by 'logtostderr' flag

	// Create Session command
	csCmd := cmdRegistry.Add(
		"CS",
		"Create Session through S8 proxy", createSession)

	csFlags := csCmd.Flags()

	csFlags.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"\tUsage: %s [OPTIONS] %s [%s OPTIONS] <IMSI>\n", os.Args[0], csCmd.Name(), csCmd.Name())
		csFlags.PrintDefaults()
	}

	csFlags.BoolVar(&testServer, "test", testServer,
		fmt.Sprintf("Start local test s8 server bound to default PGW address (%s)", testServerAddr))

	csFlags.StringVar(&localPort, "localport", localPort,
		"S8 local port to run the server")

	csFlags.StringVar(&pgwServerAddr, "server", pgwServerAddr,
		"PGW IP:port to send request with format ip:port")

	csFlags.BoolVar(&withecho, "withecho", withecho,
		fmt.Sprint("Starts s8 proxy checking PGW is alive"))

	csFlags.BoolVar(&useMconfig, "use_mconfig", false,
		"Use local gateway.mconfig configuration for local proxy (if set - starts local s6a proxy)")

	csFlags.BoolVar(&useBuiltinCli, "use_builtincli", true,
		"Use local built in client instead of the running instance on the gateway")

	rand.Seed(time.Now().UTC().UnixNano())
	csFlags.UintVar(&uTeid, "teid", uint(rand.Uint32()),
		"User Plane Teid on the magmaa side. Default is random")

	csFlags.StringVar(&apn, "apn", apn,
		"APN on the request")
}

func createSession(cmd *commands.Command, args []string) int {
	f := cmd.Flags()

	imsi, err := parseArgs(f)
	if err != nil {
		fmt.Println(err)
		cmd.Usage()
		return 1
	}

	// start and configure a test server that will act as a PGW
	if testServer {
		pgwServerAddr, err = startTestServer()
		if err != nil {
			fmt.Println(err)
			return 1
		}
		// If test server is enabled, pgwServerAddr will be overwritten
		// and S8 config too
	}

	conf := &servicers.S8ProxyConfig{
		ClientAddr: fmt.Sprintf(":%s", localPort),
		ServerAddr: pgwServerAddr,
	}

	var cli s8Cli
	// Selection of builtIn Client or S8proxy running on the gateway
	if useMconfig || useBuiltinCli {
		// use builtin proxy (ignore loccal proxy)
		if useMconfig {
			conf = servicers.GetS8ProxyConfig()
		}
		fmt.Printf("Direct connection using built in S8 client: Client Config: %+v\n", *conf)
		var localProxy *servicers.S8Proxy
		if withecho {
			fmt.Println("Send echo message on start Enabled")
			localProxy, err = servicers.NewS8ProxyWithEcho(conf)
		} else {
			localProxy, err = servicers.NewS8Proxy(conf)
		}
		if err != nil {
			fmt.Printf("BuiltIn S8 Proxy initialization error: %v\n", err)
			return 5
		}

		cli = s8BuiltIn{localProxy}
	} else {
		// TODO: use local proxy running on the gateway
		proxyAddr, _ = registry.GetServiceAddress(registry.S8_PROXY)
		cli = s8CliImpl{}
	}

	// Dummy request

	// Create Session Request messagea
	_, offset := time.Now().Zone()
	csReq := &protos.CreateSessionRequestPgw{
		PgwAddrs: pgwServerAddr,
		Imsi:     imsi,
		Msisdn:   "00111",
		Mei:      generateIMEIbasedOnIMSI(imsi),
		ServingNetwork: &protos.ServingNetwork{
			Mcc: "310",
			Mnc: "14",
		},
		RatType: 9,
		BearerContext: &protos.BearerContext{
			Id: 5,
			UserPlaneFteid: &protos.Fteid{
				Ipv4Address: generateRandomIPv4(),
				Ipv6Address: "",
				Teid:        uint32(uTeid),
			},
			Qos: &protos.QosInformation{
				Pci:                     0,
				PriorityLevel:           0,
				PreemptionCapability:    0,
				PreemptionVulnerability: 0,
				Qci:                     0,
				Gbr: &protos.Ambr{
					BrUl: 123,
					BrDl: 234,
				},
				Mbr: &protos.Ambr{
					BrUl: 567,
					BrDl: 890,
				},
			},
		},
		PdnType: protos.PDNType_IPV4,
		Paa: &protos.PdnAddressAllocation{
			Ipv4Address: "10.0.0.10",
			Ipv6Address: "",
			Ipv6Prefix:  0,
		},

		Apn:            apn,
		SelectionMode:  "",
		ApnRestriction: 0,
		Ambr: &protos.Ambr{
			BrUl: 999,
			BrDl: 888,
		},
		Uli: &protos.UserLocationInformation{
			Lac:    1,
			Ci:     2,
			Sac:    3,
			Rac:    4,
			Tac:    5,
			Eci:    6,
			MeNbi:  7,
			EMeNbi: 8,
		},
		IndicationFlag: nil,
		TimeZone: &protos.TimeZone{
			DeltaSeconds:       int32(offset),
			DaylightSavingTime: 0,
		},
	}

	printGRPCMessage("Send: ", csReq)

	res, err := cli.CreateSession(csReq)

	if err != nil {
		fmt.Printf("Create Session cli command failed: %s\n", err)
		return 9
	}
	printGRPCMessage("Received: ", res)
	return 0
}

func parseArgs(f *flag.FlagSet) (string, error) {
	if f.NArg() < 1 {
		fmt.Printf("IMSI not provided. Using default IMSI: %s\n", IMSI)
		return IMSI, nil
	}
	imsi := strings.TrimSpace(f.Arg(0))
	err := validateImsi(imsi)
	if err != nil {
		return "", err
	}
	return imsi, nil
}

func generateIMEIbasedOnIMSI(imsi string) string {
	return fmt.Sprintf("94449%s", imsi[len(imsi)-10:])
}

func generateRandomIPv4() string {
	return fmt.Sprintf("172.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255))
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

func startTestServer() (string, error) {
	// Create and run PGW
	mockPgw, err := mock_pgw.NewStarted(nil, testServerAddr)
	if err != nil {
		return "", fmt.Errorf("Error creating test server mock PGW: +%s\n", err)
	}
	fmt.Printf("Running test server: mock PGW on %s\n", mockPgw.LocalAddr().String())
	pgwServerAddr = mockPgw.LocalAddr().String()

	return pgwServerAddr, nil
}

func printGRPCMessage(prefix string, v interface{}) {
	var payload string
	if pm, ok := v.(protobuf_proto.Message); ok {
		var buf bytes.Buffer
		err := (&jsonpb.Marshaler{EmitDefaults: true, Indent: "\t", OrigName: true}).Marshal(&buf, pm)
		if err == nil {
			payload = buf.String()
		} else {
			payload = fmt.Sprintf("\n\t JSON encoding error: %v; %s", err, buf.String())
		}
	} else {
		payload = fmt.Sprintf("\n\t %T is not proto.Message; %+v", v, v)
	}
	glog.Infof("%s%T: %s", prefix, v, payload)
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
		fmt.Println("\nCommands:")
		cmdRegistry.Usage()
	}

	cmdName := flag.Arg(0)
	if len(flag.Args()) < 1 || cmdName == "" || cmdName == "help" {
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
