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
	"magma/feg/gateway/services/session_proxy/credit_control/gy"
	"magma/feg/gateway/services/testcore/ocs/mock_ocs"
	"magma/orc8r/lib/go/service"

	"github.com/golang/glog"
)

const (
	MaxUsageBytes = 2048
	MaxUsageTime  = 1000 // in second
	ValidityTime  = 60   // in second
)

var (
	serverNumber int
)

func init() {
	flag.IntVar(&serverNumber, "servernumber", 1, "Number of the server. Will use Gy[servernumber-1] configuration")
}

func main() {
	flag.Parse()

	serverIdx := serverNumber - 1
	log.Print("------ Reading Gy configuration a couple of times ------")
	// Get the server N from the list of configured servers. This is normally 0, unless multiple GX connections are configured
	gyConfigs := gy.GetGyClientConfiguration()
	if serverIdx >= len(gyConfigs) {
		log.Fatalf("ServerIndex value (%d) is bigger than the amount Gy servers configured (%d)", serverIdx, len(gyConfigs))
		return
	}
	gyCliConf := gyConfigs[serverIdx]
	gyServConf := gy.GetOCSConfiguration()[serverIdx]
	log.Print("------ Done reading Gy configuration  ------")
	log.Printf("Mock OCS using Gy server configured at index %d with adderess %s", serverIdx, gyServConf.Addr)

	serviceName := registry.MOCK_OCS
	if serverIdx > 0 {
		serviceName = fmt.Sprintf("%s%d", serviceName, serverNumber)
		log.Printf("OCS serviceName renamed to: %s", serviceName)

	}

	diamServer := mock_ocs.NewOCSDiamServer(
		gyCliConf,
		&mock_ocs.OCSConfig{
			ServerConfig:        gyServConf,
			MaxUsageOctets:      &protos.Octets{TotalOctets: MaxUsageBytes},
			MaxUsageTime:        MaxUsageTime,
			ValidityTime:        ValidityTime,
			GyInitMethod:        gy.PerSessionInit,
			FinalUnitIndication: mock_ocs.FinalUnitIndication{FinalUnitAction: protos.FinalUnitAction_Terminate},
		},
	)

	srv, err := service.NewServiceWithOptions(registry.ModuleName, serviceName)
	if err != nil {
		log.Fatalf("Error creating mock %s service: %s", serviceName, err)
	}

	lis, err := diamServer.StartListener()
	if err != nil {
		log.Fatalf("Unable to start listener for mock %s: %s", serviceName, err)
	}

	protos.RegisterMockOCSServer(srv.GrpcServer, diamServer)

	go func() {
		glog.V(2).Infof("Starting mock %s server at %s", serviceName, lis.Addr().String())
		glog.Errorf(diamServer.Start(lis).Error()) // blocks
	}()

	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running mock %s service: %s", serviceName, err)
	}
}
