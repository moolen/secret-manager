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
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	"github.com/stretchr/testify/assert"
)

func restoreRoute53Env() {
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_REGION")
}

func TestSessionProvider(t *testing.T) {
	os.Setenv("AWS_ACCESS_KEY_ID", "123")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "123")
	os.Setenv("AWS_REGION", "us-east-1")
	defer restoreRoute53Env()

	roleCreds := &sts.Credentials{
		AccessKeyId:     aws.String("666"),
		SecretAccessKey: aws.String("777"),
		SessionToken:    aws.String("secret-manager"),
	}
	envCreds := &sts.Credentials{
		AccessKeyId:     aws.String("123"),
		SecretAccessKey: aws.String("123"),
		SessionToken:    aws.String("secret-manager"),
	}
	cases := []struct {
		name            string
		accessKeyID     string
		secretAccessKey string
		role            string
		expErr          bool
		expCreds        *sts.Credentials
		region          string
		mockSTS         *mockSTS
	}{
		{
			name:     "should assume role",
			role:     "my-role",
			region:   "us-east-1",
			expErr:   false,
			expCreds: roleCreds,
			mockSTS: &mockSTS{
				AssumeRoleFn: func(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
					assert.Equal(t, "my-role", *input.RoleArn)
					return &sts.AssumeRoleOutput{
						Credentials: roleCreds,
					}, nil
				},
			},
		},
		{
			name:     "fails to assume role",
			role:     "my-role",
			region:   "us-east-1",
			expErr:   true,
			expCreds: nil,
			mockSTS: &mockSTS{
				AssumeRoleFn: func(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
					assert.Equal(t, "my-role", *input.RoleArn)
					return nil, fmt.Errorf("nop. wrong credentials")
				},
			},
		},
		{
			name:     "should use credentials from env",
			role:     "",
			region:   "us-east-1",
			expErr:   false,
			expCreds: envCreds,
			mockSTS:  nil,
		},
		{
			name:            "should use explicit credentials, not env",
			accessKeyID:     "000000",
			secretAccessKey: "888888",
			region:          "eu-central-1",
			expErr:          false,
			expCreds: &sts.Credentials{
				AccessKeyId:     aws.String("000000"),
				SecretAccessKey: aws.String("888888"),
				SessionToken:    aws.String("secret-manager"),
			},
			mockSTS: nil,
		},
	}

	for _, c := range cases {
		provider := newSessionProvider(c.accessKeyID, c.secretAccessKey, c.region, c.role)
		if c.mockSTS != nil {
			prov := c.mockSTS
			provider.StsProvider = func(sess *session.Session) stsiface.STSAPI {
				return prov
			}
		}
		sess, err := provider.GetSession()
		if c.expErr {
			assert.NotNil(t, err)
		} else {
			sessCreds, _ := sess.Config.Credentials.Get()
			if c.mockSTS != nil {
				assert.Equal(t, c.mockSTS.assumedRole, c.role)
			}
			assert.Equal(t, *c.expCreds.SecretAccessKey, sessCreds.SecretAccessKey)
			assert.Equal(t, *c.expCreds.AccessKeyId, sessCreds.AccessKeyID)
			assert.Equal(t, c.region, *sess.Config.Region)
		}
	}
}

type mockSTS struct {
	*sts.STS
	AssumeRoleFn func(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error)
	assumedRole  string
}

func (m *mockSTS) AssumeRole(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	if m.AssumeRoleFn != nil {
		m.assumedRole = *input.RoleArn
		return m.AssumeRoleFn(input)
	}

	return nil, nil
}
