// Copyright 2021 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/magma/magma/src/go/agwd/config"
	"github.com/magma/magma/src/go/protos/magma/capture"
	configpb "github.com/magma/magma/src/go/protos/magma/config"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	configFlag := flag.String(
		"c", "/etc/magma/agwd.json", "Path to config file")
	flag.Parse()

	cfgr := config.NewConfigManager()
	cfgr_err := config.LoadConfigFile(cfgr, *configFlag)
	if cfgr_err != nil {
		println("using default configuration as LoadConfigFile failed with %q", cfgr_err)
	}

	makefile := "/home/vagrant/magma/lte/gateway/python/integ_tests/defs.mk"
	dir := "/home/vagrant/magma/lte/gateway/python/integ_tests"
	out := "/home/vagrant/magma/src/go/capture/gen/resources/%s.golden"

	configConn, err := grpc.Dial(
		config.GetVagrantTarget(
			cfgr.Config().GetVagrantPrivateNetworkIp(),
			cfgr.Config().GetConfigServicePort()).
			Endpoint,
		grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	configClient := configpb.NewConfigClient(configConn)
	defer configConn.Close()

	getCfgResp, err := configClient.GetConfig(ctx, &configpb.GetConfigRequest{})
	if err != nil {
		panic(err)
	}
	spec := &configpb.CaptureConfig_MatchSpec{
		Service: "magma.sctpd.SctpdDownlink",
		Method:  "SendDl",
	}
	ulspec := &configpb.CaptureConfig_MatchSpec{
		Service: "magma.sctpd.SctpdUplink",
		Method:  "SendUl",
	}
	updatedCfg := getCfgResp.Config
	updatedCfg.CaptureConfig.MatchSpecs = []*configpb.CaptureConfig_MatchSpec{ulspec, spec}
	replaceCfgResp, err := configClient.ReplaceConfig(ctx, &configpb.ReplaceConfigRequest{Config: updatedCfg})
	if err != nil {
		panic(err)
	}
	if replaceCfgResp.Config.String() != updatedCfg.String() {
		panic("config not replaced")
	}

	captureConn, err := grpc.Dial(
		config.GetVagrantTarget(
			cfgr.Config().GetVagrantPrivateNetworkIp(),
			cfgr.Config().GetCaptureServicePort()).
			Endpoint,
		grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	captureClient := capture.NewCaptureClient(captureConn)
	defer captureConn.Close()

	runAndCaptureTests(context.Background(), dir, out, captureClient, parsePrecommitTestsFromMakefile(makefile))
}

// parsePrecommitTestsFromMakefile helper function parse the make file for the precommit test names.
func parsePrecommitTestsFromMakefile(path string) []string {
	precommitPrefix := "PRECOMMIT_TESTS = "
	extendedPrefix := "EXTENDED_TESTS = "
	prefix := "s1aptests/test_"
	suffix := " \\"

	fileb, err := os.ReadFile(path)
	if err != nil {
		println("failed to read in makefile")
		panic(err)
	}

	tests := []string{}
	sliceData := strings.Split(string(fileb), "\n")

	for _, line := range sliceData {
		if strings.HasPrefix(line, precommitPrefix) {
			tests = append(tests, strings.TrimSuffix(strings.TrimPrefix(line, precommitPrefix), suffix))
		}

		if strings.HasPrefix(line, prefix) {
			tests = append(tests, strings.TrimSuffix(line, suffix))
		}
		if strings.HasPrefix(line, extendedPrefix) {
			break
		}
	}
	return tests
}

// runAndCaptureTests iterates over the test list, runs the test and captures its output to a file.
func runAndCaptureTests(ctx context.Context, dir, outputPathFmt string, captureClient capture.CaptureClient, tests []string) {
	// flush buffer
	_, err := captureClient.Flush(ctx, &capture.FlushRequest{})
	if err != nil {
		println("failed to flush buffer")
		panic(err)
	}

	for _, test := range tests {
		testN := "TESTS=" + test
		cmd := exec.Command("make", "integ_test", testN)
		cmd.Dir = dir
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			println("failed to start test runner")
			panic(err)
		}

		resp, err := captureClient.Flush(ctx, &capture.FlushRequest{})
		if err != nil {
			println("failed to flush captured calls")
			panic(err)
		}
		// Write out golden file.
		goldenfilename := fmt.Sprintf(outputPathFmt, test)
		println(goldenfilename)
		err = os.WriteFile(goldenfilename, []byte(resp.Recording.String()), 0666)
		if err != nil {
			println("failed to write golden")
			panic(err)
		}
	}
}
