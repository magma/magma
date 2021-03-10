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

package servicers_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/service_registry/servicers"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	corev1Interface "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	namespace          = "magma_namespace"
	defaultTestTimeout = 5 * time.Second
)

func TestK8sListAllServices(t *testing.T) {
	servicer, mockClient, start, done := setupTest(t)

	req := &protos.Void{}
	response, err := servicer.ListAllServices(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, []string{"service1", "service_2"}, response.GetServices())

	err = mockClient.Services(namespace).Delete("orc8r-service-2", &metav1.DeleteOptions{})
	assert.NoError(t, err)
	refreshCache(t, start, done)

	response, err = servicer.ListAllServices(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, []string{"service1"}, response.GetServices())
}

func TestK8sFindServices(t *testing.T) {
	servicer, mockClient, start, done := setupTest(t)

	req := &protos.FindServicesRequest{
		Label: "label1",
	}
	response, err := servicer.FindServices(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, []string{"service1"}, response.GetServices())

	req.Label = "label2"
	response, err = servicer.FindServices(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, []string{"service1", "service_2"}, response.GetServices())

	err = mockClient.Services(namespace).Delete("orc8r-service1", &metav1.DeleteOptions{})
	assert.NoError(t, err)
	refreshCache(t, start, done)

	response, err = servicer.FindServices(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, []string{"service_2"}, response.GetServices())
}

func TestK8sGetServiceAddress(t *testing.T) {
	servicer, mockClient, start, done := setupTest(t)

	req := &protos.GetServiceAddressRequest{
		Service: "service1",
	}
	response, err := servicer.GetServiceAddress(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, response.GetAddress(), fmt.Sprintf("orc8r-service1:%d", registry.GrpcServicePort))

	err = mockClient.Services(namespace).Delete("orc8r-service1", &metav1.DeleteOptions{})
	assert.NoError(t, err)
	refreshCache(t, start, done)

	_, err = servicer.GetServiceAddress(context.Background(), req)
	assert.Error(t, err)
}

func TestK8sGetHttpServerAddress(t *testing.T) {
	servicer, mockClient, start, done := setupTest(t)

	req := &protos.GetHttpServerAddressRequest{
		Service: "service1",
	}
	response, err := servicer.GetHttpServerAddress(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, response.GetAddress(), fmt.Sprintf("orc8r-service1:%d", registry.HttpServerPort))

	err = mockClient.Services(namespace).Delete("orc8r-service1", &metav1.DeleteOptions{})
	assert.NoError(t, err)
	refreshCache(t, start, done)

	_, err = servicer.GetHttpServerAddress(context.Background(), req)
	assert.Error(t, err)
}

func TestK8sGetAnnotation(t *testing.T) {
	servicer, mockClient, start, done := setupTest(t)

	req := &protos.GetAnnotationRequest{
		Service:    "service1",
		Annotation: "annotation2",
	}
	response, err := servicer.GetAnnotation(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, "bar,baz", response.GetAnnotationValue())

	err = mockClient.Services(namespace).Delete("orc8r-service1", &metav1.DeleteOptions{})
	assert.NoError(t, err)
	refreshCache(t, start, done)

	_, err = servicer.GetAnnotation(context.Background(), req)
	assert.Error(t, err)
}

func createK8sServices(t *testing.T, mockClient corev1Interface.CoreV1Interface) {
	svc1 := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      "orc8r-service1",
			Labels: map[string]string{
				orc8r.PartOfLabel: orc8r.PartOfOrc8rApp,
				"label1":          "true",
				"label2":          "true",
			},
			Annotations: map[string]string{
				"annotation1": "foo",
				"annotation2": "bar,baz",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: orc8r.GRPCPortName,
					Port: registry.GrpcServicePort,
				},
				{
					Name: orc8r.HTTPPortName,
					Port: registry.HttpServerPort,
				},
			},
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: "127.0.0.1",
		},
	}
	svc2 := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      "orc8r-service-2",
			Labels: map[string]string{
				orc8r.PartOfLabel: orc8r.PartOfOrc8rApp,
				"label2":          "true",
				"label3":          "true",
			},
			Annotations: map[string]string{
				"annotation3": "roo",
				"annotation4": "par,zaz",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: orc8r.GRPCPortName,
					Port: registry.GrpcServicePort,
				},
			},
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: "127.0.0.1",
		},
	}
	svc3 := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      fmt.Sprintf("%s-%s", "nonorc8r", "service3"),
		},
		Spec: corev1.ServiceSpec{
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: "127.0.0.1",
		},
	}
	_, err := mockClient.Services(namespace).Create(svc1)
	assert.NoError(t, err)
	_, err = mockClient.Services(namespace).Create(svc2)
	assert.NoError(t, err)
	_, err = mockClient.Services(namespace).Create(svc3)
	assert.NoError(t, err)
}

func setupTest(t *testing.T) (*servicers.KubernetesServiceRegistryServicer, corev1Interface.CoreV1Interface, chan interface{}, chan interface{}) {
	mockClient := fake.NewSimpleClientset().CoreV1()
	createK8sServices(t, mockClient)

	start := make(chan interface{})
	done := make(chan interface{})
	reporter := &servicers.Reporter{
		RefreshStart: func() { start <- nil },
		RefreshDone:  func() { done <- nil },
	}

	os.Setenv(servicers.ServiceRegistryNamespaceEnvVar, namespace)
	refreshCacheFrequency := "@every 100ms" // ideally would be shorter, but seems like 1s is the effective minimum
	servicer, err := servicers.NewKubernetesServiceRegistryServicer(mockClient, refreshCacheFrequency, reporter)
	assert.NoError(t, err)

	refreshCache(t, start, done)

	return servicer, mockClient, start, done
}

func refreshCache(t *testing.T, start, done chan interface{}) {
	recvCh(t, start)
	recvCh(t, done)
}

func recvCh(t *testing.T, ch chan interface{}) {
	select {
	case <-ch:
		return
	case <-time.After(defaultTestTimeout):
		t.Fatal("receive on hook channel timed out")
	}
}
