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

// Command Line Tool to create & manage Operators, ACLs and Certificates
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/plmn_filter"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/s6a_proxy"
	"magma/feg/gateway/services/s6a_proxy/servicers"
	"magma/feg/gateway/services/s6a_proxy/servicers/test"
	"magma/orc8r/cloud/go/tools/commands"
	orcprotos "magma/orc8r/lib/go/protos"
)

const (
	S6aDiamProductEnv = "S6A_DIAM_PRODUCT"

	S6aProxyServiceName = "s6a_proxy"
	DefaultS6aDiamRealm = "epc.mnc070.mcc722.3gppnetwork.org"
	DefaultS6aDiamHost  = "feg-s6a.epc.mnc070.mcc722.3gppnetwork.org"
)

var (
	cmdRegistry    = new(commands.Map)
	proxyAddr      string
	remoteS6a      bool
	mncLen         int = 3
	s6aAddr        string
	network        string = "sctp"
	localAddr      string
	diamHost       string = "feg-s6a.epc.mnc007.mcc722.3gppnetwork.org"
	diamRealm      string = "epc.mnc007.mcc722.3gppnetwork.org"
	destHost       string = "anahss.ims.telefonica.com.ar"
	destRealm      string = "pre.mnc007.mcc722.3gppnetwork.org"
	testServer     bool
	testServerAddr string
	eutranVectors  int = 3
	utranVectors   int = 0
	useMconfig     bool
	imsiRange      uint64 = 1
	rate           int    = 0
)

type s6aCli interface {
	AuthenticationInformation(
		req *protos.AuthenticationInformationRequest) (*protos.AuthenticationInformationAnswer, error)
}

type s6aProxyCli struct{}

func (s6aProxyCli) AuthenticationInformation(
	req *protos.AuthenticationInformationRequest) (*protos.AuthenticationInformationAnswer, error) {

	return s6a_proxy.AuthenticationInformation(req)
}

type s6aBuiltIn struct {
	impl protos.S6AProxyServer
}

func (s s6aBuiltIn) AuthenticationInformation(
	req *protos.AuthenticationInformationRequest) (*protos.AuthenticationInformationAnswer, error) {

	return s.impl.AuthenticationInformation(context.Background(), req)
}

func init() {
	proxyAddr, _ = registry.GetServiceAddress(registry.S6A_PROXY)
	cmd := cmdRegistry.Add(
		"AIR",
		"Send AIR via s6a_proxy",
		air)
	f := cmd.Flags()
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, // std Usage() & PrintDefaults() use Stderr
			"\tUsage: %s [OPTIONS] %s [%s OPTIONS] <IMSI>\n", os.Args[0], cmd.Name(), cmd.Name())
		f.PrintDefaults()
	}
	f.StringVar(&proxyAddr, "proxy", proxyAddr, "s6a proxy address")
	f.BoolVar(&remoteS6a, "remote_s6a", remoteS6a, "Use orc8r to get to the s6a_proxy (Run it on AGW without proxy flag)")
	f.StringVar(&s6aAddr, "hss_addr", s6aAddr,
		"s6a server (HSS) address - overwrites proxy address and starts local s6a proxy")
	f.StringVar(&network, "network", network, "s6a server (HSS) network: tcp/sctp")
	f.StringVar(&localAddr, "local_addr", localAddr, "s6a client local address to buind to")
	f.StringVar(&diamHost, "host", diamHost, "s6a diam host")
	f.StringVar(&diamRealm, "realm", diamRealm, "s6a diam realm")
	f.StringVar(&destHost, "dhost", destHost, "s6a dest host")
	f.StringVar(&destRealm, "drealm", destRealm, "s6a dest realm")
	f.IntVar(&mncLen, "mnclen", mncLen, "IMSI's MNC part len (2 or 3)")
	f.IntVar(&mncLen, "l", mncLen, "IMSI's MNC part len (2 or 3) - short form")
	f.BoolVar(&testServer, "test", testServer,
		"Start local test s6a server bound to a specified by 'test_addr' or 'hss_addr' address")
	f.StringVar(&testServerAddr, "test_addr", testServerAddr,
		"s6a test server address (defaults to '-hss_addr' if not specified)")
	f.IntVar(&eutranVectors, "eutran_num", eutranVectors, "Number of E-UTRAN vectors to request")
	f.IntVar(&utranVectors, "utran_num", utranVectors, "Number of UTRAN vectors to request")
	f.BoolVar(&useMconfig, "use_mconfig", false,
		"Use local gateway.mconfig configuration for local proxy (if set - starts local s6a proxy)")
	f.Uint64Var(&imsiRange, "range", imsiRange, "Send multiple request with consecutive imsis")
	f.IntVar(&rate, "rate", rate, "Request per second (to be used with range)")

}

