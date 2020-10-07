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
	"testing"

	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	corev1Interface "k8s.io/client-go/kubernetes/typed/core/v1"
)

func TestK8sListAllServices(t *testing.T) {
	servicer, mockClient := setupTest(t)
	req := &protos.Void{}
	response, err := servicer.ListAllServices(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, []string{"service1", "service-2"}, response.GetServices())

	err = mockClient.Services(servicer.namespace).Delete("foo-service-2", &metav1.DeleteOptions{})
	assert.NoError(t, err)
	response, err = servicer.ListAllServices(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, []string{"service1"}, response.GetServices())
}

func TestK8sFindServices(t *testing.T) {
	servicer, mockClient := setupTest(t)
	req := &protos.FindServicesRequest{
		Label: "label1",
	}
	response, err := servicer.FindServices(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, []string{"service1"}, response.GetServices())

	req.Label = "label2"
	response, err = servicer.FindServices(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, []string{"service1", "service-2"}, response.GetServices())

	err = mockClient.Services(servicer.namespace).Delete("orc8r-service1", &metav1.DeleteOptions{})
	response, err = servicer.FindServices(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, []string{"service-2"}, response.GetServices())
}

func TestK8sGetServiceAddress(t *testing.T) {
	servicer, mockClient := setupTest(t)
	req := &protos.GetServiceAddressRequest{
		Service: "service1",
	}
	response, err := servicer.GetServiceAddress(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, response.GetAddress(), fmt.Sprintf("orc8r-service1:%d", registry.GrpcServicePort))

	mockClient.Services(servicer.namespace).Delete("orc8r-service1", &metav1.DeleteOptions{})
	_, err = servicer.GetServiceAddress(context.Background(), req)
	assert.Error(t, err)
}

func TestK8sGetHttpServerAddress(t *testing.T) {
	servicer, mockClient := setupTest(t)
	req := &protos.GetHttpServerAddressRequest{
		Service: "service1",
	}
	response, err := servicer.GetHttpServerAddress(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, response.GetAddress(), fmt.Sprintf("orc8r-service1:%d", registry.HttpServerPort))

	mockClient.Services(servicer.namespace).Delete("orc8r-service1", &metav1.DeleteOptions{})
	_, err = servicer.GetHttpServerAddress(context.Background(), req)
	assert.Error(t, err)
}

func TestK8sGetAnnotation(t *testing.T) {
	servicer, mockClient := setupTest(t)
	req := &protos.GetAnnotationRequest{
		Service:    "service1",
		Annotation: "annotation2",
	}
	response, err := servicer.GetAnnotation(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, "bar,baz", response.GetAnnotationValue())

	mockClient.Services(servicer.namespace).Delete("orc8r-service1", &metav1.DeleteOptions{})
	_, err = servicer.GetAnnotation(context.Background(), req)
	assert.Error(t, err)
}

func createK8sServices(t *testing.T, mockClient corev1Interface.CoreV1Interface, namespace string) {
	svc1 := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      "orc8r-service1",
			Labels: map[string]string{
				partOfLabel: partOfOrc8rApp,
				"label1":    "true",
				"label2":    "true",
			},
			Annotations: map[string]string{
				"annotation1": "foo",
				"annotation2": "bar,baz",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: grpcPortName,
					Port: registry.GrpcServicePort,
				},
				{
					Name: httpPortName,
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
			Name:      "foo-service-2",
			Labels: map[string]string{
				partOfLabel: partOfOrc8rApp,
				"label2":    "true",
				"label3":    "true",
			},
			Annotations: map[string]string{
				"annotation3": "roo",
				"annotation4": "par,zaz",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: grpcPortName,
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

func setupTest(t *testing.T) (*KubernetesServiceRegistryServicer, corev1Interface.CoreV1Interface) {
	mockClient := fake.NewSimpleClientset().CoreV1()
	os.Setenv(serviceRegistryNamespaceEnvVar, "magma")

	servicer, err := NewKubernetesServiceRegistryServicer(mockClient)
	assert.NoError(t, err)
	createK8sServices(t, mockClient, "magma")
	return servicer, mockClient
}
