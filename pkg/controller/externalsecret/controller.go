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

package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-logr/logr"

	"github.com/imdario/mergo"

	smmeta "github.com/itscontained/secret-manager/pkg/apis/meta/v1"
	smv1alpha1 "github.com/itscontained/secret-manager/pkg/apis/secretmanager/v1alpha1"
	"github.com/itscontained/secret-manager/pkg/internal/scheduler"
	"github.com/itscontained/secret-manager/pkg/internal/store"
	storebase "github.com/itscontained/secret-manager/pkg/internal/store/base"
	"github.com/itscontained/secret-manager/pkg/util/merge"

	corev1 "k8s.io/api/core/v1"

	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/utils/clock"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	ownerKey = ".metadata.controller"

	requeueAfter = time.Second * 30

	errStoreNotFound          = "cannot get store reference"
	errStoreSetupFailed       = "cannot setup store client"
	errGetSecretDataFailed    = "cannot get ExternalSecret data from store"
	errTemplateFailed         = "failed to merge secret with template field"
	errUpdateSecretDataFailed = "cannot create/update ExternalSecret data from store"
)

// ExternalSecretReconciler reconciles a ExternalSecret object
type ExternalSecretReconciler struct {
	client.Client
	Log       logr.Logger
	Scheme    *runtime.Scheme
	Clock     clock.Clock
	Scheduler *scheduler.Scheduler

	storeFactory store.Factory
	Reader       client.Reader
}

type externalSecretSyncer struct {
	client       client.Client
	extSecret    *smv1alpha1.ExternalSecret
	log          logr.Logger
	scheme       *runtime.Scheme
	storeFactory store.Factory
}

func (r *ExternalSecretReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("externalsecret", req.NamespacedName)
	log.V(2).Info("reconciling ExternalSecret")
	extSecret := &smv1alpha1.ExternalSecret{}
	if err := r.Get(ctx, req.NamespacedName, extSecret); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("deleting object")
			r.Scheduler.Remove(req.NamespacedName)
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to get ExternalSecret")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	es := externalSecretSyncer{
		client:       r,
		extSecret:    extSecret,
		log:          r.Log,
		scheme:       r.Scheme,
		storeFactory: r.storeFactory,
	}
	if shouldSchedule(extSecret) {
		log.V(2).Info("adding to schedule", "namespace", extSecret.Namespace, "name", extSecret.Name)
		r.Scheduler.Add(extSecret, es.sync)
	}
	// skip sync depending on refreshInterval
	if skipSync(extSecret) {
		return ctrl.Result{}, nil
	}
	err := es.sync()
	if err != nil {
		log.Error(err, "error while reconciling ExternalSecret")
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}
	return ctrl.Result{}, nil
}

func (ess *externalSecretSyncer) sync() error {
	es := ess.extSecret
	ctx := context.Background()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      es.Name,
			Namespace: es.Namespace,
		},
	}

	if es.Spec.RefreshInterval != nil {
		if es.Spec.RefreshInterval.Duration.Seconds() == 0 {
			es.Status.NextSync = &metav1.Time{
				Time: time.Unix(0, 0),
			}
		} else {
			es.Status.NextSync = &metav1.Time{
				Time: time.Now().UTC().Add(es.Spec.RefreshInterval.Duration),
			}
		}
	}

	defer func() {
		err := ess.client.Status().Update(ctx, es)
		if err != nil {
			ess.log.Error(err, "error while updating ExternalSecret Status field", "namespace", es.Namespace, "name", es.Name)
		}
	}()

	_, err := ctrl.CreateOrUpdate(ctx, ess.client, secret, func() error {
		store, err := getStore(ctx, ess.client, es)
		if err != nil {
			return fmt.Errorf("%s: %w", errStoreNotFound, err)
		}

		storeClient, err := ess.storeFactory.New(ctx, store, ess.client, ess.client, es.ObjectMeta.Namespace)
		if err != nil {
			return fmt.Errorf("%s: %w", errStoreSetupFailed, err)
		}

		secret.ObjectMeta.Labels = es.Labels
		secret.ObjectMeta.Annotations = es.Annotations
		err = controllerutil.SetControllerReference(es, &secret.ObjectMeta, ess.scheme)
		if err != nil {
			return fmt.Errorf("failed to set ExternalSecret controller reference: %w", err)
		}
		secret.Data, err = getSecret(ctx, storeClient, es)
		if err != nil {
			return fmt.Errorf("%s: %w", errGetSecretDataFailed, err)
		}
		if es.Spec.Template != nil {
			err = templateSecret(secret, es.Spec.Template)
			if err != nil {
				return fmt.Errorf("%s: %w", errTemplateFailed, err)
			}
		}
		return nil
	})
	if err != nil {
		err = fmt.Errorf("%s: %w", errUpdateSecretDataFailed, err)
		es.Status.SetConditions(smmeta.Unavailable().WithMessage(err.Error()))
		return err
	}
	es.Status.SetConditions(smmeta.Available())
	return nil
}

