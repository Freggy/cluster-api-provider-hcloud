/*
Copyright 2021.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HCloudClusterSpec defines the desired state of HCloudCluster
type HCloudClusterSpec struct {

	// The region the cluster will be created in. To see which regions are
	// available, see here https://docs.hetzner.cloud/#locations-get-all-locations
	// +optional
	Region string `json:"region"`

	// The load balancer that will be used for load balancing traffic across the cluster nodes
	// + optional
	LoadBalancer LoadBalancer `json:"loadBalancer,omitempty"`

	// Endpoint where the control plane is reachable
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint,omitempty"`
}

// HCloudClusterStatus defines the observed state of HCloudCluster
type HCloudClusterStatus struct {
	// The load balancer created by the controller
	LoadBalancer LoadBalancer `json:"loadBalancer,omitempty"`

	// Whether the cluster is ready
	// +optional
	Ready bool `json:"ready"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// HCloudCluster is the Schema for the hcloudclusters API
type HCloudCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HCloudClusterSpec   `json:"spec,omitempty"`
	Status HCloudClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HCloudClusterList contains a list of HCloudCluster
type HCloudClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HCloudCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HCloudCluster{}, &HCloudClusterList{})
}
