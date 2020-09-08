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
package aws

import (
	"context"
	"fmt"

	smv1alpha1 "github.com/itscontained/secret-manager/pkg/apis/secretmanager/v1alpha1"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/types"

	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	secretKeyAccessKeyID     = "accessKeyID"
	secretKeySecretAccessKey = "secretAccessKey"
)

func getCredentialsFromSecret(ctx context.Context, kubeclient ctrlclient.Client, name, namespace string) (string, string, error) {
	secret := &corev1.Secret{}
	ref := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
	err := kubeclient.Get(ctx, ref, secret)
	if err != nil {
		return "", "", fmt.Errorf("unable to fetch secret: %s", err)
	}

	idBytes, ok := secret.Data[secretKeyAccessKeyID]
	if !ok {
		return "", "", fmt.Errorf("no data for %q in secret '%s/%s'", secretKeyAccessKeyID, name, namespace)
	}
	id := string(idBytes)
	secBytes, ok := secret.Data[secretKeySecretAccessKey]
	if !ok {
		return "", "", fmt.Errorf("no data for %q in secret '%s/%s'", secretKeySecretAccessKey, name, namespace)
	}
	key := string(secBytes)
	return id, key, nil
}

func getCredentialsFromCredentialsRef(ctx context.Context, kubeclient ctrlclient.Client, credRef smv1alpha1.CredentialsRef) (string, string, error) {
	var awsAccessKeyID string
	var awsSecretAccessKey string
	if credRef.SecretRef != nil {
		name := credRef.SecretRef.Name
		namespace := credRef.SecretRef.Namespace
		var err error
		awsAccessKeyID, awsSecretAccessKey, err = getCredentialsFromSecret(ctx, kubeclient, name, namespace)
		if err != nil {
			return "", "", err
		}
	}
	return awsAccessKeyID, awsSecretAccessKey, nil
}
