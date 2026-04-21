package integration

import (
	"strings"
	"testing"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	xpv2 "github.com/crossplane/crossplane-runtime/v2/apis/common/v2"
	"github.com/crossplane/upjet/v2/pkg/terraform"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clusterkionv1alpha1 "github.com/enel1221/provider-kion/apis/cluster/kion/v1alpha1"
	clusterv1beta1 "github.com/enel1221/provider-kion/apis/cluster/v1beta1"
	namespacedkionv1alpha1 "github.com/enel1221/provider-kion/apis/namespaced/kion/v1alpha1"
	namespacedv1beta1 "github.com/enel1221/provider-kion/apis/namespaced/v1beta1"
	"github.com/enel1221/provider-kion/internal/clients"
	"github.com/enel1221/provider-kion/internal/testconfig"
)

func newTerraformSetupFn(t *testing.T) terraform.SetupFn {
	t.Helper()

	terraformCfg, err := testconfig.LoadTerraformConfig()
	if err != nil {
		t.Fatalf("failed to load terraform test config: %v", err)
	}

	return clients.TerraformSetupBuilder(terraformCfg.Version, terraformCfg.ProviderSource, terraformCfg.ProviderVersion)
}

func TestTerraformSetupBuilderLegacyProviderConfig(t *testing.T) {
	ctx, cancel := testContext(t)
	defer cancel()

	const credentialsNamespace = "crossplane-system"
	ensureNamespace(t, credentialsNamespace)

	secretName := uniqueName(t, "legacy-creds")
	providerConfigName := uniqueName(t, "legacy-pc")
	managedResourceName := uniqueName(t, "legacy-appconfig")
	expected := terraform.ProviderConfiguration{
		"apikey":            "legacy-api-key",
		"url":               "https://legacy.example.com",
		"skipsslvalidation": true,
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: credentialsNamespace},
		Data: map[string][]byte{
			"credentials": mustMarshalCredentials(t, map[string]any(expected)),
		},
	}
	if err := k8sClient.Create(ctx, secret); err != nil {
		t.Fatalf("failed to create secret: %v", err)
	}

	providerConfig := &clusterv1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{Name: providerConfigName},
		Spec: clusterv1beta1.ProviderConfigSpec{
			Credentials: clusterv1beta1.ProviderCredentials{
				Source: xpv1.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
					SecretRef: &xpv1.SecretKeySelector{
						SecretReference: xpv1.SecretReference{Name: secretName, Namespace: credentialsNamespace},
						Key:             "credentials",
					},
				},
			},
		},
	}
	if err := k8sClient.Create(ctx, providerConfig); err != nil {
		t.Fatalf("failed to create legacy provider config: %v", err)
	}

	managed := &clusterkionv1alpha1.AppConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: clusterkionv1alpha1.AppConfig_GroupVersionKind.GroupVersion().String(),
			Kind:       clusterkionv1alpha1.AppConfig_GroupVersionKind.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{Name: managedResourceName},
		Spec: clusterkionv1alpha1.AppConfigSpec{
			ForProvider: clusterkionv1alpha1.AppConfigParameters{},
		},
	}
	managed.Spec.ProviderConfigReference = &xpv1.Reference{Name: providerConfigName}
	if err := k8sClient.Create(ctx, managed); err != nil {
		t.Fatalf("failed to create legacy managed resource: %v", err)
	}
	managed.GetObjectKind().SetGroupVersionKind(clusterkionv1alpha1.AppConfig_GroupVersionKind)

	setupFn := newTerraformSetupFn(t)
	setup, err := setupFn(ctx, k8sClient, managed)
	if err != nil {
		t.Fatalf("TerraformSetupBuilder returned unexpected error: %v", err)
	}

	if diff := cmp.Diff(expected, setup.Configuration); diff != "" {
		t.Fatalf("Terraform configuration mismatch (-want +got):\n%s", diff)
	}

	usageList := &clusterv1beta1.ProviderConfigUsageList{}
	if err := k8sClient.List(ctx, usageList); err != nil {
		t.Fatalf("failed to list legacy provider config usages: %v", err)
	}

	found := false
	for _, usage := range usageList.Items {
		if usage.ProviderConfigReference.Name == providerConfigName && usage.ResourceReference.Name == managedResourceName {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected a ProviderConfigUsage for provider %q and resource %q", providerConfigName, managedResourceName)
	}
}

