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

package servicers

import (
	"fmt"
	"os"
	"strings"

	"magma/orc8r/lib/go/protos"

	"golang.org/x/net/context"
	corev1types "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	partOfLabel                    = "app.kubernetes.io/part-of"
	partOfOrc8rApp                 = "orc8r-app"
	serviceRegistryNamespaceEnvVar = "SERVICE_REGISTRY_NAMESPACE"
	grpcPortName                   = "grpc"
	httpPortName                   = "http"
	serviceNameDelimiter           = "-"
)

type KubernetesServiceRegistryServicer struct {
	client    corev1.CoreV1Interface
	namespace string
}

// NewKubernetesServiceRegistryServicer creates a new service registry servicer
// that is backed by Kubernetes.
func NewKubernetesServiceRegistryServicer(k8sClient corev1.CoreV1Interface) (*KubernetesServiceRegistryServicer, error) {
	namespaceEnvValue := os.Getenv(serviceRegistryNamespaceEnvVar)
	if len(namespaceEnvValue) == 0 {
		return nil, fmt.Errorf("%s was not provided as an environment variable", serviceRegistryNamespaceEnvVar)
	}
	return &KubernetesServiceRegistryServicer{
		client:    k8sClient,
		namespace: namespaceEnvValue,
	}, nil
}

// ListAllServices returns the service name of all services in the registry.
func (s *KubernetesServiceRegistryServicer) ListAllServices(ctx context.Context, req *protos.Void) (*protos.ListAllServicesResponse, error) {
	ret := &protos.ListAllServicesResponse{}
	orc8rListOption := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", partOfLabel, partOfOrc8rApp),
	}
	svcList, err := s.client.Services(s.namespace).List(orc8rListOption)
	if err != nil {
		return ret, err
	}
	for _, svc := range svcList.Items {
		ret.Services = append(ret.Services, s.parseServiceName(svc.Name))
	}
	return ret, nil
}

// FindServices returns all services in that have the provided label.
func (s *KubernetesServiceRegistryServicer) FindServices(ctx context.Context, req *protos.FindServicesRequest) (*protos.FindServicesResponse, error) {
	ret := &protos.FindServicesResponse{}
	orc8rListOption := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s,%s=true", partOfLabel, partOfOrc8rApp, req.GetLabel()),
	}
	svcList, err := s.client.Services(s.namespace).List(orc8rListOption)
	if err != nil {
		return ret, err
	}
	for _, svc := range svcList.Items {
		ret.Services = append(ret.Services, s.parseServiceName(svc.Name))
	}
	return ret, nil
}

// GetServiceAddress return the address of the gRPC server for the provided
// service.
func (s *KubernetesServiceRegistryServicer) GetServiceAddress(ctx context.Context, req *protos.GetServiceAddressRequest) (*protos.GetServiceAddressResponse, error) {
	if req == nil {
		return &protos.GetServiceAddressResponse{}, fmt.Errorf("GetServiceAddressRequest was nil")
	}
	serviceAddress, err := s.getAddressForPortName(req.GetService(), grpcPortName)
	if err != nil {
		return &protos.GetServiceAddressResponse{}, err
	}
	return &protos.GetServiceAddressResponse{
		Address: serviceAddress,
	}, nil
}

// GetHttpServerAddress returns the address of the HTTP server for the provided
// service.
func (s *KubernetesServiceRegistryServicer) GetHttpServerAddress(ctx context.Context, req *protos.GetHttpServerAddressRequest) (*protos.GetHttpServerAddressResponse, error) {
	if req == nil {
		return &protos.GetHttpServerAddressResponse{}, fmt.Errorf("GetHttpServerAddressRequest was nil")
	}
	httpServerAddress, err := s.getAddressForPortName(req.GetService(), httpPortName)
	if err != nil {
		return &protos.GetHttpServerAddressResponse{}, err
	}
	return &protos.GetHttpServerAddressResponse{
		Address: httpServerAddress,
	}, nil
}

// GetAnnotation returns the annotation value for the provided service and
// annotation.
func (s *KubernetesServiceRegistryServicer) GetAnnotation(ctx context.Context, req *protos.GetAnnotationRequest) (*protos.GetAnnotationResponse, error) {
	svc, err := s.getServiceForServiceName(req.GetService())
	if err != nil {
		return &protos.GetAnnotationResponse{}, err
	}
	for annotation, value := range svc.GetAnnotations() {
		if annotation == req.GetAnnotation() {
			return &protos.GetAnnotationResponse{
				AnnotationValue: value,
			}, nil
		}
	}
	return &protos.GetAnnotationResponse{}, fmt.Errorf("Annotation '%s' was not found for service '%s'", req.GetAnnotation(), req.GetService())
}

func (s *KubernetesServiceRegistryServicer) getAddressForPortName(service string, portName string) (string, error) {
	svc, err := s.getServiceForServiceName(service)
	if err != nil {
		return "", err
	}
	for _, port := range svc.Spec.Ports {
		if port.Name == portName {
			return fmt.Sprintf("%s:%d", svc.Name, port.Port), nil
		}
	}
	return "", fmt.Errorf("Could not find '%s' port for service '%s'", portName, service)
}

func (s *KubernetesServiceRegistryServicer) getServiceForServiceName(serviceName string) (*corev1types.Service, error) {
	// K8s services deployed via Helm have name with format
	// 'deploymentName-svcName'. Given that the mapping of module deployment
	// name to service is unknown to the registry, iterate through all
	// services and check the suffix
	orc8rListOption := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", partOfLabel, partOfOrc8rApp),
	}

	// Helm services don't allow for underscores, so convert to a dash
	formattedServiceName := strings.ReplaceAll(serviceName, "_", "-")
	svcList, err := s.client.Services(s.namespace).List(orc8rListOption)
	if err != nil {
		return nil, err
	}
	for _, svc := range svcList.Items {
		if strings.HasSuffix(svc.Name, formattedServiceName) {
			return &svc, nil
		}
	}
	return nil, fmt.Errorf("Could not find service '%s'", serviceName)
}

func (s *KubernetesServiceRegistryServicer) parseServiceName(svcName string) string {
	// K8s services deployed via Helm have name with format
	// 'deploymentName-svc-name'. Given that the deployment name for services
	// is unknown to other services running in the cluster, remove this prefix,
	// returning only the svc-name portion
	splitSvcName := strings.SplitAfterN(svcName, serviceNameDelimiter, 2)
	svcNameIndex := len(splitSvcName) - 1
	return splitSvcName[svcNameIndex]
}
