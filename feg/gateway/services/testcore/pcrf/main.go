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
	"flag"
	"fmt"
	"log"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/session_proxy/credit_control/gx"
	"magma/feg/gateway/services/testcore/pcrf/mock_pcrf"
	"magma/orc8r/lib/go/service"

	"github.com/golang/glog"
)

var (
	serverNumber int
)

func init() {
	flag.IntVar(&serverNumber, "servernumber", 1, "Number of the server. Will use Gx[servernumber-1] configuration")
}

func main() {
	flag.Parse()

	serverIdx := serverNumber - 1
	log.Print("------ Reading Gx configuration a couple of times ------")
	// Get the server N from the list of configured servers. This is normally 0, unless multiple GX connections are configured
	gxConfigs := gx.GetGxClientConfiguration()
	if serverIdx >= len(gxConfigs) {
		log.Fatalf("ServerIndex value (%d) is bigger than the amount Gx servers configured (%d)", serverIdx, len(gxConfigs))
		return
	}
	gxCliConf := gxConfigs[serverIdx]
	gxServConf := gx.GetPCRFConfiguration()[serverIdx]
	log.Print("------ Done reading Gy configuration  ------")
	log.Printf("Mock PCRF using Gx server configured at index %d with adderess %s", serverIdx, gxServConf.Addr)

	serviceName := registry.MOCK_PCRF
	if serverIdx > 0 {
		serviceName = fmt.Sprintf("%s%d", serviceName, serverNumber)
		log.Printf("PCRF serviceName renamed to: %s", serviceName)
	}

	pcrfServer := mock_pcrf.NewPCRFServer(gxCliConf, gxServConf)

	srv, err := service.NewServiceWithOptions(registry.ModuleName, serviceName)
	if err != nil {
		log.Fatalf("Error creating mock %s service: %s", serviceName, err)
	}

	lis, err := pcrfServer.StartListener()
	if err != nil {
		log.Fatalf("Unable to start listener for mock %s: %s", serviceName, err)
	}

	protos.RegisterMockPCRFServer(srv.GrpcServer, pcrfServer)

	go func() {
		glog.V(2).Infof("Starting mock %s server at %s", serviceName, lis.Addr().String())
		glog.Errorf(pcrfServer.Start(lis).Error()) // blocks
	}()

	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running mock %s service: %s", serviceName, err)
	}
}
