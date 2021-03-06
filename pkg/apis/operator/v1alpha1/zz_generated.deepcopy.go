// +build !ignore_autogenerated

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

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GroupResourceQuotaEnforcer) DeepCopyInto(out *GroupResourceQuotaEnforcer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GroupResourceQuotaEnforcer.
func (in *GroupResourceQuotaEnforcer) DeepCopy() *GroupResourceQuotaEnforcer {
	if in == nil {
		return nil
	}
	out := new(GroupResourceQuotaEnforcer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GroupResourceQuotaEnforcer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GroupResourceQuotaEnforcerList) DeepCopyInto(out *GroupResourceQuotaEnforcerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]GroupResourceQuotaEnforcer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GroupResourceQuotaEnforcerList.
func (in *GroupResourceQuotaEnforcerList) DeepCopy() *GroupResourceQuotaEnforcerList {
	if in == nil {
		return nil
	}
	out := new(GroupResourceQuotaEnforcerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GroupResourceQuotaEnforcerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GroupResourceQuotaEnforcerSpec) DeepCopyInto(out *GroupResourceQuotaEnforcerSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GroupResourceQuotaEnforcerSpec.
func (in *GroupResourceQuotaEnforcerSpec) DeepCopy() *GroupResourceQuotaEnforcerSpec {
	if in == nil {
		return nil
	}
	out := new(GroupResourceQuotaEnforcerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GroupResourceQuotaEnforcerStatus) DeepCopyInto(out *GroupResourceQuotaEnforcerStatus) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GroupResourceQuotaEnforcerStatus.
func (in *GroupResourceQuotaEnforcerStatus) DeepCopy() *GroupResourceQuotaEnforcerStatus {
	if in == nil {
		return nil
	}
	out := new(GroupResourceQuotaEnforcerStatus)
	in.DeepCopyInto(out)
	return out
}
