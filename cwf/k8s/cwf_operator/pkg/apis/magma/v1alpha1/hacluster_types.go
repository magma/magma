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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HAClusterInitState string
type CarrierWifiAccessGatewayHealthCondition string
type CarrierWifiAccessGatewayInitState string

// GatewayResource defines a gateway in the HACluster
type GatewayResource struct {
	HelmReleaseName string `json:"helmReleaseName"`
	GatewayID       string `json:"gatewayID"`
}

const (
	Initialized   HAClusterInitState = "Initialized"
	Uninitialized HAClusterInitState = "Uninitialized"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HAClusterSpec defines the desired state of HACluster
type HAClusterSpec struct {
	// GatewayResources denotes the list of all gateway resources in the
	// HACluster
	// +kubebuilder:validation:MaxItems=2
	// +kubebuilder:validation:MinItems=1
	GatewayResources []GatewayResource `json:"gatewayResources"`

	// HAPairID specifies the associated pair ID in the orchestrator with this
	// HACluster
	// +kubebuilder:validation:MinLength=1
	HAPairID string `json:"haPairID"`

	// MaxConsecutiveActiveErrors denotes the maximum number of errors the
	// HACluster's active can have fetching health status before a failover
	// occurs
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=5
	MaxConsecutiveActiveErrors int `json:"maxConsecutiveActiveErrors"`

	// Important: Run "make gen" to regenerate code after modifying this file
}

// HAClusterStatus defines the observed state of HACluster
type HAClusterStatus struct {
	// Active contains the resource name of the active gateway in the HACluster
	Active string `json:"active"`

	// ActiveInitState denotes the initialization state of the active in the
	// HACluster
	ActiveInitState HAClusterInitState `json:"activeInitState"`

	// StandbyInitState denotes the initialization state of the standby in
	// the HACluster
	StandbyInitState HAClusterInitState `json:"standbyInitState"`

	// ConsecutiveActiveErrors denotes the number of consecutive errors
	// that have occurred when the active has been called for health status
	ConsecutiveActiveErrors int `json:"consecutiveActiveErrors"`

	// ConsecutiveStandbyErrors denotes the number of consecutive errors
	// that have occurred when the standby has been called for health status
	ConsecutiveStandbyErrors int `json:"consecutiveStandbyErrors"`
	// Important: Run "make gen" to regenerate code after modifying this file
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HACluster is the Schema for the haclusters API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=haclusters,scope=Namespaced
type HACluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HAClusterSpec   `json:"spec,omitempty"`
	Status HAClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HAClusterList contains a list of HACluster
type HAClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HACluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HACluster{}, &HAClusterList{})
}
