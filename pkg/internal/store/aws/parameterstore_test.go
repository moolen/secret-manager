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
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"

	smv1alpha1 "github.com/itscontained/secret-manager/pkg/apis/secretmanager/v1alpha1"

	"github.com/stretchr/testify/assert"
)

func TestParameterStoreGetSecret(t *testing.T) {
	cases := []struct {
		ref       smv1alpha1.RemoteReference
		expErr    bool
		expSecret []byte
		ssm       func(*ssm.GetParameterInput) (*ssm.GetParameterOutput, error)
	}{
		{
			// get secret as string
			ref: smv1alpha1.RemoteReference{
				Path: "/foo/bar/baz",
			},
			expSecret: []byte("HELLO"),
			ssm: func(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
				assert.Equal(t, *input.Name, "/foo/bar/baz")
				str := "HELLO"
				return &ssm.GetParameterOutput{
					Parameter: &ssm.Parameter{
						Value: &str,
					},
				}, nil
			},
		},
		{
			// get property from json
			ref: smv1alpha1.RemoteReference{
				Path:     "/foo/bar/json",
				Property: aws.String("hello"),
			},
			expSecret: []byte("param"),
			ssm: func(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
				assert.Equal(t, *input.Name, "/foo/bar/json")
				str := `{"hello":"param"}`
				return &ssm.GetParameterOutput{
					Parameter: &ssm.Parameter{
						Value: &str,
					},
				}, nil
			},
		},
		{
			// err: get parameter
			ref: smv1alpha1.RemoteReference{
				Path: "/foo/bar/baz",
			},
			expErr: true,
			ssm: func(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
				return nil, fmt.Errorf("nop")
			},
		},
		{
			// err: parse json
			ref: smv1alpha1.RemoteReference{
				Path:     "/foo/bar/json",
				Property: aws.String("hello"),
			},
			expErr: true,
			ssm: func(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
				assert.Equal(t, *input.Name, "/foo/bar/json")
				str := `i am no json`
				return &ssm.GetParameterOutput{
					Parameter: &ssm.Parameter{
						Value: &str,
					},
				}, nil
			},
		},
		{
			// err: missing property at json
			ref: smv1alpha1.RemoteReference{
				Path:     "/foo/bar/json",
				Property: aws.String("missing prop"),
			},
			expErr: true,
			ssm: func(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
				assert.Equal(t, *input.Name, "/foo/bar/json")
				str := `{"nop":"here is nothing"}`
				return &ssm.GetParameterOutput{
					Parameter: &ssm.Parameter{
						Value: &str,
					},
				}, nil
			},
		},
	}

	for _, c := range cases {
		store := &SecureSystemsManagerStore{
			ssm: &mockSystemsManagerClient{
				getParameterFunc: c.ssm,
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

func TestParameterStoreGetSecretMap(t *testing.T) {
	cases := []struct {
		ref       smv1alpha1.RemoteReference
		expErr    bool
		expSecret map[string][]byte
		ssm       func(*ssm.GetParameterInput) (*ssm.GetParameterOutput, error)
	}{
		{
			ref: smv1alpha1.RemoteReference{
				Path: "/foo/bar/baz",
			},
			expSecret: map[string][]byte{
				"a": []byte("A"),
				"b": []byte("B"),
			},
			ssm: func(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
				assert.Equal(t, *input.Name, "/foo/bar/baz")
				str := `{"a":"A","b":"B"}`
				return &ssm.GetParameterOutput{
					Parameter: &ssm.Parameter{
						Value: &str,
					},
				}, nil
			},
		},
		{
			// err: get parameter
			ref: smv1alpha1.RemoteReference{
				Path: "/foo/bar/baz",
			},
			expErr: true,
			ssm: func(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
				return nil, fmt.Errorf("nop")
			},
		},
		{
			// err: parse json
			ref: smv1alpha1.RemoteReference{
				Path:     "/foo/bar/json",
				Property: aws.String("hello"),
			},
			expErr: true,
			ssm: func(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
				assert.Equal(t, *input.Name, "/foo/bar/json")
				str := `no JSON here`
				return &ssm.GetParameterOutput{
					Parameter: &ssm.Parameter{
						Value: &str,
					},
				}, nil
			},
		},
	}

	for _, c := range cases {
		store := &SecureSystemsManagerStore{
			ssm: &mockSystemsManagerClient{
				getParameterFunc: c.ssm,
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

type mockSystemsManagerClient struct {
	ssmiface.SSMAPI
	getParameterFunc func(*ssm.GetParameterInput) (*ssm.GetParameterOutput, error)
}

func (m *mockSystemsManagerClient) GetParameter(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	return m.getParameterFunc(input)
}
