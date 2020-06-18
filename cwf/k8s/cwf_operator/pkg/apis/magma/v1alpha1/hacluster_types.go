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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HAClusterInitState string
type CarrierWifiAccessGatewayHealthCondition string
type CarrierWifiAccessGatewayInitState string

const (
	Healthy   CarrierWifiAccessGatewayHealthCondition = "Healthy"
	Unhealthy CarrierWifiAccessGatewayHealthCondition = "Unhealthy"

	Initialized   HAClusterInitState = "Initialized"
	Uninitialized HAClusterInitState = "Uninitialized"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HAClusterSpec defines the desired state of HACluster
type HAClusterSpec struct {
	// GatewayResourceNames denotes the list of all gateway resource names in the HACluster
	// +kubebuilder:validation:MaxItems=2
	// +kubebuilder:validation:MinItems=1
	GatewayResourceNames []string `json:"gatewayResourceNames"`
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
