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

package magmacwfk8s

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"magma/feg/cloud/go/protos"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"magma/cwf/k8s/cwf-operator/controllers/magma.cwf.k8s/health_client"
	"magma/cwf/k8s/cwf-operator/controllers/magma.cwf.k8s/status_reporter"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	magmacwfk8sv1alpha1 "magma/cwf/k8s/cwf-operator/apis/magma.cwf.k8s/v1alpha1"
)

const (
	cwfAppSelectorKey      = "app.kubernetes.io/name"
	cwfAppSelectorValue    = "cwf"
	cwfInstanceSelectorKey = "app.kubernetes.io/instance"
	reconcilePeriod        = 15 * time.Second
	retryPeriod            = 5 * time.Second
	gatewayHealthService   = "health"
)

// NewReconciler returns a new reconcile.Reconciler
func NewReconciler(mgr manager.Manager) *HAClusterReconciler {
	hc := health_client.NewHealthClient()
	return &HAClusterReconciler{
		client:               mgr.GetClient(),
		scheme:               mgr.GetScheme(),
		healthClient:         hc,
		gatewayHealthService: gatewayHealthService,
		statusReporter:       status_reporter.NewStatusReporter(),
		reconcilePeriod:      reconcilePeriod,
	}
}

// HAClusterReconciler reconciles a HACluster object
type HAClusterReconciler struct {
	client               client.Client
	scheme               *runtime.Scheme
	healthClient         *health_client.HealthClient
	gatewayHealthService string
	statusReporter       *status_reporter.StatusReporter
	reconcilePeriod      time.Duration
}

