/*
Copyright 2020 The Kubernetes Authors.

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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ContainerSetSpec defines the desired state of ContainerSet
type ContainerSetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Replicas is an example field of ContainerSet. Edit ContainerSet_types.go to remove/update
	Replicas *int32 `json:"replicas,omitempty"`

	// image is the container image to run.  Image must have a tag.
	Image string `json:"image,omitempty"`
}

// ContainerSetStatus defines the observed state of ContainerSet
type ContainerSetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	HealthyReplicas *int32 `json:"healthyReplicas,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// ContainerSet is the Schema for the containersets API
type ContainerSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContainerSetSpec   `json:"spec,omitempty"`
	Status ContainerSetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ContainerSetList contains a list of ContainerSet
type ContainerSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ContainerSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ContainerSet{}, &ContainerSetList{})
}
