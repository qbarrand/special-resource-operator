/*


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

package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	helmerv1beta1 "github.com/openshift-psap/special-resource-operator/pkg/helmer/api/v1beta1"
)

// SpecialResourceSpec describes the desired state of the resource, such as the recipe to be used and a selector
// on which nodes it should be installed.
// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
// +kubebuilder:validation:Required
type SpecialResourceSpec struct {
	// Chart describes the Helm chart that needs to be installed.
	// +kubebuilder:validation:Required
	Chart helmerv1beta1.HelmChart `json:"chart"`

	// Namespace describes in which namespace the chart will be installed.
	// +kubebuilder:validation:Required
	Namespace string `json:"namespace"`

	// If Debug is true, additional debugging output will be written in the log file.
	// +kubebuilder:validation:Optional
	Debug bool `json:"debug"`

	// Set are Helm hierarchical values for this chart installation.
	// +kubebuilder:validation:Optional
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:EmbeddedResource
	Set unstructured.Unstructured `json:"set,omitempty"`

	// NodeSelector is used to determine on which nodes the soiftware stack should be installed.
	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Dependencies is a list of dependencies required by this SpecialReosurce.
	// +kubebuilder:validation:Optional
	Dependencies []SpecialResourceDependency `json:"dependencies,omitempty"`
}

// SpecialResourceDependency is a Helm chart the SpecialResource depends on.
type SpecialResourceDependency struct {
	helmerv1beta1.HelmChart `json:"chart,omitempty"`

	// Set are Helm hierarchical values for this chart installation.
	// +kubebuilder:validation:Optional
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:EmbeddedResource
	Set unstructured.Unstructured `json:"set,omitempty"`
}

// SpecialResourceStatus is the most recently observed status of the SpecialResource.
// It is populated by the system and is read-only.
// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
type SpecialResourceStatus struct {
	// State describes at which step the recipe installation is.
	State string `json:"state"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion

// SpecialResource describes a software stack for hardware accelerators on an existing Kubernetes cluster.
// +kubebuilder:resource:path=specialresources,scope=Cluster
// +kubebuilder:resource:path=specialresources,scope=Cluster,shortName=sr
type SpecialResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SpecialResourceSpec   `json:"spec,omitempty"`
	Status SpecialResourceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SpecialResourceList is a list of SpecialResource objects.
type SpecialResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// List of SpecialResources. More info:
	// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md
	Items []SpecialResource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SpecialResource{}, &SpecialResourceList{})
}
