/*
Copyright 2021 Vitaly Elyashev.

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

// CouchbaseIndexSpec defines the desired state of CouchbaseIndices
type CouchbaseIndexSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	//Secret Name, contains details of connection to couchbase cluster
	SecretName string `json:"cbClusterSecret,omitempty"`
	// IndexData is an objects, contains  all the data for  generating couchbase index
	Index IndexData `json:"indexdata,omitempty"`
}

type IndexData struct {
	BucketName string   `json:"bucketname,omitempty"`
	IndexName  string   `json:"indexname,omitempty"`
	IsPrimary  bool     `json:"isprimary,omitempty"`
	Parameters []string `json:"parameters,omitempty"`
}

// CouchbaseIndexStatus defines the observed state of CouchbaseIndex
type CouchbaseIndexStatus struct {
	Type IndexStatusConditionType `json:"conditiontype"`
}

type IndexStatusConditionType string

const (
	ConditionNotExists  IndexStatusConditionType = "NotExists"
	ConditionInProgress IndexStatusConditionType = "InProgress"
	ConditionReady      IndexStatusConditionType = "Ready"
	ConditionFailed     IndexStatusConditionType = "Failed"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status`
// +kubebuilder:printcolumn:name="BucketName",type=string,JSONPath=`.spec.indexdata.bucketname`
// +kubebuilder:printcolumn:name="IndexName",type=string,JSONPath=`.spec.indexdata.indexname`
// +kubebuilder:printcolumn:name="IsPrimary",type=boolean,JSONPath=`.spec.indexdata.isprimary`
// CouchbaseIndex is the Schema for the couchbaseindices API
type CouchbaseIndex struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CouchbaseIndexSpec   `json:"spec,omitempty"`
	Status CouchbaseIndexStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
// CouchbaseIndexList contains a list of CouchbaseIndex
type CouchbaseIndexList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CouchbaseIndex `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CouchbaseIndex{}, &CouchbaseIndexList{})
}
