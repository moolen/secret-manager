// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	metav1 "github.com/itscontained/secret-manager/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSAuth) DeepCopyInto(out *AWSAuth) {
	*out = *in
	out.Credentials = in.Credentials
	if in.TokenSecretRef != nil {
		in, out := &in.TokenSecretRef, &out.TokenSecretRef
		*out = new(metav1.SecretKeySelector)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSAuth.
func (in *AWSAuth) DeepCopy() *AWSAuth {
	if in == nil {
		return nil
	}
	out := new(AWSAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSCredentials) DeepCopyInto(out *AWSCredentials) {
	*out = *in
	out.SecretRef = in.SecretRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSCredentials.
func (in *AWSCredentials) DeepCopy() *AWSCredentials {
	if in == nil {
		return nil
	}
	out := new(AWSCredentials)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSStore) DeepCopyInto(out *AWSStore) {
	*out = *in
	in.Auth.DeepCopyInto(&out.Auth)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSStore.
func (in *AWSStore) DeepCopy() *AWSStore {
	if in == nil {
		return nil
	}
	out := new(AWSStore)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterSecretStore) DeepCopyInto(out *ClusterSecretStore) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterSecretStore.
func (in *ClusterSecretStore) DeepCopy() *ClusterSecretStore {
	if in == nil {
		return nil
	}
	out := new(ClusterSecretStore)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterSecretStore) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterSecretStoreList) DeepCopyInto(out *ClusterSecretStoreList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ExternalSecret, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterSecretStoreList.
func (in *ClusterSecretStoreList) DeepCopy() *ClusterSecretStoreList {
	if in == nil {
		return nil
	}
	out := new(ClusterSecretStoreList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterSecretStoreList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalSecret) DeepCopyInto(out *ExternalSecret) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalSecret.
func (in *ExternalSecret) DeepCopy() *ExternalSecret {
	if in == nil {
		return nil
	}
	out := new(ExternalSecret)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ExternalSecret) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalSecretList) DeepCopyInto(out *ExternalSecretList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ExternalSecret, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalSecretList.
func (in *ExternalSecretList) DeepCopy() *ExternalSecretList {
	if in == nil {
		return nil
	}
	out := new(ExternalSecretList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ExternalSecretList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalSecretSpec) DeepCopyInto(out *ExternalSecretSpec) {
	*out = *in
	out.StoreRef = in.StoreRef
	if in.RenewAfter != nil {
		in, out := &in.RenewAfter, &out.RenewAfter
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	if in.Data != nil {
		in, out := &in.Data, &out.Data
		*out = make([]KeyReference, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.DataFrom != nil {
		in, out := &in.DataFrom, &out.DataFrom
		*out = make([]RemoteReference, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalSecretSpec.
func (in *ExternalSecretSpec) DeepCopy() *ExternalSecretSpec {
	if in == nil {
		return nil
	}
	out := new(ExternalSecretSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalSecretStatus) DeepCopyInto(out *ExternalSecretStatus) {
	*out = *in
	in.ConditionedStatus.DeepCopyInto(&out.ConditionedStatus)
	if in.RenewalTime != nil {
		in, out := &in.RenewalTime, &out.RenewalTime
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalSecretStatus.
func (in *ExternalSecretStatus) DeepCopy() *ExternalSecretStatus {
	if in == nil {
		return nil
	}
	out := new(ExternalSecretStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeyReference) DeepCopyInto(out *KeyReference) {
	*out = *in
	in.RemoteRef.DeepCopyInto(&out.RemoteRef)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeyReference.
func (in *KeyReference) DeepCopy() *KeyReference {
	if in == nil {
		return nil
	}
	out := new(KeyReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectReference) DeepCopyInto(out *ObjectReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectReference.
func (in *ObjectReference) DeepCopy() *ObjectReference {
	if in == nil {
		return nil
	}
	out := new(ObjectReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteReference) DeepCopyInto(out *RemoteReference) {
	*out = *in
	if in.Path != nil {
		in, out := &in.Path, &out.Path
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.Property != nil {
		in, out := &in.Property, &out.Property
		*out = new(string)
		**out = **in
	}
	if in.Version != nil {
		in, out := &in.Version, &out.Version
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteReference.
func (in *RemoteReference) DeepCopy() *RemoteReference {
	if in == nil {
		return nil
	}
	out := new(RemoteReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretStore) DeepCopyInto(out *SecretStore) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretStore.
func (in *SecretStore) DeepCopy() *SecretStore {
	if in == nil {
		return nil
	}
	out := new(SecretStore)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SecretStore) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretStoreList) DeepCopyInto(out *SecretStoreList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ExternalSecret, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretStoreList.
func (in *SecretStoreList) DeepCopy() *SecretStoreList {
	if in == nil {
		return nil
	}
	out := new(SecretStoreList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SecretStoreList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretStoreSpec) DeepCopyInto(out *SecretStoreSpec) {
	*out = *in
	if in.Vault != nil {
		in, out := &in.Vault, &out.Vault
		*out = new(VaultStore)
		(*in).DeepCopyInto(*out)
	}
	if in.AWS != nil {
		in, out := &in.AWS, &out.AWS
		*out = new(AWSStore)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretStoreSpec.
func (in *SecretStoreSpec) DeepCopy() *SecretStoreSpec {
	if in == nil {
		return nil
	}
	out := new(SecretStoreSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretStoreStatus) DeepCopyInto(out *SecretStoreStatus) {
	*out = *in
	in.Conditions.DeepCopyInto(&out.Conditions)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretStoreStatus.
func (in *SecretStoreStatus) DeepCopy() *SecretStoreStatus {
	if in == nil {
		return nil
	}
	out := new(SecretStoreStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultAppRole) DeepCopyInto(out *VaultAppRole) {
	*out = *in
	out.SecretRef = in.SecretRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultAppRole.
func (in *VaultAppRole) DeepCopy() *VaultAppRole {
	if in == nil {
		return nil
	}
	out := new(VaultAppRole)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultAuth) DeepCopyInto(out *VaultAuth) {
	*out = *in
	if in.TokenSecretRef != nil {
		in, out := &in.TokenSecretRef, &out.TokenSecretRef
		*out = new(metav1.SecretKeySelector)
		**out = **in
	}
	if in.AppRole != nil {
		in, out := &in.AppRole, &out.AppRole
		*out = new(VaultAppRole)
		**out = **in
	}
	if in.Kubernetes != nil {
		in, out := &in.Kubernetes, &out.Kubernetes
		*out = new(VaultKubernetesAuth)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultAuth.
func (in *VaultAuth) DeepCopy() *VaultAuth {
	if in == nil {
		return nil
	}
	out := new(VaultAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultKubernetesAuth) DeepCopyInto(out *VaultKubernetesAuth) {
	*out = *in
	out.SecretRef = in.SecretRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultKubernetesAuth.
func (in *VaultKubernetesAuth) DeepCopy() *VaultKubernetesAuth {
	if in == nil {
		return nil
	}
	out := new(VaultKubernetesAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultStore) DeepCopyInto(out *VaultStore) {
	*out = *in
	in.Auth.DeepCopyInto(&out.Auth)
	if in.Version != nil {
		in, out := &in.Version, &out.Version
		*out = new(VaultKVStoreVersion)
		**out = **in
	}
	if in.Namespace != nil {
		in, out := &in.Namespace, &out.Namespace
		*out = new(string)
		**out = **in
	}
	if in.CABundle != nil {
		in, out := &in.CABundle, &out.CABundle
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultStore.
func (in *VaultStore) DeepCopy() *VaultStore {
	if in == nil {
		return nil
	}
	out := new(VaultStore)
	in.DeepCopyInto(out)
	return out
}
