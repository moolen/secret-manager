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

package framework

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// CreateAWSSecretsManagerSecret creates a sm secret with the given value
func CreateAWSSecretsManagerSecret(namespace, name, secret string) error {
	sess, err := newSession(namespace)
	if err != nil {
		return err
	}
	sm := secretsmanager.New(sess)
	// we have to use aws sdk v1! v2 panics, see here: https://gist.github.com/moolen/d453f843bfbe67ab5ff747f06678a33f
	_, err = sm.CreateSecret(&secretsmanager.CreateSecretInput{
		Name:         aws.String(name),
		SecretString: aws.String(secret),
	})
	if err != nil {
		return err
	}
	return nil
}

func newSession(namespace string) (*session.Session, error) {
	defaultResolver := endpoints.DefaultResolver()
	resolverFunc := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if service == "secretsmanager" {
			return endpoints.ResolvedEndpoint{
				URL: fmt.Sprintf("http://localstack.%s.svc.cluster.local", namespace),
			}, nil
		}
		if service == "sts" {
			return endpoints.ResolvedEndpoint{
				URL: fmt.Sprintf("http://localstack.%s.svc.cluster.local", namespace),
			}, nil
		}
		return defaultResolver.EndpointFor(service, region, optFns...)
	}
	cfg := &aws.Config{
		Region:           aws.String("us-east-1"),
		EndpointResolver: endpoints.ResolverFunc(resolverFunc),
		Credentials:      credentials.NewStaticCredentials("foobar", "foobar", "secret-manager-e2e"),
	}
	return session.NewSession(cfg)
}