// AIR Handler
func air(cmd *commands.Command, args []string) int {
	f := cmd.Flags()
	imsi := strings.TrimSpace(f.Arg(0))
	if f.NArg() != 1 || len(imsi) < 6 {
		f.Usage()
		log.Printf("A single IMSI (6+ long) must be specified.")
		return 1
	}
	imsiNum, err := strconv.ParseUint(imsi, 10, 64)
	if err != nil {
		f.Usage()
		log.Printf("Invalid IMSI '%s': %v", imsi, err)
		return 2
	}
	imsiStr := fmt.Sprintf("%d", imsiNum)
	if mncLen != 2 && mncLen != 3 {
		f.Usage()
		log.Printf("Imvalid MCC Length specified (-mccl %d). Must be 2 or 3", mncLen)
		return 3
	}
	plmnId, err := getPlmnID(imsiStr, mncLen)
	if err != nil {
		f.Usage()
		log.Print(err)
		return 31
	}
	fmt.Printf("Using IMSI: %s; MCC: %s; MNC: %s; PLMN ID: %d\n",
		imsiStr, imsiStr[:3], imsiStr[3:3+mncLen], plmnId)

	clientCfg := &diameter.DiameterClientConfig{
		Host:        diamHost,
		Realm:       diamRealm,
		ProductName: diameter.GetValueOrEnv(diameter.ProductFlag, S6aDiamProductEnv, diameter.DiamProductName),
	}
	serverCfg := &diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
		Addr:      s6aAddr,
		Protocol:  network,
		LocalAddr: localAddr},
		DestHost:  destHost,
		DestRealm: destRealm,
	}

	conf := &servicers.S6aProxyConfig{
		ClientCfg: clientCfg,
		ServerCfg: serverCfg,
		PlmnIds:   plmn_filter.PlmnIdVals{},
	}

	if testServer {
		if len(testServerAddr) == 0 {
			testServerAddr = s6aAddr
		}
		if startTestServer(serverCfg.Protocol, testServerAddr) != nil {
			return 4
		}
	}

	var cli s6aCli
	var peerAddr string
	if len(s6aAddr) > 0 || useMconfig { // use direct HSS connection if address is provided
		fmt.Println("Using builtin S6a_proxy")
		if useMconfig {
			conf = servicers.GetS6aProxyConfigs()
		}
		fmt.Printf("Direct connection:\n\tClient Config: %+v\n\tServer Config: %+v\n", *clientCfg, *serverCfg)

		localProxy, err := servicers.NewS6aProxy(conf)
		if err != nil {
			f.Usage()
			log.Printf("BuiltIn Proxy initialization error: %v", err)
			return 5
		}
		cli = s6aBuiltIn{impl: localProxy}
		peerAddr = conf.ServerCfg.Addr
	} else {
		if remoteS6a {
			fmt.Println("Using S6a_proxy through Orc8r")
			os.Setenv("USE_REMOTE_S6A_PROXY", "true")
		} else {
			fmt.Println("Using local S6a_proxy")
		}

		cli = s6aProxyCli{}
		currAddr, _ := registry.GetServiceAddress(registry.S6A_PROXY)
		if currAddr != proxyAddr {
			ch, cp, err := parseAddr(currAddr)
			if err != nil {
				log.Printf("Internal Error, invalid S6A_PROXY address '%s': %v", currAddr, err)
				cp = 9098
			}
			h, p, err := parseAddr(proxyAddr)
			if err != nil {
				if strings.HasPrefix(err.Error(), "missing port") {
					p = cp
					log.Printf("Missing S6a Proxy Address port, using %d", p)
					h = proxyAddr
				} else {
					f.Usage()
					log.Printf("Invalid S6a Proxy Address '%s': %v", proxyAddr, err)
					return 6
				}
				if len(h) == 0 {
					h = ch
					log.Printf("Missing S6a Proxy Address host, using %s", h)
				}
			}
			registry.AddService(registry.S6A_PROXY, h, p)
		}
		peerAddr = proxyAddr
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
		go func() {
			defer wg.Done()
			var errCli error
			req := &protos.AuthenticationInformationRequest{
				UserName:                      fmt.Sprintf("%0*d", lenOfImsi, imsiNum+iShadow),
				VisitedPlmn:                   plmnId[:],
				NumRequestedEutranVectors:     uint32(eutranVectors),
				ImmediateResponsePreferred:    true,
				NumRequestedUtranGeranVectors: uint32(utranVectors),
			}
			// AIR
			json, errCli := orcprotos.MarshalIntern(req)
			if errCli != nil {
				errCli := fmt.Errorf("Can not marshall request: %s", errCli)
				log.Print(errCli)
				errChann <- errCli
				return
			}
			fmt.Printf("Sending AIR to %s:\n%s\n%+#v\n\n", peerAddr, json, *req)
			r, errCli := cli.AuthenticationInformation(req)
			if errCli != nil || r == nil {
				errCli = fmt.Errorf("GRPC AIR Error: %v", errCli)
				log.Print(errCli)
				errChann <- errCli
				return
			}
			json, errCli = orcprotos.MarshalIntern(r)
			if errCli != nil {
				errCli = fmt.Errorf("Marshal Error %v for result: %+v", errCli, *r)
				errChann <- errCli
				return
			}
			fmt.Printf("Received AIA:\n%s\n%+v\n", json, *r)
		}()
	}

	// go routine to collect the errors
	airErrors := make([]error, 0)
	go func() {
		for err2 := range errChann {
			airErrors = append(airErrors, err2)
		}
		done <- struct{}{}
	}()

	// wait until all air request are done
	wg.Wait()
	close(errChann)
	// wait until all the errors are processed
	<-done
	close(done)

	// check if errors
	if len(airErrors) != 0 {
		log.Printf("Errors found: %d request failed out of %d\n", len(airErrors), imsiRange)
		return 9
	}
	log.Printf("\nAll request (%d) got a response\n", imsiRange)
	return 0
}

