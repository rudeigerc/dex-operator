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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	dexv1alpha1 "github.com/rudeigerc/dex-operator/api/v1alpha1"
)

type DeploymentReconciler struct {
	client client.Client
}

func NewDeploymentReconciler(client client.Client) DexClusterReconciler {
	return &DeploymentReconciler{
		client: client,
	}
}

//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile Deployment.
func (r *DeploymentReconciler) Reconcile(ctx context.Context, dexCluster *dexv1alpha1.DexCluster) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dexCluster.Name,
			Namespace: dexCluster.Namespace,
		},
	}

	desired := deploymentManifest(dexCluster)
	_, err := controllerutil.CreateOrUpdate(ctx, r.client, deployment, func() error {
		// immutable selector
		if deployment.ObjectMeta.CreationTimestamp.IsZero() {
			deployment.Spec.Selector = desired.Spec.Selector
		}
		deployment.Labels = desired.Labels
		deployment.Spec.Replicas = desired.Spec.Replicas
		deployment.Spec.Template = desired.Spec.Template

		err := controllerutil.SetControllerReference(dexCluster, deployment, r.client.Scheme())
		if err != nil {
			log.Error(err, "unable to set controller reference for Deployment")
			return err
		}
		return nil
	})

	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func deploymentManifest(dexCluster *dexv1alpha1.DexCluster) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: dexCluster.Name + "-",
			Namespace:    dexCluster.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "dex-operator",
				"app.kubernetes.io/instance":   dexCluster.Name,
				"app.kubernetes.io/component":  "dex-cluster",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: dexCluster.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/managed-by": "dex-operator",
					"app.kubernetes.io/instance":   dexCluster.Name,
					"app.kubernetes.io/component":  "dex-cluster",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/managed-by": "dex-operator",
						"app.kubernetes.io/instance":   dexCluster.Name,
						"app.kubernetes.io/component":  "dex-cluster",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{containerManifest(dexCluster)},
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{Name: fmt.Sprintf("%s-config", dexCluster.Name)},
									Items: []corev1.KeyToPath{{
										Key:  "config.yaml",
										Path: "config.yaml",
									}},
								},
							},
						},
					},
				},
			},
		},
	}
	return deployment
}

func containerManifest(dexCluster *dexv1alpha1.DexCluster) corev1.Container {
	ports := []corev1.ContainerPort{
		{
			Name:          "http",
			ContainerPort: 5556,
		},
		{
			Name:          "grpc",
			ContainerPort: 5557,
		},
		{
			Name:          "telemetry",
			ContainerPort: 5558,
		},
	}

	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "config",
			MountPath: "/etc/dex/cfg",
		},
	}

	return corev1.Container{
		Name:            "dex",
		Image:           dexCluster.Spec.Image,
		Command:         []string{"/usr/local/bin/dex", "serve", "/etc/dex/cfg/config.yaml"},
		ImagePullPolicy: dexCluster.Spec.ImagePullPolicy,
		Ports:           ports,
		VolumeMounts:    volumeMounts,
		ReadinessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   "/healthz/ready",
					Port:   intstr.FromString("telemetry"),
					Scheme: corev1.URISchemeHTTP,
				},
			},
		},
		LivenessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   "/healthz/live",
					Port:   intstr.FromString("telemetry"),
					Scheme: corev1.URISchemeHTTP,
				},
			},
		},
	}
}

func GetDeployment(dexCluster *dexv1alpha1.DexCluster, scheme *runtime.Scheme) *appsv1.Deployment {
	config := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-deployment", dexCluster.Name),
			Namespace: dexCluster.Namespace,
		},
	}
	return config
}
