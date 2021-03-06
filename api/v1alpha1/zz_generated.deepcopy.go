//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CouchbaseIndex) DeepCopyInto(out *CouchbaseIndex) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CouchbaseIndex.
func (in *CouchbaseIndex) DeepCopy() *CouchbaseIndex {
	if in == nil {
		return nil
	}
	out := new(CouchbaseIndex)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CouchbaseIndex) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CouchbaseIndexList) DeepCopyInto(out *CouchbaseIndexList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CouchbaseIndex, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CouchbaseIndexList.
func (in *CouchbaseIndexList) DeepCopy() *CouchbaseIndexList {
	if in == nil {
		return nil
	}
	out := new(CouchbaseIndexList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CouchbaseIndexList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CouchbaseIndexSpec) DeepCopyInto(out *CouchbaseIndexSpec) {
	*out = *in
	in.Index.DeepCopyInto(&out.Index)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CouchbaseIndexSpec.
func (in *CouchbaseIndexSpec) DeepCopy() *CouchbaseIndexSpec {
	if in == nil {
		return nil
	}
	out := new(CouchbaseIndexSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CouchbaseIndexStatus) DeepCopyInto(out *CouchbaseIndexStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CouchbaseIndexStatus.
func (in *CouchbaseIndexStatus) DeepCopy() *CouchbaseIndexStatus {
	if in == nil {
		return nil
	}
	out := new(CouchbaseIndexStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IndexData) DeepCopyInto(out *IndexData) {
	*out = *in
	if in.Parameters != nil {
		in, out := &in.Parameters, &out.Parameters
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IndexData.
func (in *IndexData) DeepCopy() *IndexData {
	if in == nil {
		return nil
	}
	out := new(IndexData)
	in.DeepCopyInto(out)
	return out
}
