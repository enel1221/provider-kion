/*
Copyright 2021 Upbound Inc.
*/

package clients

import (
	"context"
	"encoding/json"
	"testing"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	xpresource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource/fake"
	"github.com/crossplane/crossplane-runtime/v2/pkg/test"
	"github.com/crossplane/upjet/v2/pkg/terraform"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clusterv1beta1 "github.com/enel1221/provider-kion/apis/cluster/v1beta1"
	namespacedv1beta1 "github.com/enel1221/provider-kion/apis/namespaced/v1beta1"
)

// validCredentials returns a JSON-encoded credential blob for testing.
func validCredentials() []byte {
	creds := map[string]any{
		keyAPIKey:            "test-api-key",
		keyURL:               "https://kion.example.com",
		keySkipSSLValidation: true,
	}
	b, _ := json.Marshal(creds)
	return b
}

func TestTerraformSetupBuilder(t *testing.T) {
	type args struct {
		version         string
		providerSource  string
		providerVersion string
	}
	type want struct {
		configuration terraform.ProviderConfiguration
		err           error
	}
	cases := map[string]struct {
		reason string
		args   args
		mg     xpresource.Managed
		kube   client.Client
		want   want
	}{
		"LegacyManagedSecretCredentials": {
			reason: "Should extract credentials from a Secret for a legacy (cluster-scoped) managed resource.",
			args:   args{version: "0.0.1", providerSource: "registry.terraform.io/kionsoftware/kion", providerVersion: "0.3.31"},
			mg: &fake.LegacyManaged{
				ObjectMeta: metav1.ObjectMeta{Name: "test-mr"},
				LegacyProviderConfigReferencer: fake.LegacyProviderConfigReferencer{
					Ref: &xpv1.Reference{Name: "test-pc"},
				},
			},
			kube: &test.MockClient{
				MockGet: func(_ context.Context, key client.ObjectKey, obj client.Object) error {
					switch o := obj.(type) {
					case *clusterv1beta1.ProviderConfig:
						o.Spec = clusterv1beta1.ProviderConfigSpec{
							Credentials: clusterv1beta1.ProviderCredentials{
								Source: xpv1.CredentialsSourceSecret,
								CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
									SecretRef: &xpv1.SecretKeySelector{
										SecretReference: xpv1.SecretReference{
											Name:      "kion-creds",
											Namespace: "crossplane-system",
										},
										Key: "credentials",
									},
								},
							},
						}
					case *corev1.Secret:
						o.Data = map[string][]byte{
							"credentials": validCredentials(),
						}
					}
					return nil
				},
				MockStatusUpdate: test.NewMockSubResourceUpdateFn(nil),
				MockCreate:       test.NewMockCreateFn(nil),
			},
			want: want{
				configuration: terraform.ProviderConfiguration{
					keyAPIKey:            "test-api-key",
					keyURL:               "https://kion.example.com",
					keySkipSSLValidation: true,
				},
			},
		},
		"LegacyManagedNilProviderConfigRef": {
			reason: "Should return an error when providerConfigRef is nil.",
			args:   args{version: "0.0.1", providerSource: "test", providerVersion: "0.0.1"},
			mg: &fake.LegacyManaged{
				ObjectMeta: metav1.ObjectMeta{Name: "test-mr"},
				LegacyProviderConfigReferencer: fake.LegacyProviderConfigReferencer{
					Ref: nil,
				},
			},
			kube: &test.MockClient{},
			want: want{
				err: errAsNonNil,
			},
		},
		"LegacyManagedInvalidCredentialJSON": {
			reason: "Should return an error when credential data is not valid JSON.",
			args:   args{version: "0.0.1", providerSource: "test", providerVersion: "0.0.1"},
			mg: &fake.LegacyManaged{
				ObjectMeta: metav1.ObjectMeta{Name: "test-mr"},
				LegacyProviderConfigReferencer: fake.LegacyProviderConfigReferencer{
					Ref: &xpv1.Reference{Name: "test-pc"},
				},
			},
			kube: &test.MockClient{
				MockGet: func(_ context.Context, key client.ObjectKey, obj client.Object) error {
					switch o := obj.(type) {
					case *clusterv1beta1.ProviderConfig:
						o.Spec = clusterv1beta1.ProviderConfigSpec{
							Credentials: clusterv1beta1.ProviderCredentials{
								Source: xpv1.CredentialsSourceSecret,
								CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
									SecretRef: &xpv1.SecretKeySelector{
										SecretReference: xpv1.SecretReference{
											Name:      "kion-creds",
											Namespace: "crossplane-system",
										},
										Key: "credentials",
									},
								},
							},
						}
					case *corev1.Secret:
						o.Data = map[string][]byte{
							"credentials": []byte("not-json"),
						}
					}
					return nil
				},
				MockStatusUpdate: test.NewMockSubResourceUpdateFn(nil),
				MockCreate:       test.NewMockCreateFn(nil),
			},
			want: want{
				err: errAsNonNil,
			},
		},
		"PartialCredentials": {
			reason: "Should populate only the keys present in the credential JSON.",
			args:   args{version: "0.0.1", providerSource: "test", providerVersion: "0.0.1"},
			mg: &fake.LegacyManaged{
				ObjectMeta: metav1.ObjectMeta{Name: "test-mr"},
				LegacyProviderConfigReferencer: fake.LegacyProviderConfigReferencer{
					Ref: &xpv1.Reference{Name: "test-pc"},
				},
			},
			kube: &test.MockClient{
				MockGet: func(_ context.Context, key client.ObjectKey, obj client.Object) error {
					switch o := obj.(type) {
					case *clusterv1beta1.ProviderConfig:
						o.Spec = clusterv1beta1.ProviderConfigSpec{
							Credentials: clusterv1beta1.ProviderCredentials{
								Source: xpv1.CredentialsSourceSecret,
								CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
									SecretRef: &xpv1.SecretKeySelector{
										SecretReference: xpv1.SecretReference{
											Name:      "kion-creds",
											Namespace: "crossplane-system",
										},
										Key: "credentials",
									},
								},
							},
						}
					case *corev1.Secret:
						// Only apikey — no url or skipsslvalidation
						creds, _ := json.Marshal(map[string]any{keyAPIKey: "only-key"})
						o.Data = map[string][]byte{"credentials": creds}
					}
					return nil
				},
				MockStatusUpdate: test.NewMockSubResourceUpdateFn(nil),
				MockCreate:       test.NewMockCreateFn(nil),
			},
			want: want{
				configuration: terraform.ProviderConfiguration{
					keyAPIKey: "only-key",
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			setupFn := TerraformSetupBuilder(tc.args.version, tc.args.providerSource, tc.args.providerVersion)
			ps, err := setupFn(context.Background(), tc.kube, tc.mg)

			if tc.want.err != nil {
				if err == nil {
					t.Errorf("\n%s\nTerraformSetupBuilder(...): expected error, got nil", tc.reason)
				}
				return
			}
			if err != nil {
				t.Fatalf("\n%s\nTerraformSetupBuilder(...): unexpected error: %v", tc.reason, err)
			}

			if diff := cmp.Diff(tc.want.configuration, ps.Configuration); diff != "" {
				t.Errorf("\n%s\nTerraformSetupBuilder(...): -want configuration, +got:\n%s", tc.reason, diff)
			}
		})
	}
}