func (r *ExternalSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.Clock == nil {
		r.Clock = clock.RealClock{}
	}

	if r.storeFactory == nil {
		r.storeFactory = &storebase.Default{}
	}

	if r.Scheduler == nil {
		r.Scheduler = scheduler.New(r.storeFactory, r, r.Log)
	}

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.Secret{}, ownerKey, func(rawObj runtime.Object) []string {
		secret := rawObj.(*corev1.Secret)
		owner := metav1.GetControllerOf(secret)
		if owner == nil {
			return nil
		}

		if owner.APIVersion != smv1alpha1.ExtSecretGroupVersionKind.GroupVersion().String() || owner.Kind != smv1alpha1.ExtSecretKind {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&smv1alpha1.ExternalSecret{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}

func getSecret(ctx context.Context, storeClient store.Client, extSecret *smv1alpha1.ExternalSecret) (map[string][]byte, error) {
	secretDataMap := make(map[string][]byte)
	for _, remoteRef := range extSecret.Spec.DataFrom {
		secretMap, err := storeClient.GetSecretMap(ctx, remoteRef)
		if err != nil {
			if remoteRef.Name != nil {
				return nil, fmt.Errorf("path %q: %w", *remoteRef.Name, err)
			}
			return nil, fmt.Errorf("name %q: %w", *remoteRef.Name, err)
		}
		secretDataMap = merge.Merge(secretDataMap, secretMap)
	}

	for _, secretRef := range extSecret.Spec.Data {
		secretData, err := storeClient.GetSecret(ctx, secretRef.RemoteRef)
		if err != nil {
			if secretRef.RemoteRef.Name != nil {
				return nil, fmt.Errorf("path %q: %w", *secretRef.RemoteRef.Name, err)
			}
			return nil, fmt.Errorf("name %q: %w", *secretRef.RemoteRef.Name, err)
		}
		secretDataMap[secretRef.SecretKey] = secretData
	}

	for secretKey, secretData := range secretDataMap {
		dstBytes := make([]byte, base64.StdEncoding.EncodedLen(len(secretData)))
		base64.StdEncoding.Encode(dstBytes, secretData)
		secretDataMap[secretKey] = dstBytes
	}

	return secretDataMap, nil
}

func getStore(ctx context.Context, cl client.Client, extSecret *smv1alpha1.ExternalSecret) (smv1alpha1.GenericStore, error) {
	if extSecret.Kind == smv1alpha1.ClusterSecretStoreKind {
		clusterStore := &smv1alpha1.ClusterSecretStore{}
		ref := types.NamespacedName{
			Name: extSecret.Spec.StoreRef.Name,
		}
		if err := cl.Get(ctx, ref, clusterStore); err != nil {
			return nil, fmt.Errorf("ClusterSecretStore %q: %w", ref.Name, err)
		}
		return clusterStore, nil
	}
	var namespacedStore smv1alpha1.SecretStore
	ref := types.NamespacedName{
		Namespace: extSecret.Namespace,
		Name:      extSecret.Spec.StoreRef.Name,
	}
	if err := cl.Get(ctx, ref, &namespacedStore); err != nil {
		return nil, fmt.Errorf("SecretStore %q: %w", ref.Name, err)
	}
	return &namespacedStore, nil
}

func templateSecret(secret *corev1.Secret, template []byte) error {
	templatedSecret := &corev1.Secret{}
	if err := json.Unmarshal(template, templatedSecret); err != nil {
		return fmt.Errorf("error unmarshalling json: %w", err)
	}

	return mergo.Merge(secret, templatedSecret, mergo.WithOverride)
}

func shouldSchedule(extSecret *smv1alpha1.ExternalSecret) bool {
	return extSecret.Spec.RefreshInterval != nil &&
		extSecret.Spec.RefreshInterval.Seconds() >= 60
}

func skipSync(extSecret *smv1alpha1.ExternalSecret) bool {
	return extSecret.Spec.RefreshInterval != nil &&
		extSecret.Spec.RefreshInterval.Seconds() == 0 &&
		extSecret.Status.NextSync != nil
}
