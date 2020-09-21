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

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/service_registry"
	"magma/orc8r/cloud/go/services/service_registry/servicers"
	"magma/orc8r/lib/go/protos"

	"github.com/docker/docker/client"
	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, service_registry.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating service registry service %s", err)
	}
	registryModeEnvValue := os.Getenv(service_registry.ServiceRegistryModeEnvVar)

	var serviceRegistryServicer protos.ServiceRegistryServer
	switch registryModeEnvValue {
	case service_registry.DockerServiceRegistry:
		glog.Infof("Registry Mode set to %s. Creating Docker service registry", service_registry.DockerServiceRegistry)
		dockerCli, err := client.NewEnvClient()
		if err != nil {
			glog.Fatalf("Error creating docker client for service registry servicer: %s", err)
		}
		serviceRegistryServicer = servicers.NewDockerServiceRegistryServicer(dockerCli)
	case service_registry.K8sServiceRegistry:
	default:
		if registryModeEnvValue != "1" {
			glog.Infof("Registry Mode %s is invalid. Defaulting to k8s service registry", registryModeEnvValue)
		} else {
			glog.Infof("Registry Mode set to %s. Creating k8s service registry", service_registry.K8sServiceRegistry)
		}
		config, err := rest.InClusterConfig()
		if err != nil {
			glog.Fatalf("Error querying kubernetes config: %s", err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			glog.Fatalf("Error creating kubernetes clientset: %s", err)
		}
		serviceRegistryServicer, err = servicers.NewKubernetesServiceRegistryServicer(clientset.CoreV1())
		if err != nil {
			glog.Fatal(err)
		}
	}
	protos.RegisterServiceRegistryServer(srv.GrpcServer, serviceRegistryServicer)
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error while running service: %s", err)
	}
}
