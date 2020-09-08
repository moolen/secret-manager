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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

type sessionProvider struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Role            string
	StsProvider     func(*session.Session) stsiface.STSAPI
}

// TODO: add sessionCache to reuse sessions

var defaultSessionProvider = newSessionProvider

func newSessionProvider(accessKeyID, secretAccessKey, region, role string) *sessionProvider {
	return &sessionProvider{
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		Region:          region,
		Role:            role,
		StsProvider:     defaultSTSProvider,
	}
}

func defaultSTSProvider(sess *session.Session) stsiface.STSAPI {
	return sts.New(sess)
}

func (d *sessionProvider) GetSession() (*session.Session, error) {
	config := aws.NewConfig()
	sessionOpts := session.Options{
		Config: *config,
	}
	if d.AccessKeyID != "" && d.SecretAccessKey != "" {
		sessionOpts.Config.Credentials = credentials.NewStaticCredentials(d.AccessKeyID, d.SecretAccessKey, "")
		sessionOpts.SharedConfigState = session.SharedConfigDisable
	}
	sess, err := session.NewSessionWithOptions(sessionOpts)
	if err != nil {
		return nil, fmt.Errorf("unable to create aws session: %s", err)
	}
	if d.Role != "" {
		stsSvc := d.StsProvider(sess)
		result, err := stsSvc.AssumeRole(&sts.AssumeRoleInput{
			RoleArn:         aws.String(d.Role),
			RoleSessionName: aws.String("secret-manager"),
		})
		if err != nil {
			return nil, fmt.Errorf("unable to assume role: %s", err)
		}
		creds := credentials.Value{
			AccessKeyID:     *result.Credentials.AccessKeyId,
			SecretAccessKey: *result.Credentials.SecretAccessKey,
			SessionToken:    *result.Credentials.SessionToken,
		}
		sessionOpts.Config.Credentials = credentials.NewStaticCredentialsFromCreds(creds)
		sess, err = session.NewSessionWithOptions(sessionOpts)
		if err != nil {
			return nil, fmt.Errorf("unable to create aws session: %s", err)
		}
	}
	// If ambient credentials aren't permitted, always set the region, even if to
	// empty string, to avoid it falling back on the environment.
	// this has to be set after session is constructed
	if d.Region != "" || d.AccessKeyID != "" && d.SecretAccessKey != "" {
		sess.Config.WithRegion(d.Region)
	}
	sess.Handlers.Build.PushBack(request.WithAppendUserAgent("secret-manager"))
	return sess, nil
}
