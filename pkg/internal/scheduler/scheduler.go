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
	"fmt"
	"sync"

	"github.com/go-logr/logr"

	smv1alpha1 "github.com/itscontained/secret-manager/pkg/apis/secretmanager/v1alpha1"
	"github.com/itscontained/secret-manager/pkg/internal/store"

	"github.com/robfig/cron/v3"

	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Scheduler struct {
	cron         *cron.Cron
	storeFactory store.Factory
	client       client.Client
	log          logr.Logger
	mu           sync.RWMutex
	scheduleMap  map[string]cron.EntryID
}

func New(storeFactory store.Factory, client client.Client, logger logr.Logger) *Scheduler {
	return &Scheduler{
		cron:         cron.New(),
		log:          logger,
		storeFactory: storeFactory,
		client:       client,
		mu:           sync.RWMutex{},
		scheduleMap:  make(map[string]cron.EntryID),
	}
}

func (s *Scheduler) Add(extSecret *smv1alpha1.ExternalSecret, fn func() error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.log.Info("enqueing schedule for", "namespace", extSecret.Namespace, "name", extSecret.Name)
	id := identity(extSecret.Namespace, extSecret.Name)
	entryID := s.cron.Schedule(&Schedule{
		Immediate:       true,
		RefreshInterval: extSecret.Spec.RefreshInterval.Duration,
	}, Job{
		Name:           id,
		ExternalSecret: extSecret,
		Func: func() {
			s.log.Info("adding schedule for", "namespace", extSecret.Namespace, "name", extSecret.Name)
			err := fn()
			if err != nil {
				s.log.Error(err, "error running scheduled job", "namespace", extSecret.Namespace, "name", extSecret.Name)
			}
		},
	})
	s.scheduleMap[id] = entryID
}

func (s *Scheduler) Remove(nsn types.NamespacedName) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := identity(nsn.Namespace, nsn.Name)
	entry, ok := s.scheduleMap[id]
	if ok {
		s.cron.Remove(entry)
	}
}

func identity(namespace, name string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}
