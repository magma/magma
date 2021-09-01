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
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/s8_proxy/servicers"
	"magma/feg/gateway/services/s8_proxy/servicers/mock_pgw"
	"magma/orc8r/cloud/go/tools/commands"

	"github.com/golang/glog"
	"github.com/golang/protobuf/jsonpb"
	protobuf_proto "github.com/golang/protobuf/proto"
	"github.com/wmnsk/go-gtp/gtpv2"
)

var (
	cmdRegistry   = new(commands.Map)
	proxyAddr     string
	remoteS8      bool
	IMSI          string = "123456789012345"
	useMconfig    bool
	useBuiltinCli bool

	testServer     bool
	testServerAddr string = "127.0.0.1:0"
	pgwServerAddr  string

	withecho            bool   = false
	apn                 string = "internet.com"
	AGWTeidU            uint
	localPort           string = "2123"
	createDeleteTimeout int    = -1
	bearerId            uint32 = 5
	rattype             uint   = 6
	AGWTeidC            uint
	imsiRange           uint64 = 1
	rate                int    = 0
	disableGRPClog      bool   = false
)

func init() {

	// Enable logging
	flag.Set("v", "10")             // enable the most verbose logging, can be overwritten by 'v' flag
	flag.Set("logtostderr", "true") // enable printing to console, can be overwritten by 'logtostderr' flag

	// Create Session command
	csCmd := cmdRegistry.Add(
		"cs",
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

	csFlags.BoolVar(&remoteS8, "remote_s8", remoteS8, "Use orc8r to get to the s0_proxy (Run it on AGW without proxy flag)")

	csFlags.StringVar(&pgwServerAddr, "server", pgwServerAddr,
		"PGW IP:port to send request with format ip:port")

	csFlags.BoolVar(&withecho, "withecho", withecho,
		"Starts s8 proxy checking PGW is alive")

	csFlags.BoolVar(&useMconfig, "use_mconfig", false,
		"Use local gateway.mconfig configuration for local proxy (if set - starts local s6a proxy)")

	csFlags.BoolVar(&useBuiltinCli, "use_builtincli", true,
		"Use local built in client instead of the running instance on the gateway")

	rand.Seed(time.Now().UTC().UnixNano())
	csFlags.UintVar(&AGWTeidU, "uteid", uint(rand.Uint32()),
		"User Plane Teid on the magma side. Default is random")

	csFlags.UintVar(&AGWTeidC, "cteid", uint(rand.Uint32()),
		"Control Plane Teid on the magma side. Default is random")

	csFlags.StringVar(&apn, "apn", apn,
		"APN on the request")

	csFlags.IntVar(&createDeleteTimeout, "delete", -1,
		"Use in case you want to delete the session -delete 3 will wait 3 seconds before deletion")

	csFlags.UintVar(&rattype, "rat", 6,
		"Rat type (by default 6 which meanes EUTRAN)")

	csFlags.Uint64Var(&imsiRange, "range", imsiRange, "Send multiple request with consecutive imsis")

	csFlags.IntVar(&rate, "rate", rate, "Request per second (to be used with range)")

	csFlags.BoolVar(&disableGRPClog, "no_grpc_print", false,
		"disable GRPC printing logs")

	// Echo command
	eCmd := cmdRegistry.Add(
		"e",
		"Send echo request through S8 proxy", sendEcho)

	eFlags := eCmd.Flags()

	eFlags.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"\tUsage: %s [OPTIONS] %s [%s OPTIONS]\n", os.Args[0], eCmd.Name(), eCmd.Name())
		eFlags.PrintDefaults()
	}

	eFlags.BoolVar(&testServer, "test", testServer,
		fmt.Sprintf("Start local test s8 server bound to default PGW address (%s)", testServerAddr))

	eFlags.BoolVar(&remoteS8, "remote_s8", remoteS8, "Use orc8r to get to the s0_proxy (Run it on AGW without proxy flag)")

	eFlags.StringVar(&localPort, "localport", localPort,
		"S8 local port to run the server")

	eFlags.StringVar(&pgwServerAddr, "server", pgwServerAddr,
		"PGW IP:port to send request with format ip:port")

	eFlags.BoolVar(&withecho, "withecho", withecho,
		"Starts s8 proxy checking PGW is alive")

	eFlags.BoolVar(&useMconfig, "use_mconfig", false,
		"Use local gateway.mconfig configuration for local proxy (if set - starts local s8 proxy)")

	eFlags.BoolVar(&useBuiltinCli, "use_builtincli", false,
		"Use local built in client instead of the running instance on the gateway")

}

