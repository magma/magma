/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package hacluster

import (
	"context"
	"fmt"

	//"fmt"
	"math/rand"
	"net"
	"strconv"
	"testing"

	magmav1alpha1 "magma/cwf/k8s/cwf_operator/pkg/apis/magma/v1alpha1"
	"magma/feg/cloud/go/protos"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/go-logr/glogr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	healthSvcName = "testhealth"
)

// TestHAClusterControllerSingleGateway runs ReconcileHACluster.Reconcile() against a
// fake client that tracks an HACluster object with a single gateway configured
func TestHAClusterControllerSingleGateway(t *testing.T) {
	logf.SetLogger(glogr.New())
	var (
		name      = "cwf-operator"
		namespace = "test"
		gw        = "test-gw"
	)
	gateways := []string{gw}
	r := initTestReconciler(gateways, name, namespace)
	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}

	// Initial reconcile request should initialize active
	_, err := r.Reconcile(req)
	assert.NoError(t, err)
	mockCluster := &magmav1alpha1.HACluster{}
	err = r.client.Get(context.TODO(), req.NamespacedName, mockCluster)
	assert.NoError(t, err)
	assert.Equal(t, gw, mockCluster.Status.Active)

	// No monitoring action is taken with only one gateway configured
	_, err = r.Reconcile(req)
	assert.NoError(t, err)
}

// TestHAClusterControllerTwoGateways runs ReconcileHACluster.Reconcile() against a
// fake client that tracks an HACluster object with a two gateways configured
func TestHAClusterControllerTwoGateways(t *testing.T) {
	logf.SetLogger(glogr.New())
	var (
		name      = "cwf-operator"
		namespace = "test"
		gw        = "test-gw1"
		gw2       = "test-gw2"
	)
	gateways := []string{gw, gw2}
	mockServicer := &mockHealthServicer{}
	mockServicer2 := &mockHealthServicer{}
	r := initTestReconciler(gateways, name, namespace)
	addr := runMockService(t, mockServicer)
	addr2 := runMockService(t, mockServicer2)
	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}

	healthyStatus := &protos.HealthStatus{
		Health:        protos.HealthStatus_HEALTHY,
		HealthMessage: "Healthy",
	}
	unhealthyStatus := &protos.HealthStatus{
		Health:        protos.HealthStatus_UNHEALTHY,
		HealthMessage: "Unhealthy",
	}
	void := &orcprotos.Void{}

	// Create arbitrary pod and svc representing gatways
	gwPod := createPod(gw, namespace)
	gwSvc := createSvc(gw, namespace, int32(addr.Port))
	gwPod2 := createPod(gw2, namespace)
	gwSvc2 := createSvc(gw2, namespace, int32(addr2.Port))

	// Create gw1 resources
	err := r.client.Create(context.Background(), gwPod.DeepCopy())
	assert.NoError(t, err)
	err = r.client.Create(context.Background(), gwSvc.DeepCopy())
	assert.NoError(t, err)

	// Initial reconcile request should initialize active
	_, err = r.Reconcile(req)
	assert.NoError(t, err)
	mockCluster := &magmav1alpha1.HACluster{}
	err = r.client.Get(context.Background(), req.NamespacedName, mockCluster)
	assert.NoError(t, err)
	assert.Equal(t, gw, mockCluster.Status.Active)
	assert.Equal(t, magmav1alpha1.Uninitialized, mockCluster.Status.ActiveInitState)
	assert.Equal(t, magmav1alpha1.Uninitialized, mockCluster.Status.StandbyInitState)

	// Test proper error handling if a configured gateway doesn't actually exist
	mockServicer.On("GetHealthStatus", mock.Anything, mock.Anything).Return(healthyStatus, nil).Once()
	mockServicer.On("Enable", mock.Anything, mock.Anything).Return(void, nil).Once()
	_, err = r.Reconcile(req)
	assert.Error(t, err)
	mockServicer.AssertExpectations(t)

	// Create gw2 resource
	err = r.client.Create(context.Background(), gwPod2.DeepCopy())
	assert.NoError(t, err)
	err = r.client.Create(context.Background(), gwSvc2.DeepCopy())
	assert.NoError(t, err)

	// Test happy path - active already initialized
	mockServicer.On("GetHealthStatus", mock.Anything, mock.Anything).Return(healthyStatus, nil).Once()
	mockServicer2.On("GetHealthStatus", mock.Anything, mock.Anything).Return(healthyStatus, nil).Once()
	mockServicer2.On("Disable", mock.Anything, mock.Anything).Return(void, nil).Once()
	_, err = r.Reconcile(req)
	assert.NoError(t, err)
	mockServicer.AssertExpectations(t)
	mockServicer2.AssertExpectations(t)
	mockcluster := &magmav1alpha1.HACluster{}
	err = r.client.Get(context.Background(), req.NamespacedName, mockcluster)
	assert.NoError(t, err)
	assert.Equal(t, gw, mockcluster.Status.Active)
	assert.Equal(t, magmav1alpha1.Initialized, mockcluster.Status.ActiveInitState)
	assert.Equal(t, magmav1alpha1.Initialized, mockcluster.Status.StandbyInitState)

	// Test successful failover
	mockServicer.On("GetHealthStatus", mock.Anything, mock.Anything).Return(unhealthyStatus, fmt.Errorf("err connecting")).Once()
	mockServicer2.On("GetHealthStatus", mock.Anything, mock.Anything).Return(healthyStatus, nil).Once()
	mockServicer2.On("Enable", mock.Anything, mock.Anything).Return(void, nil).Once()
	mockServicer.On("Disable", mock.Anything, mock.Anything).Return(void, nil).Once()
	_, err = r.Reconcile(req)
	assert.NoError(t, err)
	mockServicer.AssertExpectations(t)
	mockServicer2.AssertExpectations(t)
	mockCluster = &magmav1alpha1.HACluster{}
	err = r.client.Get(context.Background(), req.NamespacedName, mockCluster)
	assert.NoError(t, err)
	assert.Equal(t, gw2, mockCluster.Status.Active)
	assert.Equal(t, magmav1alpha1.Initialized, mockCluster.Status.ActiveInitState)
	assert.Equal(t, magmav1alpha1.Initialized, mockCluster.Status.StandbyInitState)

	// Test failover and active init failure
	mockServicer2.On("GetHealthStatus", mock.Anything, mock.Anything).Return(unhealthyStatus, nil).Once()
	mockServicer.On("GetHealthStatus", mock.Anything, mock.Anything).Return(healthyStatus, nil).Once()
	mockServicer.On("Enable", mock.Anything, mock.Anything).Return(void, fmt.Errorf("session restart failed")).Once()
	mockServicer2.On("Disable", mock.Anything, mock.Anything).Return(void, nil).Once()
	_, err = r.Reconcile(req)
	assert.Error(t, err)
	mockServicer.AssertExpectations(t)
	mockServicer2.AssertExpectations(t)
	mockCluster = &magmav1alpha1.HACluster{}
	err = r.client.Get(context.Background(), req.NamespacedName, mockCluster)
	assert.NoError(t, err)
	assert.Equal(t, gw, mockCluster.Status.Active)
	assert.Equal(t, magmav1alpha1.Uninitialized, mockCluster.Status.ActiveInitState)
	assert.Equal(t, magmav1alpha1.Initialized, mockCluster.Status.StandbyInitState)

	// Test both unhealthy, but active uninitialized
	mockServicer.On("GetHealthStatus", mock.Anything, mock.Anything).Return(unhealthyStatus, nil).Once()
	mockServicer2.On("GetHealthStatus", mock.Anything, mock.Anything).Return(unhealthyStatus, nil).Once()
	mockServicer.On("Enable", mock.Anything, mock.Anything).Return(void, nil).Once()
	_, err = r.Reconcile(req)
	assert.NoError(t, err)
	mockServicer.AssertExpectations(t)
	mockServicer2.AssertExpectations(t)
	mockCluster = &magmav1alpha1.HACluster{}
	err = r.client.Get(context.Background(), req.NamespacedName, mockCluster)
	assert.NoError(t, err)
	assert.Equal(t, gw, mockCluster.Status.Active)
	assert.Equal(t, magmav1alpha1.Initialized, mockCluster.Status.ActiveInitState)
	assert.Equal(t, magmav1alpha1.Initialized, mockCluster.Status.StandbyInitState)
}

