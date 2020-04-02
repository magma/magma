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
	"time"

	magmav1alpha1 "magma/cwf/k8s/cwf_operator/pkg/apis/magma/v1alpha1"
	"magma/cwf/k8s/cwf_operator/pkg/health_client"

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
)

// Add creates a new HACluster Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	client := &health_client.HealthClient{}
	return add(mgr, newReconciler(mgr, client))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, healthClient health_client.HealthServiceClient) reconcile.Reconciler {
	return &ReconcileHACluster{
		client:          mgr.GetClient(),
		scheme:          mgr.GetScheme(),
		healthClient:    healthClient,
		reconcilePeriod: reconcilePeriod,
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
	client          client.Client
	healthClient    health_client.HealthServiceClient
	scheme          *runtime.Scheme
	reconcilePeriod time.Duration
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
		hacluster.Status.Active = active
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

	activeHealth, err := r.getGatewayHealthStatus(activeGateway, request.Namespace, hacluster)
	if err != nil {
		reqLogger.Error(err, "An error occurred while fetching active health status")
	}
	standbyHealth, err := r.getGatewayHealthStatus(standbyGateway, request.Namespace, hacluster)
	if err != nil {
		reqLogger.Error(err, "An error occurring while fetching standby health status")
	}

	reqLogger.Info("Fetched health status", "activeGateway", activeGateway, "health", activeHealth,
		"standbyGateway", standbyGateway, "health", standbyHealth)

	failover := false
	var updatedStatus magmav1alpha1.HAClusterStatus
	if activeHealth == magmav1alpha1.Healthy {
		updatedStatus, err = r.initActive(activeGateway, request.Namespace, hacluster.Status.InitState)
	} else if activeHealth == magmav1alpha1.Unhealthy && standbyHealth == magmav1alpha1.Healthy {
		failover = true
		reqLogger.Info("Promoting standby due to unhealthy active", "promoted active", standbyGateway)
		updatedStatus, err = r.initActive(standbyGateway, request.Namespace, magmav1alpha1.Uninitialized)
	} else {
		reqLogger.Info("Both gateways are detected to be unhealthy")
		updatedStatus, err = r.initActive(activeGateway, request.Namespace, hacluster.Status.InitState)
	}

	if err != nil && failover {
		reqLogger.Error(err, "Could not initialize promoted active")
	} else if err != nil {
		reqLogger.Error(err, "Could not initialize active")
	}
	hacluster.Status = updatedStatus
	err = r.client.Status().Update(context.TODO(), hacluster)

	if err != nil && failover {
		err = fmt.Errorf("Updating hacluster status with promoted active %s failed; %s", standbyGateway, err)
		return reconcile.Result{RequeueAfter: r.reconcilePeriod}, err
	} else if err != nil {
		return reconcile.Result{RequeueAfter: r.reconcilePeriod}, err
	}

	reqLogger.Info("Reconciled request")
	return reconcile.Result{RequeueAfter: r.reconcilePeriod}, nil
}

// getGatewayHealthStatus fetches the health status for the provided gateway.
// This status is obtained from the gateway's local health service
func (r ReconcileHACluster) getGatewayHealthStatus(gateway string, namespace string, hacluster *magmav1alpha1.HACluster) (magmav1alpha1.CarrierWifiAccessGatewayHealthCondition, error) {
	pod, err := r.getPodNamespacedNameForGateway(gateway, namespace)
	if err != nil {
		return magmav1alpha1.Unhealthy, err
	}
	// TODO: Fetch Kubernetes service:port and use as service addr
	return r.healthClient.GetHealthStatus(context.Background(), pod.Name)
}

func (r ReconcileHACluster) initActive(active string, namespace string, state magmav1alpha1.HAClusterInitState) (magmav1alpha1.HAClusterStatus, error) {
	ret := magmav1alpha1.HAClusterStatus{
		Active:    active,
		InitState: state,
	}
	if state == magmav1alpha1.Initialized {
		return ret, nil
	}
	pod, err := r.getPodNamespacedNameForGateway(active, namespace)
	if err != nil {
		return ret, err
	}
	err = r.healthClient.Enable(context.Background(), pod.Name)
	if err != nil {
		return ret, err
	}
	ret.InitState = magmav1alpha1.Initialized
	return ret, nil
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

func (r ReconcileHACluster) getStandbyGatewayName(active string, gateways []string) string {
	for _, gw := range gateways {
		if gw != active {
			return gw
		}
	}
	return ""
}
