/*
 Copyright 2018 The Operator-SDK Authors

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

/*
 Modifications:
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
	"strconv"
	"time"

	magmav1alpha1 "magma/cwf/k8s/cwf_operator/pkg/apis/magma/v1alpha1"
	"magma/cwf/k8s/cwf_operator/pkg/health_client"
	"magma/feg/cloud/go/protos"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_hacluster")

const (
	cwfAppSelectorKey      = "app.kubernetes.io/name"
	cwfAppSelectorValue    = "cwf"
	cwfInstanceSelectorKey = "app.kubernetes.io/instance"
	reconcilePeriod        = 15 * time.Second
	gatewayHealthService   = "health"
)

// Add creates a new HACluster Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileHACluster{
		client:               mgr.GetClient(),
		scheme:               mgr.GetScheme(),
		gatewayHealthService: gatewayHealthService,
		reconcilePeriod:      reconcilePeriod,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("hacluster-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource HACluster
	return c.Watch(&source.Kind{Type: &magmav1alpha1.HACluster{}}, &handler.EnqueueRequestForObject{})
}

// blank assignment to verify that ReconcileHACluster implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileHACluster{}

// ReconcileHACluster reconciles a HACluster object
type ReconcileHACluster struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client               client.Client
	scheme               *runtime.Scheme
	gatewayHealthService string
	reconcilePeriod      time.Duration
}

// Reconcile monitors gateway resources defined in the HACluster's spec
// and updates the HACluster's status, taking remediation steps if the
// the active gateway becomes unhealthy.
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileHACluster) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Cluster")

	hacluster := &magmav1alpha1.HACluster{}
	err := r.client.Get(context.TODO(), request.NamespacedName, hacluster)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("HACluster resource not found. Ignoring since object must be deleted.")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	if len(hacluster.Status.Active) == 0 {
		active := hacluster.Spec.GatewayResourceNames[0]
		newStatus := magmav1alpha1.HAClusterStatus{
			Active:           active,
			ActiveInitState:  magmav1alpha1.Uninitialized,
			StandbyInitState: magmav1alpha1.Uninitialized,
		}
		hacluster.Status = newStatus
		reqLogger.Info("No active is currently set. Setting active", "gateway", active)
		err = r.client.Status().Update(context.TODO(), hacluster)
		return reconcile.Result{RequeueAfter: r.reconcilePeriod}, err
	}
	if len(hacluster.Spec.GatewayResourceNames) == 1 {
		reqLogger.Info("Only 1 gateway resource configured. Not monitoring health")
		return reconcile.Result{RequeueAfter: r.reconcilePeriod}, nil
	}
	activeGateway := hacluster.Status.Active
	standbyGateway := r.getStandbyGatewayName(activeGateway, hacluster.Spec.GatewayResourceNames)

	activeHealth, err := r.getGatewayHealthStatus(activeGateway, request.Namespace)
	if err != nil {
		reqLogger.Error(err, "An error occurred while fetching active health status")
	} else {
		reqLogger.Info("Fetched active health status", "health", activeHealth.Health.String(), "message", activeHealth.HealthMessage)
	}
	standbyHealth, err := r.getGatewayHealthStatus(standbyGateway, request.Namespace)
	if err != nil {
		reqLogger.Error(err, "An error occurring while fetching standby health status")
	} else {
		reqLogger.Info("Fetched standby health status", "health", standbyHealth.Health.String(), "message", standbyHealth.HealthMessage)
	}

	failover := false
	var initErr error
	var updatedStatus magmav1alpha1.HAClusterStatus
	if activeHealth.Health == protos.HealthStatus_HEALTHY {
		updatedStatus, initErr = r.initCluster(activeGateway, standbyGateway, request.Namespace, hacluster.Status)
	} else if activeHealth.Health == protos.HealthStatus_UNHEALTHY && standbyHealth.Health == protos.HealthStatus_HEALTHY {
		failover = true
		reqLogger.Info("Promoting standby due to unhealthy active", "promoted active", standbyGateway)
		hacluster.Status.ActiveInitState = magmav1alpha1.Uninitialized
		hacluster.Status.StandbyInitState = magmav1alpha1.Uninitialized
		updatedStatus, initErr = r.initCluster(standbyGateway, activeGateway, request.Namespace, hacluster.Status)
	} else {
		reqLogger.Info("Both gateways are detected to be unhealthy")
		updatedStatus, initErr = r.initCluster(activeGateway, standbyGateway, request.Namespace, hacluster.Status)
	}

	if initErr != nil {
		reqLogger.Error(err, "failover occurred", strconv.FormatBool(failover))
	}
	hacluster.Status = updatedStatus
	updateErr := r.client.Status().Update(context.TODO(), hacluster)
	if updateErr != nil && failover {
		err = fmt.Errorf("Updating hacluster status with promoted active %s failed; %s", standbyGateway, updateErr)
		return reconcile.Result{RequeueAfter: r.reconcilePeriod}, updateErr
	} else if updateErr != nil {
		return reconcile.Result{RequeueAfter: r.reconcilePeriod}, updateErr
	}

	reqLogger.Info("Reconciled request")
	return reconcile.Result{RequeueAfter: r.reconcilePeriod}, initErr
}

// getGatewayHealthStatus fetches the health status for the provided gateway.
// This status is obtained from the gateway's local health service
func (r ReconcileHACluster) getGatewayHealthStatus(gateway string, namespace string) (*protos.HealthStatus, error) {
	// Get pod first to avoid sending an RPC that will timeout
	_, err := r.getPodNamespacedNameForGateway(gateway, namespace)
	if err != nil {
		return &protos.HealthStatus{
			Health:        protos.HealthStatus_UNHEALTHY,
			HealthMessage: fmt.Sprintf("could not find pod for gw: %s", gateway),
		}, err
	}
	svc, err := r.getHealthServiceAddressForResource(gateway, namespace)
	if err != nil {
		return &protos.HealthStatus{
			Health:        protos.HealthStatus_UNHEALTHY,
			HealthMessage: fmt.Sprintf("could not find svc endpoint for local health on gateway: %s", gateway),
		}, err
	}
	health, err := health_client.GetHealthStatus(svc)
	if err != nil {
		// GRPC doesn't return the proto if an error is returned.
		// To make the health logic simpler, create the proto
		health = &protos.HealthStatus{
			Health:        protos.HealthStatus_UNHEALTHY,
			HealthMessage: "error fetching health status",
		}
	}
	return health, err
}

func (r ReconcileHACluster) initCluster(active string, standby string, namespace string, status magmav1alpha1.HAClusterStatus) (magmav1alpha1.HAClusterStatus, error) {
	ret := magmav1alpha1.HAClusterStatus{
		Active:           active,
		ActiveInitState:  status.ActiveInitState,
		StandbyInitState: status.StandbyInitState,
	}
	var initActiveErr error
	var initStandbyErr error
	if status.ActiveInitState != magmav1alpha1.Initialized {
		ret.ActiveInitState, initActiveErr = r.initGatewayWithRole(active, namespace, true)
	}
	if status.StandbyInitState != magmav1alpha1.Initialized {
		ret.StandbyInitState, initStandbyErr = r.initGatewayWithRole(standby, namespace, false)
	}
	return ret, r.constructAggregateClusterError(initActiveErr, initStandbyErr)
}

func (r ReconcileHACluster) initGatewayWithRole(gateway string, namespace string, active bool) (magmav1alpha1.HAClusterInitState, error) {
	// Get pod first to avoid sending an RPC that will timeout
	_, err := r.getPodNamespacedNameForGateway(gateway, namespace)
	if err != nil {
		return magmav1alpha1.Uninitialized, err
	}
	svc, err := r.getHealthServiceAddressForResource(gateway, namespace)
	if err != nil {
		return magmav1alpha1.Uninitialized, err
	}
	if active {
		err = health_client.Enable(svc)
	} else {
		err = health_client.Disable(svc)
	}
	if err != nil {
		return magmav1alpha1.Uninitialized, err
	}
	return magmav1alpha1.Initialized, nil
}

func (r ReconcileHACluster) getPodNamespacedNameForGateway(gateway string, namespace string) (*types.NamespacedName, error) {
	listOption := []client.ListOption{
		client.InNamespace(namespace),
		client.MatchingLabels{
			cwfAppSelectorKey:      cwfAppSelectorValue,
			cwfInstanceSelectorKey: gateway,
		},
	}
	podlist := &corev1.PodList{}
	err := r.client.List(context.TODO(), podlist, listOption...)
	if err != nil || len(podlist.Items) == 0 {
		return nil, err
	}
	pod := podlist.Items[0]
	return &types.NamespacedName{
		Namespace: namespace,
		Name:      pod.Name,
	}, nil
}

func (r ReconcileHACluster) getHealthServiceAddressForResource(gateway string, namespace string) (string, error) {
	healthService := &corev1.Service{}
	serviceName := types.NamespacedName{
		Name:      fmt.Sprintf("%s-%s", gateway, r.gatewayHealthService),
		Namespace: namespace,
	}
	err := r.client.Get(context.TODO(), serviceName, healthService)
	if err != nil {
		return "", err
	}
	// Sanity check to ensure the operator only connects to a ClusterIP service
	if healthService.Spec.Type != corev1.ServiceTypeClusterIP {
		return "", fmt.Errorf("%s is not a ClusterIP service", healthService.Name)
	}
	return fmt.Sprintf("%s:%d", healthService.Spec.ClusterIP, healthService.Spec.Ports[0].Port), nil
}

func (r ReconcileHACluster) getStandbyGatewayName(active string, gateways []string) string {
	for _, gw := range gateways {
		if gw != active {
			return gw
		}
	}
	return ""
}

func (r ReconcileHACluster) constructAggregateClusterError(activeErr error, standbyErr error) error {
	if activeErr == nil && standbyErr == nil {
		return nil
	} else if activeErr != nil {
		return fmt.Errorf("initializing the active failed: %s", activeErr)
	} else if standbyErr != nil {
		return fmt.Errorf("initializing the standby failed: %s", standbyErr)
	}
	return fmt.Errorf("initializing both the active and standby failed; active error: %s, standby error: %s", activeErr, standbyErr)
}
