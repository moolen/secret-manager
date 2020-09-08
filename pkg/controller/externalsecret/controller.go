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
	"fmt"
	"time"

	"github.com/go-logr/logr"

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

	// requeueAfter = time.Second * 30

	errStoreNotFound          = "cannot get store reference"
	errStoreSetupFailed       = "cannot setup store client"
	errGetSecretDataFailed    = "cannot get ExternalSecret data from store"
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

	// recon func used by scheduler
	fn := func() error {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      extSecret.Name,
				Namespace: extSecret.Namespace,
			},
		}

		// TODO: frontend specific strategy
		result, err := ctrl.CreateOrUpdate(ctx, r.Client, secret, func() error {
			store, err := r.getStore(ctx, extSecret)
			if err != nil {
				return fmt.Errorf("%s: %w", errStoreNotFound, err)
			}

			storeClient, err := r.storeFactory.New(ctx, store, r.Client, req.Namespace)
			if err != nil {
				return fmt.Errorf("%s: %w", errStoreSetupFailed, err)
			}

			secret.ObjectMeta.Labels = extSecret.Labels
			secret.ObjectMeta.Annotations = extSecret.Annotations
			err = controllerutil.SetControllerReference(extSecret, &secret.ObjectMeta, r.Scheme)
			if err != nil {
				return fmt.Errorf("failed to set ExternalSecret controller reference: %w", err)
			}
			secret.Data, err = r.getSecret(ctx, storeClient, extSecret)
			if err != nil {
				return fmt.Errorf("%s: %w", errGetSecretDataFailed, err)
			}
			return nil
		})
		if err != nil {
			log.Error(err, "error while reconciling ExternalSecret", "namespace", extSecret.Namespace, "name", extSecret.Name)
			extSecret.Status.SetConditions(smmeta.Unavailable().WithMessage(err.Error()))
			if extSecret.Spec.RenewAfter != nil {
				extSecret.Status.RenewalTime = &metav1.Time{Time: time.Now().Add(extSecret.Spec.RenewAfter.Duration)}
			}
			_ = r.Status().Update(ctx, extSecret)
			return fmt.Errorf("%s: %w", errUpdateSecretDataFailed, err)
		}
		log.Info("successfully reconcile ExternalSecret", "operation", result)
		extSecret.Status.SetConditions(smmeta.Available())
		if extSecret.Spec.RenewAfter != nil {
			extSecret.Status.RenewalTime = &metav1.Time{Time: time.Now().Add(extSecret.Spec.RenewAfter.Duration)}
		}
		_ = r.Status().Update(ctx, extSecret)
		return nil
	}

	if extSecret.Spec.RenewAfter != nil {
		log.Info("adding to schedule", "namespace", extSecret.Namespace, "name", extSecret.Name)
		r.Scheduler.Add(extSecret, fn)
	}

	err := fn()
	if err != nil {
		log.Error(err, "error while reconciling ExternalSecret")
		extSecret.Status.SetConditions(smmeta.Unavailable().WithMessage(err.Error()))
		if extSecret.Spec.RenewAfter != nil {
			extSecret.Status.RenewalTime = &metav1.Time{Time: time.Now().Add(extSecret.Spec.RenewAfter.Duration)}
		}
		_ = r.Status().Update(ctx, extSecret)
		return ctrl.Result{}, nil
	}

	extSecret.Status.SetConditions(smmeta.Available())
	if extSecret.Spec.RenewAfter != nil {
		extSecret.Status.RenewalTime = &metav1.Time{Time: time.Now().Add(extSecret.Spec.RenewAfter.Duration)}
	}
	_ = r.Status().Update(ctx, extSecret)
	return ctrl.Result{}, nil
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

func (r *ExternalSecretReconciler) getSecret(ctx context.Context, storeClient store.Client, extSecret *smv1alpha1.ExternalSecret) (map[string][]byte, error) {
	secretDataMap := make(map[string][]byte)
	for _, remoteRef := range extSecret.Spec.DataFrom {
		secretMap, err := storeClient.GetSecretMap(ctx, remoteRef)
		if err != nil {
			return nil, fmt.Errorf("path %q: %w", remoteRef.Path, err)
		}
		secretDataMap = merge.Merge(secretDataMap, secretMap)
	}

	for _, secretRef := range extSecret.Spec.Data {
		secretData, err := storeClient.GetSecret(ctx, secretRef.RemoteRef)
		if err != nil {
			return nil, fmt.Errorf("path %q: %w", secretRef.RemoteRef.Path, err)
		}
		secretDataMap[secretRef.SecretKey] = secretData
	}

	for secretKey, secretData := range secretDataMap {
		dstBytes := make([]byte, base64.RawStdEncoding.EncodedLen(len(secretData)))
		base64.RawStdEncoding.Encode(dstBytes, secretData)
		secretDataMap[secretKey] = dstBytes
	}

	return secretDataMap, nil
}

func (r *ExternalSecretReconciler) getStore(ctx context.Context, extSecret *smv1alpha1.ExternalSecret) (smv1alpha1.GenericStore, error) {
	if extSecret.Kind == smv1alpha1.ClusterSecretStoreKind {
		clusterStore := &smv1alpha1.ClusterSecretStore{}
		ref := types.NamespacedName{
			Name: extSecret.Spec.StoreRef.Name,
		}
		if err := r.Get(ctx, ref, clusterStore); err != nil {
			return nil, fmt.Errorf("ClusterSecretStore %q: %w", ref.Name, err)
		}
		return clusterStore, nil
	}

	namespacedStore := &smv1alpha1.SecretStore{}
	ref := types.NamespacedName{
		Namespace: extSecret.Namespace,
		Name:      extSecret.Spec.StoreRef.Name,
	}
	if err := r.Get(ctx, ref, namespacedStore); err != nil {
		return nil, fmt.Errorf("SecretStore %q: %w", ref.Name, err)
	}
	return namespacedStore, nil
}