func TestTerraformSetupBuilderNamespacedProviderConfig(t *testing.T) {
	ctx, cancel := testContext(t)
	defer cancel()

	namespace := uniqueName(t, "team")
	ensureNamespace(t, namespace)

	secretName := uniqueName(t, "ns-creds")
	providerConfigName := uniqueName(t, "ns-pc")
	managedResourceName := uniqueName(t, "ns-appconfig")
	expected := terraform.ProviderConfiguration{
		"apikey":            "modern-api-key",
		"url":               "https://namespaced.example.com",
		"skipsslvalidation": false,
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: namespace},
		Data: map[string][]byte{
			"credentials": mustMarshalCredentials(t, map[string]any(expected)),
		},
	}
	if err := k8sClient.Create(ctx, secret); err != nil {
		t.Fatalf("failed to create namespaced secret: %v", err)
	}

	providerConfig := &namespacedv1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{Name: providerConfigName, Namespace: namespace},
		Spec: namespacedv1beta1.ProviderConfigSpec{
			Credentials: namespacedv1beta1.ProviderCredentials{
				Source: xpv1.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
					SecretRef: &xpv1.SecretKeySelector{
						SecretReference: xpv1.SecretReference{Name: secretName},
						Key:             "credentials",
					},
				},
			},
		},
	}
	if err := k8sClient.Create(ctx, providerConfig); err != nil {
		t.Fatalf("failed to create namespaced provider config: %v", err)
	}

	managed := &namespacedkionv1alpha1.AppConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: namespacedkionv1alpha1.AppConfig_GroupVersionKind.GroupVersion().String(),
			Kind:       namespacedkionv1alpha1.AppConfig_GroupVersionKind.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{Name: managedResourceName, Namespace: namespace},
		Spec: namespacedkionv1alpha1.AppConfigSpec{
			ManagedResourceSpec: xpv2.ManagedResourceSpec{
				ProviderConfigReference: &xpv1.ProviderConfigReference{Name: providerConfigName, Kind: "ProviderConfig"},
			},
			ForProvider: namespacedkionv1alpha1.AppConfigParameters{},
		},
	}
	if err := k8sClient.Create(ctx, managed); err != nil {
		t.Fatalf("failed to create namespaced managed resource: %v", err)
	}
	managed.GetObjectKind().SetGroupVersionKind(namespacedkionv1alpha1.AppConfig_GroupVersionKind)

	setupFn := newTerraformSetupFn(t)
	setup, err := setupFn(ctx, k8sClient, managed)
	if err != nil {
		t.Fatalf("TerraformSetupBuilder returned unexpected error: %v", err)
	}

	if diff := cmp.Diff(expected, setup.Configuration); diff != "" {
		t.Fatalf("Terraform configuration mismatch (-want +got):\n%s", diff)
	}

	usageList := &namespacedv1beta1.ProviderConfigUsageList{}
	if err := k8sClient.List(ctx, usageList, client.InNamespace(namespace)); err != nil {
		t.Fatalf("failed to list namespaced provider config usages: %v", err)
	}

	found := false
	for _, usage := range usageList.Items {
		if usage.ProviderConfigReference.Name == providerConfigName && usage.ResourceReference.Name == managedResourceName {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected a namespaced ProviderConfigUsage for provider %q and resource %q", providerConfigName, managedResourceName)
	}
}

