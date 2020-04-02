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
	"math/rand"
	"strconv"
	"testing"

	magmav1alpha1 "magma/cwf/k8s/cwf_operator/pkg/apis/magma/v1alpha1"

	"github.com/go-logr/glogr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
	mockHealthClient := &mockHealthClient{}
	gateways := []string{gw}
	r := initTestReconciler(gateways, name, namespace, mockHealthClient)
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
	initializedCluster := &magmav1alpha1.HACluster{}
	err = r.client.Get(context.TODO(), req.NamespacedName, initializedCluster)
	assert.NoError(t, err)
	assert.Equal(t, gw, initializedCluster.Status.Active)

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
	healthClient := &mockHealthClient{}
	r := initTestReconciler(gateways, name, namespace, healthClient)
	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}
	// Create 2 arbitrary pods representing gateways
	gwPod := createPod(gw, namespace)
	gwPod2 := createPod(gw2, namespace)
	err := r.client.Create(context.Background(), gwPod2.DeepCopy())
	assert.NoError(t, err)
	err = r.client.Create(context.Background(), gwPod.DeepCopy())
	assert.NoError(t, err)

	// Initial reconcile request should initialize active
	_, err = r.Reconcile(req)
	assert.NoError(t, err)
	initializedCluster := &magmav1alpha1.HACluster{}
	err = r.client.Get(context.Background(), req.NamespacedName, initializedCluster)
	assert.NoError(t, err)
	assert.Equal(t, gw, initializedCluster.Status.Active)

	// Test happy path
	healthClient.On("GetHealthStatus", mock.Anything, gwPod.Name).Return(magmav1alpha1.Healthy, nil).Once()
	healthClient.On("GetHealthStatus", mock.Anything, gwPod2.Name).Return(magmav1alpha1.Healthy, nil).Once()
	healthClient.On("Enable", mock.Anything, gwPod.Name).Return(nil).Once()
	_, err = r.Reconcile(req)
	assert.NoError(t, err)
	healthClient.AssertExpectations(t)
	initializedCluster = &magmav1alpha1.HACluster{}
	err = r.client.Get(context.Background(), req.NamespacedName, initializedCluster)
	assert.NoError(t, err)
	assert.Equal(t, gw, initializedCluster.Status.Active)
	assert.Equal(t, magmav1alpha1.Initialized, initializedCluster.Status.InitState)

	// Test failover
	healthClient = &mockHealthClient{}
	r.healthClient = healthClient
	healthClient.On("GetHealthStatus", mock.Anything, gwPod.Name).Return(magmav1alpha1.Unhealthy, fmt.Errorf("err connecting")).Once()
	healthClient.On("GetHealthStatus", mock.Anything, gwPod2.Name).Return(magmav1alpha1.Healthy, nil).Once()
	healthClient.On("Enable", mock.Anything, gwPod2.Name).Return(nil).Once()
	_, err = r.Reconcile(req)
	assert.NoError(t, err)
	healthClient.AssertExpectations(t)
	initializedCluster = &magmav1alpha1.HACluster{}
	err = r.client.Get(context.Background(), req.NamespacedName, initializedCluster)
	assert.NoError(t, err)
	assert.Equal(t, gw2, initializedCluster.Status.Active)
	assert.Equal(t, magmav1alpha1.Initialized, initializedCluster.Status.InitState)

	// Test failover and init failure
	healthClient = &mockHealthClient{}
	r.healthClient = healthClient
	healthClient.On("GetHealthStatus", mock.Anything, gwPod2.Name).Return(magmav1alpha1.Unhealthy, nil).Once()
	healthClient.On("GetHealthStatus", mock.Anything, gwPod.Name).Return(magmav1alpha1.Healthy, nil).Once()
	healthClient.On("Enable", mock.Anything, gwPod.Name).Return(fmt.Errorf("session restart failed")).Once()

	_, err = r.Reconcile(req)
	assert.NoError(t, err)
	healthClient.AssertExpectations(t)
	initializedCluster = &magmav1alpha1.HACluster{}
	err = r.client.Get(context.Background(), req.NamespacedName, initializedCluster)
	assert.NoError(t, err)
	assert.Equal(t, gw, initializedCluster.Status.Active)
	assert.Equal(t, magmav1alpha1.Uninitialized, initializedCluster.Status.InitState)

	// Test both unhealthy, but active uninitialized
	healthClient = &mockHealthClient{}
	r.healthClient = healthClient
	healthClient.On("GetHealthStatus", mock.Anything, gwPod.Name).Return(magmav1alpha1.Unhealthy, nil).Once()
	healthClient.On("GetHealthStatus", mock.Anything, gwPod2.Name).Return(magmav1alpha1.Unhealthy, nil).Once()
	healthClient.On("Enable", mock.Anything, gwPod.Name).Return(nil).Once()
	_, err = r.Reconcile(req)
	assert.NoError(t, err)
	healthClient.AssertExpectations(t)
	initializedCluster = &magmav1alpha1.HACluster{}
	err = r.client.Get(context.Background(), req.NamespacedName, initializedCluster)
	assert.NoError(t, err)
	assert.Equal(t, gw, initializedCluster.Status.Active)
	assert.Equal(t, magmav1alpha1.Initialized, initializedCluster.Status.InitState)
}

func initTestReconciler(gateways []string, name string, namespace string, healthClient *mockHealthClient) *ReconcileHACluster {
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
		client:       fakeClient,
		healthClient: healthClient,
		scheme:       s,
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

type mockHealthClient struct {
	mock.Mock
}

func (m *mockHealthClient) GetHealthStatus(ctx context.Context, serviceAddr string) (magmav1alpha1.CarrierWifiAccessGatewayHealthCondition, error) {
	args := m.Called(ctx, serviceAddr)
	return args.Get(0).(magmav1alpha1.CarrierWifiAccessGatewayHealthCondition), args.Error(1)
}

func (m *mockHealthClient) Enable(ctx context.Context, serviceAddr string) error {
	args := m.Called(ctx, serviceAddr)
	return args.Error(0)
}
