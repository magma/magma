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
//
package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/magma/magma/src/go/agwd/config"
	pipelinedpb "github.com/magma/magma/src/go/protos/magma/pipelined"
	"google.golang.org/grpc"
)


// Script to dial Pipelined/GetStats gRPC and print out the results
// useful for build/debug purposes
func main() {
	configFlag := flag.String(
		"c", "/etc/magma/agwd.json", "Path to config file")
	flag.Parse()

	cfgr := config.NewConfigManager()
	cfgr_err := config.LoadConfigFile(cfgr, *configFlag)
	if cfgr_err != nil {
		fmt.Printf("using default configuration as LoadConfigFile failed with %q\n", cfgr_err)
	}

	target := config.ParseTarget(cfgr.Config().GetPipelinedServiceTarget())
	conn, err := grpc.Dial(target.Endpoint, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("could not dial %s: %q\n", target, err))
	}
	defer conn.Close()
	ctx := context.Background()
	req := &pipelinedpb.GetStatsRequest{
		Cookie: 0,
		CookieMask: 0,
	}

	res, err := pipelinedpb.NewPipelinedClient(conn).GetStats(ctx, req)
	if err != nil {
		panic(fmt.Sprintf("could not call GetStats %s: %q\n", target, err))
	}
	fmt.Printf("response is %+v", res)
}
