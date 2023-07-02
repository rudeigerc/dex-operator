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

package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dexv1alpha1 "github.com/rudeigerc/dex-operator/api/v1alpha1"
	"github.com/rudeigerc/dex-operator/internal/reconciler"
)

// DexClusterReconciler reconciles a DexCluster object
type DexClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=dex.rudeigerc.dev,resources=dexclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dex.rudeigerc.dev,resources=dexclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dex.rudeigerc.dev,resources=dexclusters/finalizers,verbs=update

// Reconcile DexCluster.
func (r *DexClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var dexCluster dexv1alpha1.DexCluster
	if err := r.Get(ctx, req.NamespacedName, &dexCluster); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log := ctrl.LoggerFrom(ctx).WithValues("dexCluster", klog.KObj(&dexCluster))
	ctx = ctrl.LoggerInto(ctx, log)

	log.V(2).Info("Reconciling DexCluster")

	var res ctrl.Result
	var err error

	configMapReconciler := reconciler.NewConfigMapReconciler(r.Client)
	res, err = configMapReconciler.Reconcile(ctx, &dexCluster)
	if err != nil {
		return res, err
	}

	deploymentReconciler := reconciler.NewDeploymentReconciler(r.Client)
	res, err = deploymentReconciler.Reconcile(ctx, &dexCluster)
	if err != nil {
		return res, err
	}

	serviceReconciler := reconciler.NewServiceReconciler(r.Client)
	res, err = serviceReconciler.Reconcile(ctx, &dexCluster)
	if err != nil {
		return res, err
	}

	oldStatus := dexCluster.Status

	if !equality.Semantic.DeepEqual(oldStatus, &dexCluster.Status) {
		if err := r.Status().Update(ctx, &dexCluster); err != nil {
			log.Error(err, "unable to update DexCluster status")
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DexClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dexv1alpha1.DexCluster{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