// errAsNonNil is a sentinel used to indicate we expect a non-nil error
// without asserting its exact message.
var errAsNonNil = errSentinel("any")

type errSentinel string

func (e errSentinel) Error() string { return string(e) }

func TestLegacyToModernProviderConfigSpec(t *testing.T) {
	cases := map[string]struct {
		reason string
		pc     *clusterv1beta1.ProviderConfig
		want   *namespacedv1beta1.ProviderConfigSpec
	}{
		"NilInput": {
			reason: "Should return nil when input is nil.",
			pc:     nil,
			want:   nil,
		},
		"ConvertSecretSource": {
			reason: "Should convert cluster ProviderConfig spec to namespaced spec.",
			pc: &clusterv1beta1.ProviderConfig{
				Spec: clusterv1beta1.ProviderConfigSpec{
					Credentials: clusterv1beta1.ProviderCredentials{
						Source: xpv1.CredentialsSourceSecret,
						CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
							SecretRef: &xpv1.SecretKeySelector{
								SecretReference: xpv1.SecretReference{
									Name:      "my-secret",
									Namespace: "default",
								},
								Key: "creds",
							},
						},
					},
				},
			},
			want: &namespacedv1beta1.ProviderConfigSpec{
				Credentials: namespacedv1beta1.ProviderCredentials{
					Source: xpv1.CredentialsSourceSecret,
					CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
						SecretRef: &xpv1.SecretKeySelector{
							SecretReference: xpv1.SecretReference{
								Name:      "my-secret",
								Namespace: "default",
							},
							Key: "creds",
						},
					},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := legacyToModernProviderConfigSpec(tc.pc)
			if err != nil {
				t.Fatalf("\n%s\nlegacyToModernProviderConfigSpec(...): unexpected error: %v", tc.reason, err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("\n%s\nlegacyToModernProviderConfigSpec(...): -want, +got:\n%s", tc.reason, diff)
			}
		})
	}
}

func TestEnrichLocalSecretRefs(t *testing.T) {
	cases := map[string]struct {
		reason    string
		pc        *namespacedv1beta1.ProviderConfig
		mg        xpresource.Managed
		wantNS    string
		expectNil bool
	}{
		"SetsNamespaceFromManagedResource": {
			reason: "Should set the SecretRef namespace from the managed resource's namespace.",
			pc: &namespacedv1beta1.ProviderConfig{
				Spec: namespacedv1beta1.ProviderConfigSpec{
					Credentials: namespacedv1beta1.ProviderCredentials{
						Source: xpv1.CredentialsSourceSecret,
						CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
							SecretRef: &xpv1.SecretKeySelector{
								SecretReference: xpv1.SecretReference{
									Name: "my-secret",
								},
								Key: "creds",
							},
						},
					},
				},
			},
			mg: &fake.ModernManaged{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-mr",
					Namespace: "team-a",
				},
			},
			wantNS: "team-a",
		},
		"NilSecretRef": {
			reason: "Should be a no-op when SecretRef is nil.",
			pc: &namespacedv1beta1.ProviderConfig{
				Spec: namespacedv1beta1.ProviderConfigSpec{
					Credentials: namespacedv1beta1.ProviderCredentials{
						Source: xpv1.CredentialsSourceNone,
					},
				},
			},
			mg: &fake.ModernManaged{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-mr",
					Namespace: "team-a",
				},
			},
			expectNil: true,
		},
		"NilProviderConfig": {
			reason: "Should be a no-op when ProviderConfig is nil.",
			pc:     nil,
			mg: &fake.ModernManaged{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-mr",
					Namespace: "team-a",
				},
			},
			expectNil: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			enrichLocalSecretRefs(tc.pc, tc.mg)
			if tc.expectNil {
				return
			}
			if got := tc.pc.Spec.Credentials.SecretRef.Namespace; got != tc.wantNS {
				t.Errorf("\n%s\nenrichLocalSecretRefs(...): namespace = %q, want %q", tc.reason, got, tc.wantNS)
			}
		})
	}
}