func createSession(cmd *commands.Command, args []string) int {
	if disableGRPClog {
		fmt.Println("no_grpc flag detected. GRPC will NOT be displayed")
	}
	cli, f, err := initialize(cmd, args)
	if err != nil {
		fmt.Print(err)
		return 1
	}

	imsi, err := parseArgs(f)
	if err != nil {
		fmt.Println(err)
		cmd.Usage()
		return 2
	}
	imsiNum, err := strconv.ParseUint(imsi, 10, 64)
	if err != nil {
		fmt.Printf("imsi is not numeric %+v: %s\n", imsi, err)
		return 2
	}
	errChann := make(chan error)
	done := make(chan struct{})
	lenOfImsi := len(imsi)
	wg := sync.WaitGroup{}
	// run all the producers in parallel
	fmt.Printf("Start sending requests %d\n", imsiRange)
	for i := uint64(0); i < imsiRange; i++ {

		wg.Add(1)
		iShadow := i
		// wait to adjust the rate
		if rate > 0 && i != 0 && i%uint64(rate) == 0 {
			fmt.Printf("\nWait 1 sec to send next group of %d request\n\n", rate)
			time.Sleep(time.Second)
		}

		_, offset := time.Now().Zone()
		go func() {
			var errCli error
			defer wg.Done()
			// Create Session Request message
			currentImsi := fmt.Sprintf("%0*d", lenOfImsi, imsiNum+iShadow)
			currentMei := fmt.Sprint(uint64(2000000000000) + iShadow)
			currentCAgwTeid := uint32(AGWTeidC + uint(iShadow))

			csReq := &protos.CreateSessionRequestPgw{
				PgwAddrs: pgwServerAddr,
				Imsi:     currentImsi,
				Msisdn:   currentMei,
				Mei:      generateIMEIbasedOnIMSI(currentImsi),
				CAgwTeid: currentCAgwTeid,
				ServingNetwork: &protos.ServingNetwork{
					Mcc: "310",
					Mnc: "14",
				},
				RatType: protos.RATType_EUTRAN,
				BearerContext: &protos.BearerContext{
					Id: bearerId,
					UserPlaneFteid: &protos.Fteid{
						Ipv4Address: "11.11.11.11",
						Ipv6Address: "",
						Teid:        uint32(AGWTeidU),
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

				Apn:           apn,
				SelectionMode: protos.SelectionModeType_APN_provided_subscription_verified,
				Ambr: &protos.Ambr{
					BrUl: 999,
					BrDl: 888,
				},
				Uli: &protos.UserLocationInformation{
					Tac: 5,
					Eci: 6,
				},

				ProtocolConfigurationOptions: &protos.ProtocolConfigurationOptions{
					ConfigProtocol: uint32(gtpv2.ConfigProtocolPPPWithIP),
					ProtoOrContainerId: []*protos.PcoProtocolOrContainerId{
						{
							Id:       uint32(gtpv2.ProtoIDIPCP),
							Contents: []byte{0x01, 0x00, 0x00, 0x10, 0x03, 0x06, 0x01, 0x01, 0x01, 0x01, 0x81, 0x06, 0x02, 0x02, 0x02, 0x02},
						},
						{
							Id:       uint32(gtpv2.ProtoIDPAP),
							Contents: []byte{0x01, 0x00, 0x00, 0x0c, 0x03, 0x66, 0x6f, 0x6f, 0x03, 0x62, 0x61, 0x72},
						},
						{
							Id:       uint32(gtpv2.ContIDMSSupportOfNetworkRequestedBearerControlIndicator),
							Contents: nil,
						},
					},
				},
				IndicationFlag: nil,
				TimeZone: &protos.TimeZone{
					DeltaSeconds:       int32(offset),
					DaylightSavingTime: 1,
				},
			}

			fmt.Println("\n *** Create Session Test ***")
			printGRPCMessage("Sending GRPC message: ", csReq)

			csRes, errCli := cli.CreateSession(csReq)

			if errCli != nil {
				errCli = fmt.Errorf("=> Create Session cli command failed: %s\n", errCli)
				fmt.Print(errCli)
				errChann <- errCli
				return
			}
			printGRPCMessage("Received GRPC message: ", csRes)

			// check if message was received but GTP message received was in fact an error
			if csRes.GtpError != nil {
				errCli = fmt.Errorf("Received a GTP error (see the GRPC message before): %d\n", csRes.GtpError.Cause)
				fmt.Println(errCli)
				errChann <- errCli
				return
			}

			// Delete recently created session (if enableD)
			if createDeleteTimeout != -1 {
				fmt.Printf("\n=> Sleeping for %ds before deleting....", createDeleteTimeout)
				time.Sleep(time.Duration(createDeleteTimeout) * time.Second)
				fmt.Println(" Done")

				fmt.Println("\n *** Delete Session Test ***")
				dsReq := &protos.DeleteSessionRequestPgw{
					PgwAddrs:       pgwServerAddr,
					Imsi:           imsi,
					BearerId:       bearerId,
					CAgwTeid:       uint32(AGWTeidC),
					CPgwTeid:       csRes.CPgwFteid.Teid,
					ServingNetwork: csReq.ServingNetwork,
					Uli:            csReq.Uli,
				}
				printGRPCMessage("Sending GRPC message: ", dsReq)
				dsRes, errCli := cli.DeleteSession(dsReq)
				if errCli != nil {
					errCli = fmt.Errorf("=> Delete session failed: %s\n", errCli)
					fmt.Println(errCli)
					errChann <- errCli
					return
				}
				printGRPCMessage("Received GRPC message: ", dsRes)
			}
		}()
	}

	// go routine to collect the errors
	errsCli := make([]error, 0)
	go func() {
		for err2 := range errChann {
			errsCli = append(errsCli, err2)
		}
		done <- struct{}{}
	}()

	// wait until all request are done
	wg.Wait()
	close(errChann)
	// wait until all the errors are processed
	<-done
	close(done)

	// check if errors
	if len(errsCli) != 0 {
		fmt.Printf("Errors found: %d request failed out of %d\n", len(errsCli), imsiRange)
		return 9
	}
	fmt.Printf("\nAll request (%d) got a response\n", imsiRange)
	return 0
}

func sendEcho(cmd *commands.Command, args []string) int {
	cli, _, err := initialize(cmd, args)
	if err != nil {
		fmt.Print(err)
		return 1
	}
	fmt.Printf("=> Sending Echo request to %s\n", pgwServerAddr)
	dsReq := &protos.EchoRequest{PgwAddrs: pgwServerAddr}
	dsRes, err := cli.SendEcho(dsReq)
	if err != nil {
		fmt.Printf("=> Delete session failed: %s\n", err)
		return 9
	}
	printGRPCMessage("=> Received: ", dsRes)
	return 0
}

func initialize(cmd *commands.Command, args []string) (s8Cli, *flag.FlagSet, error) {
	f := cmd.Flags()
	var err error
	var cli s8Cli

	// start and configure a test server that will act as a PGW
	if testServer {
		// ONLY USE BUILTIN CLI
		useBuiltinCli = true
		pgwServerAddr, err = startTestServer()
		if err != nil {
			return nil, nil, err
		}
		// If test server is enabled, pgwServerAddr will be overwritten
		// and S8 config too
	}

	conf := &servicers.S8ProxyConfig{
		ClientAddr: fmt.Sprintf(":%s", localPort),
	}

	// Selection of builtIn Client or S8proxy running on the gateway
	if useMconfig || useBuiltinCli {
		// use builtin proxy (ignore loccal proxy)
		fmt.Println("Using builtin S8_proxy")
		if useMconfig {
			conf = servicers.GetS8ProxyConfig()
		}
		fmt.Printf("=> Direct connection using built in S8 client: Client Config: %+v\n", *conf)
		var localProxy *servicers.S8Proxy
		if withecho {
			fmt.Println("=> Send echo message on start Enabled")
			localProxy, err = servicers.NewS8ProxyWithEcho(conf)
		} else {
			localProxy, err = servicers.NewS8Proxy(conf)
		}
		if err != nil {
			return nil, nil, fmt.Errorf("=> BuiltIn S8 Proxy initialization error: %v\n", err)
		}
		cli = s8BuiltIn{localProxy}
	} else {
		if remoteS8 {
			fmt.Println("Using S8_proxy through Orc8r")
			os.Setenv("USE_REMOTE_S8_PROXY", "true")
		} else {
			fmt.Println("Using local S8_proxy")
		}
		proxyAddr, _ = registry.GetServiceAddress(registry.S8_PROXY)
		cli = s8CliImpl{}
	}

	return cli, f, nil
}

func parseArgs(f *flag.FlagSet) (string, error) {
	if f.NArg() < 1 {
		fmt.Printf("=> IMSI not provided. Using default IMSI: %s\n", IMSI)
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
		return fmt.Errorf("=> The IMSI specified must be 6 - 15 digits long\n\n")
	}
	_, err := strconv.ParseUint(imsi, 10, 64)
	if err != nil {
		return fmt.Errorf("Invalid IMSI '%s': %v\n\n", imsi, err)
	}
	return nil
}

func startTestServer() (string, error) {
	// Create and run PGW
	mockPgw, err := mock_pgw.NewStarted(context.Background(), testServerAddr)
	if err != nil {
		return "", fmt.Errorf("Error creating test server mock PGW: +%s\n", err)
	}
	fmt.Printf("Running test server: mock PGW on %s\n", mockPgw.LocalAddr().String())
	pgwServerAddr = mockPgw.LocalAddr().String()

	return pgwServerAddr, nil
}

func printGRPCMessage(prefix string, v interface{}) {
	if disableGRPClog {
		return
	}
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
		fmt.Println("Example:")
		fmt.Println("./s8_cli cs -server 172.16.1.2:2123 -use_builtincli -delete 3 -apn roam  001020000000066 -logtostderr")
		fmt.Printf(
			"\nUsage: \033[1m%s command [OPTIONS] <IMSI> [DEFAULTS]\033[0m\n\n",
			filepath.Base(cmd))
		fmt.Println("Defaults:")
		flag.PrintDefaults()
		fmt.Println("\nCommands:")
		cmdRegistry.Usage()
	}

	cmdName := flag.Arg(0)
	if cmdName == "" || cmdName == "help" {
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