//+kubebuilder:rbac:groups=magma.cwf.k8s.magma.cwf,resources=haclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=magma.cwf.k8s.magma.cwf,resources=haclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=magma.cwf.k8s.magma.cwf,resources=haclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HACluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *HAClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx, "Request.Namespace", req.Namespace, "Request.Name", req.Name)
	reqLogger.Info("Reconciling Cluster")

	hacluster := &magmacwfk8sv1alpha1.HACluster{}
	err := r.client.Get(context.TODO(), req.NamespacedName, hacluster)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("HACluster resource not found. Ignoring since object must be deleted.")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	if len(hacluster.Status.Active) == 0 {
		active := hacluster.Spec.GatewayResources[0].HelmReleaseName
		newStatus := magmacwfk8sv1alpha1.HAClusterStatus{
			Active:           active,
			ActiveInitState:  magmacwfk8sv1alpha1.Uninitialized,
			StandbyInitState: magmacwfk8sv1alpha1.Uninitialized,
		}
		hacluster.Status = newStatus
		reqLogger.Info("No active is currently set. Setting active", "gateway", active)
		err = r.client.Status().Update(context.TODO(), hacluster)
		return reconcile.Result{RequeueAfter: r.reconcilePeriod}, err
	}
	if len(hacluster.Spec.GatewayResources) == 1 {
		reqLogger.Info("Only 1 gateway resource configured. Not monitoring health")
		return reconcile.Result{RequeueAfter: r.reconcilePeriod}, nil
	}
	activeGateway := hacluster.Status.Active
	standbyGateway := r.getStandbyGatewayName(activeGateway, hacluster.Spec.GatewayResources)

	activeHealth, err := r.getGatewayHealthStatus(activeGateway, req.Namespace)
	if err != nil {
		reqLogger.Error(
			err,
			"An error occurred while fetching active health status",
			"consecutive error count",
			hacluster.Status.ConsecutiveActiveErrors,
			"max consecutive errors for failover",
			hacluster.Spec.MaxConsecutiveActiveErrors,
		)
		hacluster.Status.ConsecutiveActiveErrors++
		if hacluster.Status.ConsecutiveActiveErrors < hacluster.Spec.MaxConsecutiveActiveErrors {
			updateErr := r.client.Status().Update(context.TODO(), hacluster)
			if updateErr != nil {
				reqLogger.Error(updateErr, "An error occurred while updating consecutive error total for active")
			}
			return reconcile.Result{RequeueAfter: retryPeriod}, updateErr
		}
	} else {
		hacluster.Status.ConsecutiveActiveErrors = 0
		reqLogger.Info("Fetched active health status", "health", activeHealth.Health.String(), "message", activeHealth.HealthMessage)
	}
	standbyHealth, err := r.getGatewayHealthStatus(standbyGateway, req.Namespace)
	if err != nil {
		hacluster.Status.ConsecutiveStandbyErrors++
		reqLogger.Error(err, "An error occurring while fetching standby health status")
	} else {
		hacluster.Status.ConsecutiveStandbyErrors = 0
		reqLogger.Info("Fetched standby health status", "health", standbyHealth.Health.String(), "message", standbyHealth.HealthMessage)
	}

	failover := false
	var initErr error
	var updatedStatus magmacwfk8sv1alpha1.HAClusterStatus
	isActiveHealthy := activeHealth.Health == protos.HealthStatus_HEALTHY
	isStandbyHealthy := standbyHealth.Health == protos.HealthStatus_HEALTHY

	if isActiveHealthy {
		updatedStatus, initErr = r.initCluster(activeGateway, standbyGateway, req.Namespace, hacluster.Status, false)
	} else if isStandbyHealthy {
		reqLogger.Info("Promoting standby due to unhealthy active", "promoted active", standbyGateway)
		hacluster.Status.ActiveInitState = magmacwfk8sv1alpha1.Uninitialized
		hacluster.Status.StandbyInitState = magmacwfk8sv1alpha1.Uninitialized
		hacluster.Status.ConsecutiveActiveErrors = 0
		hacluster.Status.ConsecutiveStandbyErrors = 0
		updatedStatus, initErr = r.initCluster(standbyGateway, activeGateway, req.Namespace, hacluster.Status, true)
	} else {
		reqLogger.Info("Both gateways are detected to be unhealthy")
		// If both gateways are unhealthy, there is a chance that a restart
		// occurred, rendering the gateway(s) uninitialized. As a precaution,
		// update the status to uninitialized and re-init. If the gateways are
		// already initialized, this will be a no-op
		hacluster.Status.ActiveInitState = magmacwfk8sv1alpha1.Uninitialized
		hacluster.Status.StandbyInitState = magmacwfk8sv1alpha1.Uninitialized
		updatedStatus, initErr = r.initCluster(activeGateway, standbyGateway, req.Namespace, hacluster.Status, false)
	}

	if initErr != nil {
		reqLogger.Error(initErr, "did failover occur", strconv.FormatBool(failover))
	}
	if !reflect.DeepEqual(hacluster.Status, updatedStatus) {
		hacluster.Status = updatedStatus
		updateErr := r.client.Status().Update(context.TODO(), hacluster)
		if updateErr != nil {
			if failover {
				reqLogger.Error(updateErr, "Updating hacluster status with promoted active failed", "promoted active", standbyGateway)
				return reconcile.Result{}, updateErr
			}
			reqLogger.Error(updateErr, "Updating hacluster status failed")
			return reconcile.Result{}, updateErr
		}
	}
	if failover {
		go r.statusReporter.UpdateHAClusterStatus(hacluster.Status, hacluster.Spec, standbyHealth, activeHealth)
	} else {
		go r.statusReporter.UpdateHAClusterStatus(hacluster.Status, hacluster.Spec, activeHealth, standbyHealth)
	}

	reqLogger.Info("Reconciled request")
	return reconcile.Result{RequeueAfter: r.reconcilePeriod}, nil
}

func (r *HAClusterReconciler) getStandbyGatewayName(active string, gateways []magmacwfk8sv1alpha1.GatewayResource) string {
	for _, gw := range gateways {
		if gw.HelmReleaseName != active {
			return gw.HelmReleaseName
		}
	}
	return ""
}