func initTestReconciler(gateways []string, name string, namespace string) *ReconcileHACluster {
	// An haCluster resource with metadata and spec.
	haCluster := &magmav1alpha1.HACluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: magmav1alpha1.HAClusterSpec{
			GatewayResourceNames: gateways,
		},
	}
	// Objects to track in the fake client.
	objs := []runtime.Object{
		haCluster,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(magmav1alpha1.SchemeGroupVersion, haCluster)
	// Create a fake client to mock API calls.
	fakeClient := fake.NewFakeClientWithScheme(s, objs...)

	// Create a ReconcileHaCluster object with the scheme, fake client, and
	// fake v1 client.
	return &ReconcileHACluster{
		client:               fakeClient,
		scheme:               s,
		gatewayHealthService: healthSvcName,
	}
}

func createPod(instanceName string, namespace string) *corev1.Pod {
	return &corev1.Pod{ObjectMeta: metav1.ObjectMeta{
		Namespace: namespace,
		Labels: map[string]string{
			cwfAppSelectorKey:      cwfAppSelectorValue,
			cwfInstanceSelectorKey: instanceName,
		},
		Name: instanceName + "-" + strconv.Itoa(rand.Int()),
	}, Spec: corev1.PodSpec{Containers: []corev1.Container{{Image: "blah"}}}}
}

func createSvc(instanceName string, namespace string, port int32) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      instanceName + "-" + healthSvcName,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Port: port,
				},
			},
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: "127.0.0.1",
		},
	}
}

type mockHealthServicer struct {
	mock.Mock
}

func (m *mockHealthServicer) GetHealthStatus(ctx context.Context, void *orcprotos.Void) (*protos.HealthStatus, error) {
	args := m.Called(ctx, void)
	return args.Get(0).(*protos.HealthStatus), args.Error(1)
}

func (m *mockHealthServicer) Disable(ctx context.Context, req *protos.DisableMessage) (*orcprotos.Void, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*orcprotos.Void), args.Error(1)
}

func (m *mockHealthServicer) Enable(ctx context.Context, void *orcprotos.Void) (*orcprotos.Void, error) {
	args := m.Called(ctx, void)
	return args.Get(0).(*orcprotos.Void), args.Error(1)
}

func runMockService(t *testing.T, servicer *mockHealthServicer) *net.TCPAddr {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create listener: %s", err)
	}
	addr, err := net.ResolveTCPAddr("tcp", lis.Addr().String())
	if err != nil {
		t.Fatalf("failed to resolve TCP address: %s", err)
	}
	opts := []grpc.ServerOption{}
	srv := grpc.NewServer(opts...)
	protos.RegisterServiceHealthServer(srv, servicer)
	go runServer(t, srv, lis)
	return addr
}

func runServer(t *testing.T, srv *grpc.Server, lis net.Listener) {
	err := srv.Serve(lis)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
}
