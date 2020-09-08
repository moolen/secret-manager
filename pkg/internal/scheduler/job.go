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

package scheduler

import (
	"time"

	smv1alpha1 "github.com/itscontained/secret-manager/pkg/apis/secretmanager/v1alpha1"
)

type Job struct {
	Name           string
	ExternalSecret *smv1alpha1.ExternalSecret
	Func           func()
}

func (j Job) Run() {
	j.Func()
}

type Schedule struct {
	Immediate  bool
	RenewAfter time.Duration
}

func NewSchedule(interval string) (*Schedule, error) {
	dur, err := time.ParseDuration(interval)
	if err != nil {
		return nil, err
	}
	return &Schedule{
		RenewAfter: dur,
	}, nil
}

func (s *Schedule) Next(t time.Time) time.Time {
	if s.Immediate {
		s.Immediate = false
		return t
	}
	return t.Add(s.RenewAfter)
}
