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
	"sync"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"golang.org/x/net/context"
	corev1types "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	ServiceRegistryNamespaceEnvVar = "SERVICE_REGISTRY_NAMESPACE"

	orc8rServiceNamePrefix = "orc8r-"
)

type KubernetesServiceRegistryServicer struct {
	sync.RWMutex
	client    corev1.CoreV1Interface
	namespace string
	cache     []corev1types.Service
	reporter  *Reporter
}

// NewKubernetesServiceRegistryServicer creates a new service registry servicer
// that is backed by Kubernetes.
//
// Takes an argument for how frequently to refresh the local cache of tracked
// services.
// Ref: https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format
func NewKubernetesServiceRegistryServicer(k8sClient corev1.CoreV1Interface, refreshCacheFrequency string, reporter *Reporter) (*KubernetesServiceRegistryServicer, error) {
	namespaceEnvValue := os.Getenv(ServiceRegistryNamespaceEnvVar)
	if len(namespaceEnvValue) == 0 {
		return nil, fmt.Errorf("%s was not provided as an environment variable", ServiceRegistryNamespaceEnvVar)
	}

	k := &KubernetesServiceRegistryServicer{client: k8sClient, namespace: namespaceEnvValue, reporter: reporter}

	c := cron.New()
	_, err := c.AddFunc(refreshCacheFrequency, k.refreshServicesCache)
	if err != nil {
		return nil, err
	}
	c.Start()

	// Seed registry with initial values
	go k.refreshServicesCache()

	return k, nil
}

// ListAllServices returns the service name of all services in the registry.
func (k *KubernetesServiceRegistryServicer) ListAllServices(ctx context.Context, req *protos.Void) (*protos.ListAllServicesResponse, error) {
	k.RLock()
	defer k.RUnlock()
	ret := &protos.ListAllServicesResponse{}
	for _, svc := range k.cache {
		formattedName := convertK8sServiceNameToMagmaServiceName(svc.Name)
		ret.Services = append(ret.Services, formattedName)
	}
	return ret, nil
}

// FindServices returns all services in that have the provided label.
func (k *KubernetesServiceRegistryServicer) FindServices(ctx context.Context, req *protos.FindServicesRequest) (*protos.FindServicesResponse, error) {
	k.RLock()
	defer k.RUnlock()
	ret := &protos.FindServicesResponse{}
	for _, svc := range k.cache {
		if hasLabel(svc, req.Label) {
			formattedName := convertK8sServiceNameToMagmaServiceName(svc.Name)
			ret.Services = append(ret.Services, formattedName)
		}
	}
	return ret, nil
}

// GetServiceAddress return the address of the gRPC server for the provided
// service.
func (k *KubernetesServiceRegistryServicer) GetServiceAddress(ctx context.Context, req *protos.GetServiceAddressRequest) (*protos.GetServiceAddressResponse, error) {
	k.RLock()
	defer k.RUnlock()
	if req == nil {
		return &protos.GetServiceAddressResponse{}, fmt.Errorf("GetServiceAddressRequest was nil")
	}
	serviceAddress, err := k.getAddressForPortName(req.GetService(), orc8r.GRPCPortName)
	if err != nil {
		return &protos.GetServiceAddressResponse{}, err
	}
	return &protos.GetServiceAddressResponse{
		Address: serviceAddress,
	}, nil
}

// GetHttpServerAddress returns the address of the HTTP server for the provided
// service.
func (k *KubernetesServiceRegistryServicer) GetHttpServerAddress(ctx context.Context, req *protos.GetHttpServerAddressRequest) (*protos.GetHttpServerAddressResponse, error) {
	k.RLock()
	defer k.RUnlock()
	if req == nil {
		return &protos.GetHttpServerAddressResponse{}, fmt.Errorf("GetHttpServerAddressRequest was nil")
	}
	httpServerAddress, err := k.getAddressForPortName(req.GetService(), orc8r.HTTPPortName)
	if err != nil {
		return &protos.GetHttpServerAddressResponse{}, err
	}
	return &protos.GetHttpServerAddressResponse{
		Address: httpServerAddress,
	}, nil
}

// GetAnnotation returns the annotation value for the provided service and
// annotation.
func (k *KubernetesServiceRegistryServicer) GetAnnotation(ctx context.Context, req *protos.GetAnnotationRequest) (*protos.GetAnnotationResponse, error) {
	k.RLock()
	defer k.RUnlock()
	svc, err := k.getServiceForServiceName(req.GetService())
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

func (k *KubernetesServiceRegistryServicer) getAddressForPortName(service string, portName string) (string, error) {
	svc, err := k.getServiceForServiceName(service)
	if err != nil {
		return "", err
	}
	for _, port := range svc.Spec.Ports {
		if port.Name == portName {
			return fmt.Sprintf("%s:%d", svc.Name, port.Port), nil
		}
	}
	return "", fmt.Errorf("could not find '%s' port for service '%s'", portName, service)
}

func (k *KubernetesServiceRegistryServicer) getServiceForServiceName(serviceName string) (*corev1types.Service, error) {
	formattedSvcName := convertMagmaServiceNameToK8sServiceName(serviceName)
	for _, svc := range k.cache {
		if svc.Name == formattedSvcName {
			return &svc, nil
		}
	}
	return nil, fmt.Errorf("could not find service '%s'", serviceName)
}

func (k *KubernetesServiceRegistryServicer) refreshServicesCache() {
	if k.reporter != nil {
		k.reporter.RefreshStart()
	}

	opts := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", orc8r.PartOfLabel, orc8r.PartOfOrc8rApp),
	}
	services, err := k.client.Services(k.namespace).List(opts)
	if err != nil {
		// Log error and leave previous cache intact
		err = errors.Wrap(err, "refresh service registry cache of K8s services")
		glog.Error(err)
		return
	}
	for _, s := range services.Items {
		s.Name = convertK8sServiceNameToMagmaServiceName(s.Name)
	}

	k.Lock()
	k.cache = services.Items
	k.Unlock()

	glog.V(1).Infof("Refreshed service registry cache. Found %d services.", len(k.cache))
	if k.reporter != nil {
		k.reporter.RefreshDone()
	}
}

// Orc8r helm services are formatted as orc8r-<svc-name>. Magma convention is
// to use underscores in service names, so remove prefix and convert any
// hyphens in the k8s service name.
func convertK8sServiceNameToMagmaServiceName(serviceName string) string {
	trimmedSvcName := strings.TrimPrefix(serviceName, orc8rServiceNamePrefix)
	return strings.ReplaceAll(trimmedSvcName, "-", "_")
}

// Orc8r helm services are formatted as orc8r-<svc-name>. Magma convention is
// to use underscores in service names, so add prefix and convert any
// underscores to hyphens
func convertMagmaServiceNameToK8sServiceName(serviceName string) string {
	k8sSvcNameSuffix := strings.ReplaceAll(serviceName, "_", "-")
	return fmt.Sprintf("%s%s", orc8rServiceNamePrefix, k8sSvcNameSuffix)
}

// hasLabel returns true if the service has the passed label and the label's
// value is "true".
func hasLabel(service corev1types.Service, label string) bool {
	for l, v := range service.ObjectMeta.Labels {
		if l == label && v == "true" {
			return true
		}
	}
	return false
}

// Reporter reports service registry events.
type Reporter struct {
	RefreshStart func()
	RefreshDone  func()
}
