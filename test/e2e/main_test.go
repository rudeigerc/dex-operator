package e2e

import (
	"os"
	"testing"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	log "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"

	dexv1alpha1 "github.com/rudeigerc/dex-operator/api/v1alpha1"
)

var (
	testEnv          env.Environment
	image            = os.Getenv("IMG")
	e2eTestAssetsDir = os.Getenv("E2ETEST_ASSETS_DIR")
)

func TestMain(m *testing.M) {
	utilruntime.Must(dexv1alpha1.AddToScheme(scheme.Scheme))

	log.SetLogger(klog.NewKlogr())

	testEnv = env.New()
	kindClusterName := envconf.RandomName("dex-operator", 16)
	namespace := envconf.RandomName("dex-operator-ns", 16)

	testEnv.Setup(
		envfuncs.CreateKindClusterWithConfig(kindClusterName, "kindest/node:v1.27.3", "kind-config.yaml"),
		envfuncs.CreateNamespace(namespace),
		envfuncs.LoadDockerImageToCluster(kindClusterName, image),
		envfuncs.SetupCRDs(e2eTestAssetsDir, "dex-operator.yaml"),
	).Finish(
		envfuncs.DeleteNamespace(namespace),
		envfuncs.TeardownCRDs(e2eTestAssetsDir, "dex-operator.yaml"),
		envfuncs.DestroyKindCluster(kindClusterName),
	)
	os.Exit(testEnv.Run(m))
}