func getPlmnID(imsi string, mncLen int) ([3]byte, error) {
	imsiBytes := [6]byte{}
	for i := 0; i < 6; i++ {
		v, err := strconv.Atoi(imsi[i : i+1])
		if err != nil {
			return [3]byte{}, fmt.Errorf("Invalid Digit '%s' in IMSI '%s': %v", imsi[i:i+1], imsi, err)
		}
		imsiBytes[i] = byte(v)
	}
	// see https://www.arib.or.jp/english/html/overview/doc/STD-T63v10_70/5_Appendix/Rel11/29/29272-bb0.pdf#page=73
	plmnId := [3]byte{
		imsiBytes[0] | (imsiBytes[1] << 4),
		imsiBytes[2] | (imsiBytes[5] << 4),
		imsiBytes[3] | (imsiBytes[4] << 4)}
	if mncLen < 3 {
		plmnId[1] |= 0xF0
	}
	return plmnId, nil
}

func parseAddr(addr string) (string, int, error) {
	h, pStr, err := net.SplitHostPort(proxyAddr)
	if err != nil {
		return "", 0, fmt.Errorf("%s: for given address: %s", err, addr)
	}
	p, err := strconv.Atoi(pStr)
	return h, p, err
}

func startTestServer(protocol, address string) error {
	fmt.Printf("Starting Test S6a server on %s: %s\n", protocol, address)
	err := test.StartTestS6aServer(protocol, address, false)
	if err != nil {
		log.Printf("Test S6a server stert error: %v", err)
		return err
	}
	time.Sleep(time.Millisecond * 200)
	return nil
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
