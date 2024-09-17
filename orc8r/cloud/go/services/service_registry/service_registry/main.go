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
	"os"

	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	servicers "magma/orc8r/cloud/go/services/service_registry/servicers/protected"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"
)

const (
	defaultK8sQPS   = 50
	defaultK8sBurst = 50

	pollFrequencyConfigKey = "poll_frequency"
)

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, registry.ServiceRegistryServiceName)
	if err != nil {
		glog.Fatalf("Error creating service registry service %s", err)
	}
	registryModeEnvValue := os.Getenv(registry.ServiceRegistryModeEnvVar)
	if registryModeEnvValue == registry.K8sRegistryMode {
		glog.Infof("Registry Mode set to %s. Creating k8s service registry", registry.K8sRegistryMode)
		config, err := rest.InClusterConfig()
		if err != nil {
			glog.Fatalf("Error querying kubernetes config: %s", err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			glog.Fatalf("Error creating kubernetes clientset: %s", err)
		}
		servicer, err := servicers.NewKubernetesServiceRegistryServicer(clientset.CoreV1(), srv.Config.MustGetString(pollFrequencyConfigKey), nil)
		if err != nil {
			glog.Fatal(err)
		}
		protos.RegisterServiceRegistryServer(srv.ProtectedGrpcServer, servicer)
	} else {
		glog.Infof("Registry Mode set to %s. Not creating service registry servicer", registryModeEnvValue)
	}

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running service: %s", err)
	}
}
