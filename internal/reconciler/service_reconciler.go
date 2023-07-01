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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	dexv1alpha1 "github.com/rudeigerc/dex-operator/api/v1alpha1"
)

type ServiceReconciler struct {
	client client.Client
}

func NewServiceReconciler(client client.Client) DexClusterReconciler {
	return &ServiceReconciler{
		client: client,
	}
}

// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
func (r *ServiceReconciler) Reconcile(ctx context.Context, dexCluster *dexv1alpha1.DexCluster) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dexCluster.Name,
			Namespace: dexCluster.Namespace,
		},
	}
	_, err := controllerutil.CreateOrUpdate(ctx, r.client, service, func() error {
		service.Spec = serviceManifest(dexCluster).Spec
		err := controllerutil.SetControllerReference(dexCluster, service, r.client.Scheme())
		if err != nil {
			log.Error(err, "unable to set controller reference for Service")
			return err
		}
		return nil
	})

	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func serviceManifest(dexCluster *dexv1alpha1.DexCluster) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dexCluster.Name,
			Namespace: dexCluster.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app.kubernetes.io/managed-by": "dex-operator",
				"app.kubernetes.io/instance":   dexCluster.Name,
				"app.kubernetes.io/component":  "dex-cluster",
			},
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     5556,
					Protocol: corev1.ProtocolTCP,
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
	return service
}
