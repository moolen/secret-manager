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

package v1alpha1

import (
	smmeta "github.com/itscontained/secret-manager/pkg/apis/meta/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SecretStoreSpec holds the configuration for a single secret store
// +kubebuilder:validation:MinProperties=1
// +kubebuilder:validation:MaxProperties=1
type SecretStoreSpec struct {
	// Vault configures this store to sync secrets using a HashiCorp Vault
	// KV backend.
	// +kubebuilder:validation:Optional
	Vault *VaultStore `json:"vault,omitempty"`

	// +kubebuilder:validation:Optional
	AWSSecretManager *AWSSecretManagerStore `json:"AWSSecretManager,omitempty"`
	// +kubebuilder:validation:Optional
	AWSParameterStore *AWSParameterStoreStore `json:"AWSParameterStore,omitempty"`
}

// AWSSecretManagerStore is a store for AWS Secrets Manager
type AWSSecretManagerStore struct {
	// +kubebuilder:validation:Optional
	Region string `json:"region,omitempty"`

	// +kubebuilder:validation:Optional
	Role string `json:"role,omitempty"`

	// +kubebuilder:validation:Optional
	Secret string `json:"secret,omitempty"`

	// Credentials must have accessKeyID and secretAccessKey
	// to allow secret-manager to access the API
	// +kubebuilder:validation:Optional
	Credentials CredentialsRef `json:"credentials,omitempty"`
}

// AWSParameterStoreStore is a store for AWS Seure Systems Manager Parameter Store
type AWSParameterStoreStore struct {
	// +kubebuilder:validation:Optional
	Region string `json:"region,omitempty"`

	// +kubebuilder:validation:Optional
	Role string `json:"role,omitempty"`

	Parameter string `json:"parameter"`

	// Credentials must have accessKeyID and secretAccessKey
	// to allow secret-manager to access the API
	// +kubebuilder:validation:Optional
	Credentials CredentialsRef `json:"credentials,omitempty"`
}

// CredentialsRef references credentials for a backend
type CredentialsRef struct {
	SecretRef *SecretRef `json:"secretRef"`
	// TODO: consider adding configMapRef
}

// SecretRef is a reference to a secret
// the secret must have specific keys depending on the backend
type SecretRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// Configures an store to sync secrets using a HashiCorp Vault
// KV backend.
type VaultStore struct {
	// Auth configures how secret-manager authenticates with the Vault server.
	Auth VaultAuth `json:"auth"`

	// Server is the connection address for the Vault server, e.g: "https://vault.example.com:8200".
	Server string `json:"server"`

	// Path is the mount path of the Vault KV backend endpoint, e.g:
	// "secret". The "/data" path suffix for fetching secrets from Vault
	// is optional and will be appended if not present in specified path.
	Path string `json:"path"`

	// Version is the Vault KV secret engine version. This can be either "v1" or
	// "v2". Version defaults to "v2".
	// +optional
	Version *VaultKVStoreVersion `json:"version,omitempty"`

	// Name of the vault namespace. Namespaces is a set of features within Vault Enterprise that allows Vault environments to support Secure Multi-tenancy. e.g: "ns1"
	// More about namespaces can be found here https://www.vaultproject.io/docs/enterprise/namespaces
	// +optional
	Namespace *string `json:"namespace,omitempty"`

	// PEM encoded CA bundle used to validate Vault server certificate. Only used
	// if the Server URL is using HTTPS protocol. This parameter is ignored for
	// plain HTTP protocol connection. If not set the system root certificates
	// are used to validate the TLS connection.
	// +optional
	CABundle []byte `json:"caBundle,omitempty"`
}

type VaultKVStoreVersion string

const (
	VaultKVStoreV1 VaultKVStoreVersion = "v1"
	VaultKVStoreV2 VaultKVStoreVersion = "v2"
)

// Configuration used to authenticate with a Vault server.
// Only one of `tokenSecretRef`, `appRole` or `kubernetes` may be specified.
type VaultAuth struct {
	// TokenSecretRef authenticates with Vault by presenting a token.
	// +optional
	TokenSecretRef *smmeta.SecretKeySelector `json:"tokenSecretRef,omitempty"`

	// AppRole authenticates with Vault using the App Role auth mechanism,
	// with the role and secret stored in a Kubernetes Secret resource.
	// +optional
	AppRole *VaultAppRole `json:"appRole,omitempty"`

	// Kubernetes authenticates with Vault by passing the ServiceAccount
	// token stored in the named Secret resource to the Vault server.
	// +optional
	Kubernetes *VaultKubernetesAuth `json:"kubernetes,omitempty"`
}

// VaultAppRole authenticates with Vault using the App Role auth mechanism,
// with the role and secret stored in a Kubernetes Secret resource.
type VaultAppRole struct {
	// Path where the App Role authentication backend is mounted in Vault, e.g:
	// "approle"
	// +kubebuilder:default=approle
	Path string `json:"path"`

	// RoleID configured in the App Role authentication backend when setting
	// up the authentication backend in Vault.
	RoleID string `json:"roleId"`

	// Reference to a key in a Secret that contains the App Role secret used
	// to authenticate with Vault.
	// The `key` field must be specified and denotes which entry within the Secret
	// resource is used as the app role secret.
	SecretRef smmeta.SecretKeySelector `json:"secretRef"`
}

// Authenticate against Vault using a Kubernetes ServiceAccount token stored in
// a Secret.
type VaultKubernetesAuth struct {
	// Path where the Kubernetes authentication backend is mounted in Vault, e.g:
	// "kubernetes"
	// +kubebuilder:default=kubernetes
	Path string `json:"mountPath"`

	// The required Secret field containing a Kubernetes ServiceAccount JWT used
	// for authenticating with Vault. Use of 'ambient credentials' is not
	// supported.
	SecretRef smmeta.SecretKeySelector `json:"secretRef"`

	// A required field containing the Vault Role to assume. A Role binds a
	// Kubernetes ServiceAccount with a set of Vault policies.
	Role string `json:"role"`
}

type SecretStoreStatus struct {
	// List of status conditions to indicate the status of SecretStore.
	// Known condition types are `Ready`.
	// +optional
	Conditions smmeta.ConditionedStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// SecretStore is the Schema for the SecretStore API
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={secretmanager}
type SecretStore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretStoreSpec   `json:"spec,omitempty"`
	Status SecretStoreStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SecretStoreList contains a list of SecretStore
type SecretStoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExternalSecret `json:"items"`
}

// +kubebuilder:object:root=true

// ClusterSecretStore is the Schema for the ClusterSecretStore API
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={secretmanager}
type ClusterSecretStore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretStoreSpec   `json:"spec,omitempty"`
	Status SecretStoreStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterSecretStoreList contains a list of ClusterSecretStore
type ClusterSecretStoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExternalSecret `json:"items"`
}
