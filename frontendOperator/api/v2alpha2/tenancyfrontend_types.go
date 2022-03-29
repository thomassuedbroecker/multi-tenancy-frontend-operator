/*
Copyright 2022.

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

package v2alpha2

import (
	"github.com/thomassuedbroecker/multi-tenancy-frontend-operator/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TenancyFrontendSpec defines the desired state of TenancyFrontend
type TenancyFrontendSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Size is an example field of TenancyFrontend. Edit tenancyfrontend_types.go to remove/update
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	Size int32 `json:"size"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=15
	DisplayName string `json:"displayname,omitempty"`
	// +kubebuilder:validation:MaxLength=15
	// +kubebuilder:default:=Movies
	CatalogName string `json:"catalogname,omitempty"`
}

// TenancyFrontendStatus defines the observed state of TenancyFrontend
type TenancyFrontendStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// TenancyFrontend is the Schema for the tenancyfrontends API
type TenancyFrontend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TenancyFrontendSpec   `json:"spec,omitempty"`
	Status TenancyFrontendStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TenancyFrontendList contains a list of TenancyFrontend
type TenancyFrontendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TenancyFrontend `json:"items"`
}

// functions

func init() {
	SchemeBuilder.Register(&TenancyFrontend{}, &TenancyFrontendList{})
}

// ConvertTo converts this v2alpha2 to v1beta1. (upgrade)
func (src *TenancyFrontend) ConvertTo(dstRaw conversion.Hub) error {

	dst := dstRaw.(*v1beta1.TenancyFrontend)
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.DisplayName = src.Spec.DisplayName
	dst.Spec.Size = src.Spec.Size
	dst.Spec.CatalogName = src.Spec.CatalogName

	return nil
}

// ConvertFrom converts from the Hub version (v1beta1) to (v2alpha2). (downgrade)
func (dst *TenancyFrontend) ConvertFrom(srcRaw conversion.Hub) error {

	src := srcRaw.(*v1beta1.TenancyFrontend)

	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.DisplayName = src.Spec.DisplayName
	dst.Spec.Size = src.Spec.Size
	dst.Spec.CatalogName = src.Spec.CatalogName

	return nil
}