func TestTerraformSetupBuilderClusterProviderConfig(t *testing.T) {
	ctx, cancel := testContext(t)
	defer cancel()

	const credentialsNamespace = "crossplane-system"
	namespace := uniqueName(t, "team")
	ensureNamespace(t, credentialsNamespace)
	ensureNamespace(t, namespace)

	secretName := uniqueName(t, "cluster-creds")
	providerConfigName := uniqueName(t, "cluster-pc")
	managedResourceName := uniqueName(t, "cluster-appconfig")
	expected := terraform.ProviderConfiguration{
		"apikey":            "cluster-api-key",
		"url":               "https://cluster.example.com",
		"skipsslvalidation": true,
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: credentialsNamespace},
		Data: map[string][]byte{
			"credentials": mustMarshalCredentials(t, map[string]any(expected)),
		},
	}
	if err := k8sClient.Create(ctx, secret); err != nil {
		t.Fatalf("failed to create cluster-scoped secret: %v", err)
	}

	providerConfig := &namespacedv1beta1.ClusterProviderConfig{
		ObjectMeta: metav1.ObjectMeta{Name: providerConfigName},
		Spec: namespacedv1beta1.ProviderConfigSpec{
			Credentials: namespacedv1beta1.ProviderCredentials{
				Source: xpv1.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
					SecretRef: &xpv1.SecretKeySelector{
						SecretReference: xpv1.SecretReference{Name: secretName, Namespace: credentialsNamespace},
						Key:             "credentials",
					},
				},
			},
		},
	}
	if err := k8sClient.Create(ctx, providerConfig); err != nil {
		t.Fatalf("failed to create cluster provider config: %v", err)
	}

	managed := &namespacedkionv1alpha1.AppConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: namespacedkionv1alpha1.AppConfig_GroupVersionKind.GroupVersion().String(),
			Kind:       namespacedkionv1alpha1.AppConfig_GroupVersionKind.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{Name: managedResourceName, Namespace: namespace},
		Spec: namespacedkionv1alpha1.AppConfigSpec{
			ManagedResourceSpec: xpv2.ManagedResourceSpec{
				ProviderConfigReference: &xpv1.ProviderConfigReference{Name: providerConfigName, Kind: "ClusterProviderConfig"},
			},
			ForProvider: namespacedkionv1alpha1.AppConfigParameters{},
		},
	}
	if err := k8sClient.Create(ctx, managed); err != nil {
		t.Fatalf("failed to create managed resource: %v", err)
	}
	managed.GetObjectKind().SetGroupVersionKind(namespacedkionv1alpha1.AppConfig_GroupVersionKind)

	setupFn := newTerraformSetupFn(t)
	setup, err := setupFn(ctx, k8sClient, managed)
	if err != nil {
		t.Fatalf("TerraformSetupBuilder returned unexpected error: %v", err)
	}

	if diff := cmp.Diff(expected, setup.Configuration); diff != "" {
		t.Fatalf("Terraform configuration mismatch (-want +got):\n%s", diff)
	}
}

func TestTerraformSetupBuilderMissingSecret(t *testing.T) {
	ctx, cancel := testContext(t)
	defer cancel()

	namespace := uniqueName(t, "team")
	ensureNamespace(t, namespace)

	providerConfigName := uniqueName(t, "missing-secret-pc")
	managedResourceName := uniqueName(t, "missing-secret-appconfig")

	providerConfig := &namespacedv1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{Name: providerConfigName, Namespace: namespace},
		Spec: namespacedv1beta1.ProviderConfigSpec{
			Credentials: namespacedv1beta1.ProviderCredentials{
				Source: xpv1.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
					SecretRef: &xpv1.SecretKeySelector{
						SecretReference: xpv1.SecretReference{Name: "missing-secret"},
						Key:             "credentials",
					},
				},
			},
		},
	}
	if err := k8sClient.Create(ctx, providerConfig); err != nil {
		t.Fatalf("failed to create provider config: %v", err)
	}

	managed := &namespacedkionv1alpha1.AppConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: namespacedkionv1alpha1.AppConfig_GroupVersionKind.GroupVersion().String(),
			Kind:       namespacedkionv1alpha1.AppConfig_GroupVersionKind.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{Name: managedResourceName, Namespace: namespace},
		Spec: namespacedkionv1alpha1.AppConfigSpec{
			ManagedResourceSpec: xpv2.ManagedResourceSpec{
				ProviderConfigReference: &xpv1.ProviderConfigReference{Name: providerConfigName, Kind: "ProviderConfig"},
			},
			ForProvider: namespacedkionv1alpha1.AppConfigParameters{},
		},
	}
	if err := k8sClient.Create(ctx, managed); err != nil {
		t.Fatalf("failed to create managed resource: %v", err)
	}
	managed.GetObjectKind().SetGroupVersionKind(namespacedkionv1alpha1.AppConfig_GroupVersionKind)

	setupFn := newTerraformSetupFn(t)
	_, err := setupFn(ctx, k8sClient, managed)
	if err == nil {
		t.Fatal("expected TerraformSetupBuilder to fail when the referenced secret is missing")
	}
	if !strings.Contains(err.Error(), "cannot extract credentials") {
		t.Fatalf("expected missing secret error to contain %q, got %v", "cannot extract credentials", err)
	}
}
