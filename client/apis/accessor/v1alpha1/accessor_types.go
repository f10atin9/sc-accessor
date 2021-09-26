/*
Copyright 2021 f10atin9.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AccessorSpec defines the desired state of Accessor
type AccessorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// which storageClass is affected by accessor
	StorageClass string `json:"storage_class"`

	AllowedNamespace []string `json:"allowed_namespace" `
	AllowedWorkspace []string `json:"allowed_workspace" `
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Accessor is the Schema for the accessors API
type Accessor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AccessorSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// AccessorList contains a list of Accessor
type AccessorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Accessor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Accessor{}, &AccessorList{})
}
