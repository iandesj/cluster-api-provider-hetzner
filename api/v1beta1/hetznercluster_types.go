/*
Copyright 2021 The Kubernetes Authors.

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
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// ClusterFinalizer allows ReconcileHetznerCluster to clean up HCloud
	// resources associated with HetznerCluster before removing it from the
	// apiserver.
	ClusterFinalizer = "hetznercluster.infrastructure.cluster.x-k8s.io"

	// LoadBalancerAlgorithmTypeRoundRobin default for the Kubernetes Api Server loadbalancer.
	LoadBalancerAlgorithmTypeRoundRobin = LoadBalancerAlgorithmType("round_robin")

	// LoadBalancerAlgorithmTypeLeastConnections default for Loadbalancer.
	LoadBalancerAlgorithmTypeLeastConnections = LoadBalancerAlgorithmType("least_connections")
)

// HetznerClusterSpec defines the desired state of HetznerCluster.
type HetznerClusterSpec struct {
	// HCloudNetworkSpec defines the Network for Hetzner Cloud. If left empty no private Network is configured.
	// +optional
	HCloudNetworkSpec HCloudNetworkSpec `json:"hcloudNetwork"`

	// ControlPlaneRegion consists of a list of HCloud Regions (fsn, nbg, hel). Because HCloud Networks
	// have a very low latency we could assume in some use-cases that a region is behaving like a zone
	// https://kubernetes.io/docs/reference/labels-annotations-taints/#topologykubernetesiozone
	// +kubebuilder:validation:MinItems=1
	ControlPlaneRegion []Region `json:"controlPlaneRegions"`

	// define cluster wide SSH keys. Valid values are a valid SSH key name, or a valid ID.
	SSHKey []SSHKeySpec `json:"sshKeys,omitempty"`
	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint *clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`

	// ControlPlaneLoadBalancer is optional configuration for customizing control plane behavior. Naming convention is from upstream cluster-api project.
	// +optional
	ControlPlaneLoadBalancer LoadBalancerSpec `json:"controlPlaneLoadBalancer,omitempty"`

	// +optional
	HCloudPlacementGroupSpec []HCloudPlacementGroupSpec `json:"hcloudPlacementGroups,omitempty"`

	// HetznerSecretRef is a reference to a token to be used when reconciling this cluster.
	// This is generated in the Security section under API TOKENS. Read & Write is necessary.
	HetznerSecretRef HetznerSecretRef `json:"hetznerSecretRef"`
}

// HetznerSecretRef defines all the name of the secret and the relevant keys needed to access Hetzner API.
type HetznerSecretRef struct {
	Name string              `json:"name"`
	Key  HetznerSecretKeyRef `json:"key"`
}

// HetznerSecretKeyRef defines the key name of the HetznerSecret.
type HetznerSecretKeyRef struct {
	HCloudToken string `json:"hcloudToken"`
}

// LoadBalancerSpec defines the desired state of the Control Plane Loadbalancer.
type LoadBalancerSpec struct {
	// +optional
	Name *string `json:"name,omitempty"`

	// Could be round_robin or least_connection. The default value is "round_robin".
	// +optional
	// +kubebuilder:validation:Enum=round_robin;least_connections
	// +kubebuilder:default=round_robin
	Algorithm LoadBalancerAlgorithmType `json:"algorithm,omitempty"`

	// Loadbalancer type
	// +optional
	// +kubebuilder:validation:Enum=lb11;lb21;lb31
	// +kubebuilder:default=lb11
	Type string `json:"type,omitempty"`

	// API Server port. It must be valid ports range (1-65535). If omitted, default value is 6443.
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=6443
	Port int `json:"port,omitempty"`

	// Defines how traffic will be routed from the Load Balancer to your target server.
	// +optional
	Targets []LoadBalancerTargetSpec `json:"extraTargets,omitempty"`

	Region Region `json:"region"`
}

// LoadBalancerStatus defines the obeserved state of the control plane loadbalancer.
type LoadBalancerStatus struct {
	ID         int    `json:"id,omitempty"`
	IPv4       string `json:"ipv4,omitempty"`
	IPv6       string `json:"ipv6,omitempty"`
	InternalIP string `json:"internalIP,omitempty"`
	Target     []int  `json:"targets,omitempty"`
	Protected  bool   `json:"protected,omitempty"`
}

// HetznerClusterStatus defines the observed state of HetznerCluster.
type HetznerClusterStatus struct {
	// +kubebuilder:default=false
	Ready bool `json:"ready"`

	// +optional
	Network *NetworkStatus `json:"networkStatus,omitempty"`

	ControlPlaneLoadBalancer *LoadBalancerStatus `json:"controlPlaneLoadBalancer,omitempty"`
	// +optional
	HCloudPlacementGroup []HCloudPlacementGroupStatus `json:"hcloudPlacementGroups,omitempty"`
	FailureDomains       clusterv1.FailureDomains     `json:"failureDomains,omitempty"`
	Conditions           clusterv1.Conditions         `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=hetznerclusters,scope=Namespaced,categories=cluster-api,shortName=capihc
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this HetznerCluster belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Cluster infrastructure is ready for Nodes"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".spec.controlPlaneEndpoint",description="API Endpoint",priority=1
// +k8s:defaulter-gen=true

// HetznerCluster is the Schema for the hetznercluster API.
type HetznerCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HetznerClusterSpec   `json:"spec,omitempty"`
	Status HetznerClusterStatus `json:"status,omitempty"`
}

// HetznerClusterList contains a list of HetznerCluster
// +kubebuilder:object:root=true
// +k8s:defaulter-gen=true
type HetznerClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HetznerCluster `json:"items"`
}

// GetConditions returns the observations of the operational state of the HetznerCluster resource.
func (r *HetznerCluster) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the HetznerCluster to the predescribed clusterv1.Conditions.
func (r *HetznerCluster) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&HetznerCluster{}, &HetznerClusterList{})
}
