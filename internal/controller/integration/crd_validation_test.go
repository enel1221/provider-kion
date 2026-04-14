package integration

import (
	"strings"
	"testing"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clusterkionv1alpha1 "github.com/enel1221/provider-kion/apis/cluster/kion/v1alpha1"
	namespacedkionv1alpha1 "github.com/enel1221/provider-kion/apis/namespaced/kion/v1alpha1"
)

func TestClusterNewResourceValidation(t *testing.T) {
	ctx, cancel := testContext(t)
	defer cancel()

	tests := []struct {
		name    string
		object  func(string) client.Object
		wantErr string
	}{
		{
			name: "AppConfigAcceptsGeneratedShape",
			object: func(name string) client.Object {
				return &clusterkionv1alpha1.AppConfig{
					ObjectMeta: metav1.ObjectMeta{Name: name},
					Spec:       clusterkionv1alpha1.AppConfigSpec{ForProvider: clusterkionv1alpha1.AppConfigParameters{}},
				}
			},
		},
		{
			name: "CustomAccountAcceptsReferenceFields",
			object: func(name string) client.Object {
				return &clusterkionv1alpha1.CustomAccount{
					ObjectMeta: metav1.ObjectMeta{Name: name},
					Spec: clusterkionv1alpha1.CustomAccountSpec{ForProvider: clusterkionv1alpha1.CustomAccountParameters{
						AccountNumber:      stringPtr("CUSTOM-100"),
						Name:               stringPtr("Custom Account"),
						PayerID:            floatPtr(1),
						ProjectIDSelector:  &xpv1.Selector{MatchLabels: map[string]string{"testing.upbound.io/example-name": "project"}},
						StartDatecode:      stringPtr("2024-03"),
						SkipAccessChecking: boolPtr(true),
					}},
				}
			},
		},
		{
			name: "CustomAccountRejectsMissingName",
			object: func(name string) client.Object {
				return &clusterkionv1alpha1.CustomAccount{
					ObjectMeta: metav1.ObjectMeta{Name: name},
					Spec: clusterkionv1alpha1.CustomAccountSpec{ForProvider: clusterkionv1alpha1.CustomAccountParameters{
						AccountNumber: stringPtr("CUSTOM-200"),
						PayerID:       floatPtr(1),
					}},
				}
			},
			wantErr: "spec.forProvider.name is a required parameter",
		},
		{
			name: "CustomVariableAcceptsRequiredFields",
			object: func(name string) client.Object {
				return &clusterkionv1alpha1.CustomVariable{
					ObjectMeta: metav1.ObjectMeta{Name: name},
					Spec: clusterkionv1alpha1.CustomVariableSpec{ForProvider: clusterkionv1alpha1.CustomVariableParameters{
						Name:                   stringPtr("environment"),
						Description:            stringPtr("Environment type"),
						KeyValidationMessage:   stringPtr("invalid key"),
						KeyValidationRegex:     stringPtr("^[a-z]+$"),
						Type:                   stringPtr("string"),
						ValueValidationMessage: stringPtr("invalid value"),
						ValueValidationRegex:   stringPtr("^(dev|prod)$"),
					}},
				}
			},
		},
		{
			name: "CustomVariableRejectsMissingDescription",
			object: func(name string) client.Object {
				return &clusterkionv1alpha1.CustomVariable{
					ObjectMeta: metav1.ObjectMeta{Name: name},
					Spec: clusterkionv1alpha1.CustomVariableSpec{ForProvider: clusterkionv1alpha1.CustomVariableParameters{
						Name:                   stringPtr("environment"),
						KeyValidationMessage:   stringPtr("invalid key"),
						KeyValidationRegex:     stringPtr("^[a-z]+$"),
						Type:                   stringPtr("string"),
						ValueValidationMessage: stringPtr("invalid value"),
						ValueValidationRegex:   stringPtr("^(dev|prod)$"),
					}},
				}
			},
			wantErr: "spec.forProvider.description is a required parameter",
		},
		{
			name: "CustomVariableOverrideAcceptsReferenceFields",
			object: func(name string) client.Object {
				return &clusterkionv1alpha1.CustomVariableOverride{
					ObjectMeta: metav1.ObjectMeta{Name: name},
					Spec: clusterkionv1alpha1.CustomVariableOverrideSpec{ForProvider: clusterkionv1alpha1.CustomVariableOverrideParameters{
						CustomVariableIDSelector: &xpv1.Selector{MatchLabels: map[string]string{"testing.upbound.io/example-name": "custom-variable"}},
						EntityID:                 stringPtr("123"),
						EntityType:               stringPtr("project"),
						ValueString:              stringPtr("prod"),
					}},
				}
			},
		},
		{
			name: "CustomVariableOverrideRejectsMissingEntityID",
			object: func(name string) client.Object {
				return &clusterkionv1alpha1.CustomVariableOverride{
					ObjectMeta: metav1.ObjectMeta{Name: name},
					Spec: clusterkionv1alpha1.CustomVariableOverrideSpec{ForProvider: clusterkionv1alpha1.CustomVariableOverrideParameters{
						EntityType:  stringPtr("project"),
						ValueString: stringPtr("prod"),
					}},
				}
			},
			wantErr: "spec.forProvider.entityId is a required parameter",
		},
		{
			name: "ProjectNoteAcceptsReferenceFields",
			object: func(name string) client.Object {
				return &clusterkionv1alpha1.ProjectNote{
					ObjectMeta: metav1.ObjectMeta{Name: name},
					Spec: clusterkionv1alpha1.ProjectNoteSpec{ForProvider: clusterkionv1alpha1.ProjectNoteParameters{
						CreateUserID:      floatPtr(15),
						Name:              stringPtr("Project Overview"),
						ProjectIDSelector: &xpv1.Selector{MatchLabels: map[string]string{"testing.upbound.io/example-name": "project"}},
						Text:              stringPtr("Project documentation"),
					}},
				}
			},
		},
		{
			name: "ProjectNoteRejectsMissingText",
			object: func(name string) client.Object {
				return &clusterkionv1alpha1.ProjectNote{
					ObjectMeta: metav1.ObjectMeta{Name: name},
					Spec: clusterkionv1alpha1.ProjectNoteSpec{ForProvider: clusterkionv1alpha1.ProjectNoteParameters{
						CreateUserID: floatPtr(15),
						Name:         stringPtr("Project Overview"),
					}},
				}
			},
			wantErr: "spec.forProvider.text is a required parameter",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			resourceName := uniqueName(t, "cluster")
			err := k8sClient.Create(ctx, testCase.object(resourceName))
			if testCase.wantErr == "" {
				if err != nil {
					t.Fatalf("expected create to succeed, got error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected create to fail with %q", testCase.wantErr)
			}
			if !strings.Contains(err.Error(), testCase.wantErr) {
				t.Fatalf("expected error to contain %q, got %v", testCase.wantErr, err)
			}
		})
	}
}

