//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// IBMDEV: Done (also set cluster scope)

// GroupResourceQuotaEnforcerSpec defines the desired state of GroupResourceQuotaEnforcer
type GroupResourceQuotaEnforcerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// ImageRegistry is the image repository
    ImageRegistry string `json:"imageRegistry"`
	// ImageRegistry is the optional secret to use when pulling from the ImageRegistry
	ImagePullSecret string `json:"imagePullSecret"`
	// InstanceNamespace is the namespace in which to place the namespaced resources (e.g. workloads) for the CR instance
	InstanceNamespace string `json:"instanceNamespace"`
	// TLS Certificate Secret
	CertSecret string `json:"certSecret"`
}

// GroupResourceQuotaEnforcerStatus defines the observed state of GroupResourceQuotaEnforcer
type GroupResourceQuotaEnforcerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Nodes are the names of the grq pods
    // +listType=set
    Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GroupResourceQuotaEnforcer is the Schema for the groupresourcequotaenforcers API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=groupresourcequotaenforcers,scope=Cluster
type GroupResourceQuotaEnforcer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GroupResourceQuotaEnforcerSpec   `json:"spec,omitempty"`
	Status GroupResourceQuotaEnforcerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GroupResourceQuotaEnforcerList contains a list of GroupResourceQuotaEnforcer
type GroupResourceQuotaEnforcerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GroupResourceQuotaEnforcer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GroupResourceQuotaEnforcer{}, &GroupResourceQuotaEnforcerList{})
}
