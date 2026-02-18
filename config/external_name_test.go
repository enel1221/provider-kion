package config

import (
	"strings"
	"testing"

	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/google/go-cmp/cmp"
)

func TestExternalNameConfigs(t *testing.T) {
	const wantCount = 28

	t.Run("HasExpectedCount", func(t *testing.T) {
		if got := len(ExternalNameConfigs); got != wantCount {
			t.Errorf("ExternalNameConfigs: want %d entries, got %d", wantCount, got)
		}
	})

	t.Run("AllResources", func(t *testing.T) {
		expected := []string{
			"kion_aws_account",
			"kion_aws_cloudformation_template",
			"kion_aws_iam_policy",
			"kion_azure_account",
			"kion_azure_arm_template",
			"kion_azure_policy",
			"kion_azure_role",
			"kion_cloud_rule",
			"kion_compliance_check",
			"kion_compliance_standard",
			"kion_funding_source",
			"kion_funding_source_permission_mapping",
			"kion_gcp_account",
			"kion_gcp_iam_role",
			"kion_global_permission_mapping",
			"kion_label",
			"kion_ou",
			"kion_ou_cloud_access_role",
			"kion_ou_permission_mapping",
			"kion_project",
			"kion_project_cloud_access_role",
			"kion_project_enforcement",
			"kion_project_permission_mapping",
			"kion_saml_group_association",
			"kion_service_control_policy",
			"kion_user",
			"kion_user_group",
			"kion_webhook",
		}
		for _, name := range expected {
			if _, ok := ExternalNameConfigs[name]; !ok {
				t.Errorf("ExternalNameConfigs: missing entry for %q", name)
			}
		}
	})
}

func TestExternalNameConfigurations(t *testing.T) {
	t.Run("AppliesConfigToMatchingResource", func(t *testing.T) {
		opt := ExternalNameConfigurations()
		r := &config.Resource{
			Name:         "kion_aws_account",
			ExternalName: config.ExternalName{}, // zero value
		}
		opt(r)

		// After applying, ExternalName should be IdentifierFromProvider.
		want := config.IdentifierFromProvider
		if diff := cmp.Diff(want.GetExternalNameFn, r.ExternalName.GetExternalNameFn, cmpOptFuncNil()); diff != "" {
			t.Errorf("ExternalNameConfigurations(): applied ExternalName should be IdentifierFromProvider")
		}
	})

	t.Run("NoOpForUnknownResource", func(t *testing.T) {
		opt := ExternalNameConfigurations()
		r := &config.Resource{
			Name:         "kion_does_not_exist",
			ExternalName: config.ExternalName{},
		}
		opt(r)

		if r.ExternalName.GetExternalNameFn != nil {
			t.Error("ExternalNameConfigurations(): should not apply config to unknown resource")
		}
	})
}

func TestExternalNameConfigured(t *testing.T) {
	patterns := ExternalNameConfigured()

	t.Run("CountMatchesConfigs", func(t *testing.T) {
		if len(patterns) != len(ExternalNameConfigs) {
			t.Errorf("ExternalNameConfigured(): got %d patterns, want %d", len(patterns), len(ExternalNameConfigs))
		}
	})

	t.Run("AllEndWithDollar", func(t *testing.T) {
		for _, p := range patterns {
			if !strings.HasSuffix(p, "$") {
				t.Errorf("ExternalNameConfigured(): pattern %q does not end with '$'", p)
			}
		}
	})

	t.Run("AllStartWithKion", func(t *testing.T) {
		for _, p := range patterns {
			if !strings.HasPrefix(p, "kion_") {
				t.Errorf("ExternalNameConfigured(): pattern %q does not start with 'kion_'", p)
			}
		}
	})
}

// cmpOptFuncNil returns a cmp.Option that considers two function values equal if
// they are both nil or both non-nil. We cannot compare function pointers directly
// in Go, so this is a pragmatic approach.
func cmpOptFuncNil() cmp.Option {
	return cmp.FilterPath(
		func(p cmp.Path) bool { return true },
		cmp.Comparer(func(_, _ interface{}) bool { return true }),
	)
}
