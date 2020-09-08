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
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"

	smv1alpha1 "github.com/itscontained/secret-manager/pkg/apis/secretmanager/v1alpha1"

	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type SecureSystemsManagerStore struct {
	ssm ssmiface.SSMAPI
}

func NewSecureSystemsManager(ctx context.Context, kubeclient ctrlclient.Client, store smv1alpha1.GenericStore, namespace string) (*SecureSystemsManagerStore, error) {
	awsAccessKeyID, awsSecretAccessKey, err := getCredentialsFromCredentialsRef(ctx, kubeclient, store.GetSpec().AWSParameterStore.Credentials)
	if err != nil {
		return nil, err
	}
	sess, err := defaultSessionProvider(
		awsAccessKeyID,
		awsSecretAccessKey,
		store.GetSpec().AWSParameterStore.Region,
		store.GetSpec().AWSParameterStore.Role).GetSession()
	if err != nil {
		return nil, err
	}
	svc := ssm.New(sess)
	return &SecureSystemsManagerStore{
		ssm: svc,
	}, nil
}

func (s SecureSystemsManagerStore) GetSecret(ctx context.Context, ref smv1alpha1.RemoteReference) ([]byte, error) {
	out, err := s.ssm.GetParameter(&ssm.GetParameterInput{
		Name: &ref.Path,
	})
	if err != nil {
		return nil, fmt.Errorf("could not read parameter %q from AWS Parameter Store", ref.Path)
	}
	if ref.Property != nil {
		m := make(map[string]string)
		err = json.Unmarshal([]byte(*out.Parameter.Value), &m)
		if err != nil {
			return nil, fmt.Errorf("could not read property %s from secret %q from AWS SecretsManager: %s", *ref.Property, ref.Path, err)
		}
		val, ok := m[*ref.Property]
		if !ok {
			return nil, fmt.Errorf("property %s in secret %q from AWS SecretsManager does not exist", *ref.Property, ref.Path)
		}
		return []byte(val), nil
	}
	return []byte(*out.Parameter.Value), nil
}
func (s SecureSystemsManagerStore) GetSecretMap(ctx context.Context, ref smv1alpha1.RemoteReference) (map[string][]byte, error) {
	out, err := s.ssm.GetParameter(&ssm.GetParameterInput{
		Name: &ref.Path,
	})
	if err != nil {
		return nil, fmt.Errorf("could not read secret %q from AWS SecretsManager", ref.Path)
	}
	m := make(map[string]string)
	err = json.Unmarshal([]byte(*out.Parameter.Value), &m)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json from secret %q from AWS SecretsManager: %s", ref.Path, err)
	}
	byteMap := make(map[string][]byte, len(m))
	for k, v := range m {
		byteMap[k] = []byte(v)
	}
	return byteMap, nil
}
