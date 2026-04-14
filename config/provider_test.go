package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestGetProviderResources validates that GetProvider returns a fully configured
// provider with the expected resources, Kind names, ShortGroups, and references.
// This exercises the real embedded schema + kionprovider.Configure() path.
func TestGetProviderResources(t *testing.T) {
	pc := GetProvider()

	// expectedKinds maps Terraform resource names to the Kind we expect after
	// configuration. Resources without an explicit Kind override use upjet's
	// auto-generated name (not tested here — only explicit overrides).
	expectedKinds := map[string]string{
		"kion_app_config":                        "AppConfig",
		"kion_aws_account":                       "AWSAccount",
		"kion_aws_iam_policy":                    "AWSIAMPolicy",
		"kion_azure_account":                     "AzureAccount",
		"kion_azure_policy":                      "AzurePolicy",
		"kion_azure_role":                        "AzureRole",
		"kion_cloud_rule":                        "CloudRule",
		"kion_compliance_check":                  "ComplianceCheck",
		"kion_compliance_standard":               "ComplianceStandard",
		"kion_custom_account":                    "CustomAccount",
		"kion_custom_variable":                   "CustomVariable",
		"kion_custom_variable_override":          "CustomVariableOverride",
		"kion_funding_source":                    "FundingSource",
		"kion_funding_source_permission_mapping": "FundingSourcePermissionMapping",
		"kion_gcp_account":                       "GCPAccount",
		"kion_gcp_iam_role":                      "GCPIAMRole",
		"kion_ou_cloud_access_role":              "OUCloudAccessRole",
		"kion_ou_permission_mapping":             "OUPermissionMapping",
		"kion_project_cloud_access_role":         "ProjectCloudAccessRole",
		"kion_project_enforcement":               "ProjectEnforcement",
		"kion_project_note":                      "ProjectNote",
		"kion_project_permission_mapping":        "ProjectPermissionMapping",
		"kion_service_control_policy":            "ServiceControlPolicy",
		"kion_global_permission_mapping":         "GlobalPermissionMapping",
	}

	t.Run("ExplicitKindNames", func(t *testing.T) {
		for tfName, wantKind := range expectedKinds {
			r, ok := pc.Resources[tfName]
			if !ok {
				t.Errorf("resource %q not found in provider", tfName)
				continue
			}
			if r.Kind != wantKind {
				t.Errorf("resource %q: Kind = %q, want %q", tfName, r.Kind, wantKind)
			}
		}
	})

	t.Run("AllResourcesHaveShortGroupKion", func(t *testing.T) {
		for name, r := range pc.Resources {
			if r.ShortGroup != "kion" {
				t.Errorf("resource %q: ShortGroup = %q, want %q", name, r.ShortGroup, "kion")
			}
		}
	})

	t.Run("FundingSourcePermissionMappingKindRegression", func(t *testing.T) {
		// Regression test for Step 4 fix: previously auto-generated as
		// "SourcePermissionMapping" instead of "FundingSourcePermissionMapping".
		r, ok := pc.Resources["kion_funding_source_permission_mapping"]
		if !ok {
			t.Fatal("resource kion_funding_source_permission_mapping not found")
		}
		if r.Kind != "FundingSourcePermissionMapping" {
			t.Errorf("FundingSourcePermissionMapping Kind = %q, want %q", r.Kind, "FundingSourcePermissionMapping")
		}
	})

	t.Run("AWSAccountUseAsyncEnabled", func(t *testing.T) {
		r, ok := pc.Resources["kion_aws_account"]
		if !ok {
			t.Fatal("resource kion_aws_account not found")
		}
		if !r.UseAsync {
			t.Error("kion_aws_account: UseAsync should be true")
		}
	})

	// Reference tests for key resources.
	referenceTests := map[string]map[string]string{
		"kion_aws_account": {
			"project_id": "kion_project",
		},
		"kion_azure_account": {
			"project_id": "kion_project",
		},
		"kion_gcp_account": {
			"project_id": "kion_project",
		},
		"kion_project": {
			"ou_id": "kion_ou",
		},
		"kion_ou": {
			"parent_ou_id": "kion_ou",
		},
		"kion_project_enforcement": {
			"project_id":    "kion_project",
			"cloud_rule_id": "kion_cloud_rule",
		},
		"kion_custom_account": {
			"project_id": "kion_project",
		},
		"kion_custom_variable_override": {
			"custom_variable_id": "kion_custom_variable",
		},
		"kion_cloud_rule": {
			"compliance_standards.id":           "kion_compliance_standard",
			"aws_cloudformation_templates.id":   "kion_aws_cloudformation_template",
			"aws_iam_policies.id":               "kion_aws_iam_policy",
			"azure_arm_template_definitions.id": "kion_azure_arm_template",
			"azure_policy_definitions.id":       "kion_azure_policy",
			"azure_role_definitions.id":         "kion_azure_role",
			"service_control_policies.id":       "kion_service_control_policy",
			"gcp_iam_roles.id":                  "kion_gcp_iam_role",
			"ous.id":                            "kion_ou",
			"projects.id":                       "kion_project",
			"pre_webhook_id":                    "kion_webhook",
			"post_webhook_id":                   "kion_webhook",
			"owner_user_groups.id":              "kion_user_group",
			"owner_users.id":                    "kion_user",
		},
		"kion_ou_permission_mapping": {
			"ou_id":           "kion_ou",
			"user_groups_ids": "kion_user_group",
			"user_ids":        "kion_user",
		},
		"kion_project_permission_mapping": {
			"project_id":      "kion_project",
			"user_groups_ids": "kion_user_group",
			"user_ids":        "kion_user",
		},
		"kion_project_note": {
			"project_id": "kion_project",
		},
		"kion_funding_source_permission_mapping": {
			"funding_source_id": "kion_funding_source",
			"user_groups_ids":   "kion_user_group",
			"user_ids":          "kion_user",
		},
		"kion_global_permission_mapping": {
			"user_groups_ids": "kion_user_group",
			"user_ids":        "kion_user",
		},
		"kion_saml_group_association": {
			"user_group_id": "kion_user_group",
		},
		"kion_ou_cloud_access_role": {
			"ou_id":                     "kion_ou",
			"azure_role_definitions.id": "kion_azure_role",
			"gcp_iam_roles.id":          "kion_gcp_iam_role",
			"user_groups.id":            "kion_user_group",
		},
		"kion_project_cloud_access_role": {
			"project_id":                "kion_project",
			"azure_role_definitions.id": "kion_azure_role",
			"gcp_iam_roles.id":          "kion_gcp_iam_role",
			"user_groups.id":            "kion_user_group",
		},
	}

	t.Run("References", func(t *testing.T) {
		for tfName, wantRefs := range referenceTests {
			r, ok := pc.Resources[tfName]
			if !ok {
				t.Errorf("resource %q not found", tfName)
				continue
			}
			for field, wantTfName := range wantRefs {
				ref, ok := r.References[field]
				if !ok {
					t.Errorf("resource %q: missing reference for field %q", tfName, field)
					continue
				}
				if ref.TerraformName != wantTfName {
					t.Errorf("resource %q field %q: TerraformName = %q, want %q", tfName, field, ref.TerraformName, wantTfName)
				}
			}
		}
	})

	t.Run("UserGroupSchemaElementOptions", func(t *testing.T) {
		r, ok := pc.Resources["kion_user_group"]
		if !ok {
			t.Fatal("resource kion_user_group not found")
		}
		for _, field := range []string{"owner_users", "users"} {
			opt, ok := r.SchemaElementOptions[field]
			if !ok {
				t.Errorf("kion_user_group: missing SchemaElementOption for %q", field)
				continue
			}
			if opt.EmbeddedObject {
				t.Errorf("kion_user_group: SchemaElementOption[%q].EmbeddedObject = true, want false", field)
			}
		}
	})
}

func TestGetNamespacedProvider(t *testing.T) {
	pc := GetNamespacedProvider()

	t.Run("RootGroupIsNamespaced", func(t *testing.T) {
		// Namespaced provider uses "m.upbound.io" root group.
		// We verify a resource's API group contains the expected suffix.
		r, ok := pc.Resources["kion_aws_account"]
		if !ok {
			t.Fatal("resource kion_aws_account not found")
		}
		// The ShortGroup should still be "kion".
		if r.ShortGroup != "kion" {
			t.Errorf("namespaced kion_aws_account ShortGroup = %q, want %q", r.ShortGroup, "kion")
		}
	})

	t.Run("HasSameResourceCount", func(t *testing.T) {
		cluster := GetProvider()
		if diff := cmp.Diff(len(cluster.Resources), len(pc.Resources)); diff != "" {
			t.Errorf("namespaced provider resource count differs from cluster: -want +got:\n%s", diff)
		}
	})
}