func TestNamespacedNewResourceValidation(t *testing.T) {
	namespace := uniqueName(t, "validation")
	ensureNamespace(t, namespace)

	ctx, cancel := testContext(t)
	defer cancel()

	tests := []struct {
		name    string
		object  func(string, string) client.Object
		wantErr string
	}{
		{
			name: "AppConfigAcceptsGeneratedShape",
			object: func(namespace, name string) client.Object {
				return &namespacedkionv1alpha1.AppConfig{
					ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
					Spec:       namespacedkionv1alpha1.AppConfigSpec{ForProvider: namespacedkionv1alpha1.AppConfigParameters{}},
				}
			},
		},
		{
			name: "CustomAccountAcceptsReferenceFields",
			object: func(namespace, name string) client.Object {
				return &namespacedkionv1alpha1.CustomAccount{
					ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
					Spec: namespacedkionv1alpha1.CustomAccountSpec{ForProvider: namespacedkionv1alpha1.CustomAccountParameters{
						AccountNumber:      stringPtr("CUSTOM-100"),
						Name:               stringPtr("Custom Account"),
						PayerID:            floatPtr(1),
						ProjectIDSelector:  &xpv1.NamespacedSelector{MatchLabels: map[string]string{"testing.upbound.io/example-name": "project"}},
						StartDatecode:      stringPtr("2024-03"),
						SkipAccessChecking: boolPtr(true),
					}},
				}
			},
		},
		{
			name: "CustomAccountRejectsMissingName",
			object: func(namespace, name string) client.Object {
				return &namespacedkionv1alpha1.CustomAccount{
					ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
					Spec: namespacedkionv1alpha1.CustomAccountSpec{ForProvider: namespacedkionv1alpha1.CustomAccountParameters{
						AccountNumber: stringPtr("CUSTOM-200"),
						PayerID:       floatPtr(1),
					}},
				}
			},
			wantErr: "spec.forProvider.name is a required parameter",
		},
		{
			name: "CustomVariableAcceptsRequiredFields",
			object: func(namespace, name string) client.Object {
				return &namespacedkionv1alpha1.CustomVariable{
					ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
					Spec: namespacedkionv1alpha1.CustomVariableSpec{ForProvider: namespacedkionv1alpha1.CustomVariableParameters{
						Name:                   stringPtr("environment"),
						Description:            stringPtr("Environment type"),
						KeyValidationMessage:   stringPtr("invalid key"),
						KeyValidationRegex:     stringPtr("^[a-z]+$"),
						Type:                   stringPtr("string"),
						ValueValidationMessage: stringPtr("invalid value"),
						ValueValidationRegex:   stringPtr("^(dev|prod)$"),
					}},
				}
			},
		},
		{
			name: "CustomVariableRejectsMissingDescription",
			object: func(namespace, name string) client.Object {
				return &namespacedkionv1alpha1.CustomVariable{
					ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
					Spec: namespacedkionv1alpha1.CustomVariableSpec{ForProvider: namespacedkionv1alpha1.CustomVariableParameters{
						Name:                   stringPtr("environment"),
						KeyValidationMessage:   stringPtr("invalid key"),
						KeyValidationRegex:     stringPtr("^[a-z]+$"),
						Type:                   stringPtr("string"),
						ValueValidationMessage: stringPtr("invalid value"),
						ValueValidationRegex:   stringPtr("^(dev|prod)$"),
					}},
				}
			},
			wantErr: "spec.forProvider.description is a required parameter",
		},
		{
			name: "CustomVariableOverrideAcceptsReferenceFields",
			object: func(namespace, name string) client.Object {
				return &namespacedkionv1alpha1.CustomVariableOverride{
					ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
					Spec: namespacedkionv1alpha1.CustomVariableOverrideSpec{ForProvider: namespacedkionv1alpha1.CustomVariableOverrideParameters{
						CustomVariableIDSelector: &xpv1.NamespacedSelector{MatchLabels: map[string]string{"testing.upbound.io/example-name": "custom-variable"}},
						EntityID:                 stringPtr("123"),
						EntityType:               stringPtr("project"),
						ValueString:              stringPtr("prod"),
					}},
				}
			},
		},
		{
			name: "CustomVariableOverrideRejectsMissingEntityID",
			object: func(namespace, name string) client.Object {
				return &namespacedkionv1alpha1.CustomVariableOverride{
					ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
					Spec: namespacedkionv1alpha1.CustomVariableOverrideSpec{ForProvider: namespacedkionv1alpha1.CustomVariableOverrideParameters{
						EntityType:  stringPtr("project"),
						ValueString: stringPtr("prod"),
					}},
				}
			},
			wantErr: "spec.forProvider.entityId is a required parameter",
		},
		{
			name: "ProjectNoteAcceptsReferenceFields",
			object: func(namespace, name string) client.Object {
				return &namespacedkionv1alpha1.ProjectNote{
					ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
					Spec: namespacedkionv1alpha1.ProjectNoteSpec{ForProvider: namespacedkionv1alpha1.ProjectNoteParameters{
						CreateUserID:      floatPtr(15),
						Name:              stringPtr("Project Overview"),
						ProjectIDSelector: &xpv1.NamespacedSelector{MatchLabels: map[string]string{"testing.upbound.io/example-name": "project"}},
						Text:              stringPtr("Project documentation"),
					}},
				}
			},
		},
		{
			name: "ProjectNoteRejectsMissingText",
			object: func(namespace, name string) client.Object {
				return &namespacedkionv1alpha1.ProjectNote{
					ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
					Spec: namespacedkionv1alpha1.ProjectNoteSpec{ForProvider: namespacedkionv1alpha1.ProjectNoteParameters{
						CreateUserID: floatPtr(15),
						Name:         stringPtr("Project Overview"),
					}},
				}
			},
			wantErr: "spec.forProvider.text is a required parameter",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			resourceName := uniqueName(t, "namespaced")
			err := k8sClient.Create(ctx, testCase.object(namespace, resourceName))
			if testCase.wantErr == "" {
				if err != nil {
					t.Fatalf("expected create to succeed, got error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected create to fail with %q", testCase.wantErr)
			}
			if !strings.Contains(err.Error(), testCase.wantErr) {
				t.Fatalf("expected error to contain %q, got %v", testCase.wantErr, err)
			}
		})
	}
}
