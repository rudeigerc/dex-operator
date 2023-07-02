/*
Copyright 2023 Yuchen Cheng.

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

package reconciler

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	dexv1alpha1 "github.com/rudeigerc/dex-operator/api/v1alpha1"
)

type ConfigMapReconciler struct {
	client client.Client
}

func NewConfigMapReconciler(client client.Client) DexClusterReconciler {
	return &ConfigMapReconciler{
		client: client,
	}
}

//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile ConfigMap.
func (r *ConfigMapReconciler) Reconcile(ctx context.Context, dexCluster *dexv1alpha1.DexCluster) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-config", dexCluster.Name),
			Namespace: dexCluster.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.client, configMap, func() error {
		if configMap.Data == nil {
			configMap.Data = make(map[string]string)
		}
		configMap.Data["config.yaml"] = dexCluster.Spec.Config
		err := controllerutil.SetControllerReference(dexCluster, configMap, r.client.Scheme())
		if err != nil {
			log.Error(err, "unable to set controller reference for ConfigMap")
			return err
		}
		return nil
	})

	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
