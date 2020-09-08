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
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"

	smv1alpha1 "github.com/itscontained/secret-manager/pkg/apis/secretmanager/v1alpha1"

	"github.com/stretchr/testify/assert"
)

func TestSecretsManagerGetSecret(t *testing.T) {
	cases := []struct {
		ref       smv1alpha1.RemoteReference
		expErr    bool
		expSecret []byte
		sm        func(*secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
	}{
		{
			// get secret as string
			ref: smv1alpha1.RemoteReference{
				Path: "/foo/bar/baz",
			},
			expSecret: []byte("HELLO"),
			sm: func(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
				assert.Equal(t, *input.SecretId, "/foo/bar/baz")
				str := "HELLO"
				return &secretsmanager.GetSecretValueOutput{
					SecretString: &str,
				}, nil
			},
		},
		{
			// get property from json
			ref: smv1alpha1.RemoteReference{
				Path:     "/foo/bar/json",
				Property: aws.String("hello"),
			},
			expSecret: []byte("secretsmanager"),
			sm: func(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
				assert.Equal(t, *input.SecretId, "/foo/bar/json")
				str := `{"hello":"secretsmanager"}`
				return &secretsmanager.GetSecretValueOutput{
					SecretString: &str,
				}, nil
			},
		},
		{
			// err: parse json
			ref: smv1alpha1.RemoteReference{
				Path:     "/foo/bar/json",
				Property: aws.String("hello"),
			},
			expErr: true,
			sm: func(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
				assert.Equal(t, *input.SecretId, "/foo/bar/json")
				str := `I AM NO JSON`
				return &secretsmanager.GetSecretValueOutput{
					SecretString: &str,
				}, nil
			},
		},
		{
			// err: json property missing
			ref: smv1alpha1.RemoteReference{
				Path:     "/foo/bar/json",
				Property: aws.String("i dont exist"),
			},
			expErr: true,
			sm: func(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
				assert.Equal(t, *input.SecretId, "/foo/bar/json")
				str := `{"nop":"nothing here"}`
				return &secretsmanager.GetSecretValueOutput{
					SecretString: &str,
				}, nil
			},
		},
		{
			// err at sm
			ref: smv1alpha1.RemoteReference{
				Path: "/foo/bar/baz",
			},
			expErr: true,
			sm: func(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
				return nil, fmt.Errorf("nop")
			},
		},
	}

	for _, c := range cases {
		store := &SecretsManagerStore{
			secretsManager: &mockSecretsManagerClient{
				getSecretValueFunc: c.sm,
			},
		}
		sec, err := store.GetSecret(context.Background(), c.ref)
		if c.expErr {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, sec, c.expSecret)
		}
	}
}

func TestSecretsManagerGetSecretMap(t *testing.T) {
	cases := []struct {
		ref       smv1alpha1.RemoteReference
		expErr    bool
		expSecret map[string][]byte
		sm        func(*secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
	}{
		{
			// get map from json
			ref: smv1alpha1.RemoteReference{
				Path: "/foo/bar/json",
			},
			expSecret: map[string][]byte{
				"a": []byte("A"),
				"b": []byte("B"),
			},
			sm: func(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
				assert.Equal(t, *input.SecretId, "/foo/bar/json")
				str := `{"a":"A","b":"B"}`
				return &secretsmanager.GetSecretValueOutput{
					SecretString: &str,
				}, nil
			},
		},
		{
			// err: sm
			ref: smv1alpha1.RemoteReference{
				Path: "/foo/bar/json",
			},
			expErr: true,
			sm: func(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
				return nil, fmt.Errorf("nop")
			},
		},
		{
			// err: no json
			ref: smv1alpha1.RemoteReference{
				Path: "/foo/bar/json",
			},
			expErr: true,
			sm: func(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
				assert.Equal(t, *input.SecretId, "/foo/bar/json")
				str := `nothing here to be parsed`
				return &secretsmanager.GetSecretValueOutput{
					SecretString: &str,
				}, nil
			},
		},
	}

	for _, c := range cases {
		store := &SecretsManagerStore{
			secretsManager: &mockSecretsManagerClient{
				getSecretValueFunc: c.sm,
			},
		}
		sec, err := store.GetSecretMap(context.Background(), c.ref)
		if c.expErr {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, sec, c.expSecret)
		}
	}
}

type mockSecretsManagerClient struct {
	secretsmanageriface.SecretsManagerAPI
	getSecretValueFunc func(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
}

func (m *mockSecretsManagerClient) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	return m.getSecretValueFunc(input)
}
