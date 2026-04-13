package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	clusterapis "github.com/enel1221/provider-kion/apis/cluster"
	namespacedapis "github.com/enel1221/provider-kion/apis/namespaced"
)

var (
	testEnv    *envtest.Environment
	testScheme = k8sruntime.NewScheme()
	k8sClient  client.Client
)

func TestMain(m *testing.M) {
	os.Exit(runSuite(m))
}

func runSuite(m *testing.M) int {
	utilruntime.Must(clientgoscheme.AddToScheme(testScheme))
	utilruntime.Must(apiextensionsv1.AddToScheme(testScheme))
	utilruntime.Must(clusterapis.AddToScheme(testScheme))
	utilruntime.Must(namespacedapis.AddToScheme(testScheme))

	repoRoot, err := repositoryRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve repository root: %v\n", err)
		return 1
	}

	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join(repoRoot, "package", "crds")},
		ErrorIfCRDPathMissing: true,
	}

	cfg, err := testEnv.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start envtest: %v\n", err)
		return 1
	}

	k8sClient, err = client.New(cfg, client.Options{Scheme: testScheme})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to construct envtest client: %v\n", err)
		_ = testEnv.Stop()
		return 1
	}

	code := m.Run()
	if err := testEnv.Stop(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to stop envtest: %v\n", err)
		return 1
	}
	return code
}

func repositoryRoot() (string, error) {
	_, currentFile, _, ok := goruntime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to locate current test file")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "..")), nil
}

func testContext(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithTimeout(context.Background(), 30*time.Second)
}

func ensureNamespace(t *testing.T, name string) {
	t.Helper()
	ctx, cancel := testContext(t)
	defer cancel()

	err := k8sClient.Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name}})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		t.Fatalf("failed to ensure namespace %q: %v", name, err)
	}
}

func uniqueName(t *testing.T, prefix string) string {
	t.Helper()
	suffix := strings.ToLower(strings.NewReplacer("/", "-", "_", "-", " ", "-").Replace(t.Name()))
	if len(suffix) > 40 {
		suffix = suffix[len(suffix)-40:]
	}
	name := fmt.Sprintf("%s-%s", prefix, suffix)
	if len(name) > 63 {
		return name[:63]
	}
	return name
}

func boolPtr(value bool) *bool {
	return &value
}

func floatPtr(value float64) *float64 {
	return &value
}

func stringPtr(value string) *string {
	return &value
}

func mustMarshalCredentials(t *testing.T, credentials map[string]any) []byte {
	t.Helper()
	data, err := json.Marshal(credentials)
	if err != nil {
		t.Fatalf("failed to marshal credentials: %v", err)
	}
	return data
}