// getGatewayHealthStatus fetches the health status for the provided gateway.
// This status is obtained from the gateway's local health service
func (r HAClusterReconciler) getGatewayHealthStatus(helmReleaseName string, namespace string) (*protos.HealthStatus, error) {
	// Get pod first to avoid sending an RPC that will timeout
	_, err := r.getPodNamespacedNameForGateway(helmReleaseName, namespace)
	if err != nil {
		return &protos.HealthStatus{
			Health:        protos.HealthStatus_UNHEALTHY,
			HealthMessage: fmt.Sprintf("could not find pod for gw release: %s", helmReleaseName),
		}, err
	}
	svc, port, err := r.getHealthServiceAddressForResource(helmReleaseName, namespace)
	if err != nil {
		return &protos.HealthStatus{
			Health:        protos.HealthStatus_UNHEALTHY,
			HealthMessage: fmt.Sprintf("could not find svc endpoint for local health on gateway: %s", helmReleaseName),
		}, err
	}
	health, err := r.healthClient.GetHealthStatus(svc, port)
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

func (r HAClusterReconciler) initCluster(active string, standby string, namespace string, status magmacwfk8sv1alpha1.HAClusterStatus, didFailover bool) (magmacwfk8sv1alpha1.HAClusterStatus, error) {
	ret := magmacwfk8sv1alpha1.HAClusterStatus{
		Active:                   active,
		ActiveInitState:          status.ActiveInitState,
		StandbyInitState:         status.StandbyInitState,
		ConsecutiveActiveErrors:  status.ConsecutiveActiveErrors,
		ConsecutiveStandbyErrors: status.ConsecutiveStandbyErrors,
	}
	var initActiveErr error
	var initStandbyErr error
	var recreateGateway bool

	// Init standby first to avoid potential dual-ownership of VIP
	if status.StandbyInitState != magmacwfk8sv1alpha1.Initialized {
		ret.StandbyInitState, recreateGateway, initStandbyErr = r.initGatewayWithRole(standby, namespace, false)
		// If we are unable to disable the standby after failover, recreate the
		// gateway to ensure the VIP is transferred properly to the active
		if didFailover && recreateGateway {
			r.recreateGateway(standby, namespace)
		}
	}
	if status.ActiveInitState != magmacwfk8sv1alpha1.Initialized {
		ret.ActiveInitState, _, initActiveErr = r.initGatewayWithRole(active, namespace, true)
	}

	return ret, r.constructAggregateClusterError(initActiveErr, initStandbyErr)
}

// initGatewayWithRole inits a gateway with active or standby role based off of
// the existing cluster initialization status. This function returns the
// updated cluster status, an error, and a boolean signaling if the returned
// error (if any) was an RPC error
func (r HAClusterReconciler) initGatewayWithRole(gateway string, namespace string, active bool) (magmacwfk8sv1alpha1.HAClusterInitState, bool, error) {
	// Get pod first to avoid sending an RPC that will timeout
	_, err := r.getPodNamespacedNameForGateway(gateway, namespace)
	if err != nil {
		return magmacwfk8sv1alpha1.Uninitialized, false, err
	}
	svc, port, err := r.getHealthServiceAddressForResource(gateway, namespace)
	if err != nil {
		return magmacwfk8sv1alpha1.Uninitialized, false, err
	}
	if active {
		err = r.healthClient.Enable(svc, port)
	} else {
		err = r.healthClient.Disable(svc, port)
	}
	if err != nil {
		return magmacwfk8sv1alpha1.Uninitialized, true, err
	}
	return magmacwfk8sv1alpha1.Initialized, false, nil
}

func (r HAClusterReconciler) recreateGateway(gateway string, namespace string) error {
	podName, err := r.getPodNamespacedNameForGateway(gateway, namespace)
	if err != nil {
		return err
	}
	pod := &corev1.Pod{}
	err = r.client.Get(context.TODO(), *podName, pod)
	if err != nil {
		return err
	}
	return r.client.Delete(context.TODO(), pod)
}

func (r HAClusterReconciler) getPodNamespacedNameForGateway(gateway string, namespace string) (*types.NamespacedName, error) {
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

func (r HAClusterReconciler) getHealthServiceAddressForResource(gateway string, namespace string) (string, int, error) {
	healthService := &corev1.Service{}
	serviceName := types.NamespacedName{
		Name:      fmt.Sprintf("%s-%s", gateway, r.gatewayHealthService),
		Namespace: namespace,
	}
	err := r.client.Get(context.TODO(), serviceName, healthService)
	if err != nil {
		return "", 0, err
	}
	// Sanity check to ensure the operator only connects to a ClusterIP service
	if healthService.Spec.Type != corev1.ServiceTypeClusterIP {
		return "", 0, fmt.Errorf("%s is not a ClusterIP service", healthService.Name)
	}
	return serviceName.Name, int(healthService.Spec.Ports[0].Port), nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HAClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&magmacwfk8sv1alpha1.HACluster{}).
		Complete(r)
}

func (r HAClusterReconciler) constructAggregateClusterError(activeErr error, standbyErr error) error {
	if activeErr == nil && standbyErr == nil {
		return nil
	} else if activeErr != nil {
		return fmt.Errorf("initializing the active failed: %s", activeErr)
	} else if standbyErr != nil {
		return fmt.Errorf("initializing the standby failed: %s", standbyErr)
	}
	return fmt.Errorf("initializing both the active and standby failed; active error: %s, standby error: %s", activeErr, standbyErr)
}
