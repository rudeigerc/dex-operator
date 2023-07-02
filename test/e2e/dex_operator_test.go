package e2e

import (
	"context"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	dexv1alpha1 "github.com/rudeigerc/dex-operator/api/v1alpha1"
)

func TestOperatorSetup(t *testing.T) {
	name := envconf.RandomName("dexcluster", 8)

	feature := features.New("DexCluster Operator").
		Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			dexCluster := newDexCluster(cfg.Namespace(), name)
			if err := cfg.Client().Resources().Create(ctx, dexCluster); err != nil {
				t.Fatal(err)
			}
			time.Sleep(2 * time.Second)
			return ctx
		}).
		Assess("DexCluster creation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			var dexCluster dexv1alpha1.DexCluster
			if err := cfg.Client().Resources().Get(ctx, name, cfg.Namespace(), &dexCluster); err != nil {
				t.Fatal(err)
			}
			if &dexCluster != nil {
				t.Logf("DexCluster found: %s", dexCluster.Name)
			}
			return context.WithValue(ctx, name, &dexCluster)
		}).
		// Assess("Deployment available", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		// 	client, err := cfg.NewClient()
		// 	if err != nil {
		// 		t.Fatal(err)
		// 	}

		// 	var deployment appsv1.Deployment
		// 	if err := client.Resources().Get(ctx, name, cfg.Namespace(), &deployment); err != nil {
		// 		t.Fatal(err)
		// 	}
		// 	if err := wait.For(conditions.New(client.Resources()).DeploymentConditionMatch(&deployment, appsv1.DeploymentAvailable, corev1.ConditionTrue), wait.WithTimeout(time.Minute*2)); err != nil {
		// 		t.Fatal(err)
		// 	}
		// 	return ctx
		// }).
		Teardown(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client, err := cfg.NewClient()
			if err != nil {
				t.Fatal(err)
			}

			dexCluster := ctx.Value(name).(*dexv1alpha1.DexCluster)
			if err := client.Resources().Delete(ctx, dexCluster); err != nil {
				t.Fatal(err)
			}
			err = wait.For(conditions.New(client.Resources()).ResourceDeleted(dexCluster), wait.WithTimeout(time.Minute*1))
			if err != nil {
				t.Fatal(err)
			}
			return ctx
		}).
		Feature()

	testEnv.Test(t, feature)
}

func newDexCluster(namespace string, name string) *dexv1alpha1.DexCluster {
	return &dexv1alpha1.DexCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: dexv1alpha1.DexClusterSpec{
			Image:           "dexidp/dex:v2.37.0",
			ImagePullPolicy: corev1.PullIfNotPresent,
			Config: `
			issuer: http://127.0.0.1:5556/dex
			storage:
				type: sqlite3
				config:
					file: /tmp/dex.db
			web:
				http: 0.0.0.0:5556
			telemetry:
				http: 0.0.0.0:5558
			connectors:
			- type: mockCallback
				id: mock
				name: Example
			`,
		},
	}
}
