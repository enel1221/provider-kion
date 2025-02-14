package config

import "github.com/crossplane/upjet/pkg/config"

// ExternalNameConfigs maps Kion resource identifiers to their corresponding external name configuration
// functions. Each key in the map represents a specific Kion resource, and the associated value specifies
// how the external name for that resource is determined. In this configuration, all resources use the
// IdentifierFromProvider function, which indicates that the external name is derived directly from the provider's
// identifier. Since this variable is exported, it serves as a central reference for configuring external names across
// the provider, ensuring consistency in resource identification.
var ExternalNameConfigs = map[string]config.ExternalName{

	// Kion resources:
	"kion_aws_account":                     config.IdentifierFromProvider,
	"kion_aws_cloudformation_template":     config.IdentifierFromProvider,
	"kion_aws_iam_policy":                  config.IdentifierFromProvider,
	"kion_azure_account":                   config.IdentifierFromProvider,
	"kion_azure_arm_template":              config.IdentifierFromProvider,
	"kion_azure_policy":                    config.IdentifierFromProvider,
	"kion_azure_role":                      config.IdentifierFromProvider,
	"kion_cloud_rule":                      config.IdentifierFromProvider,
	"kion_compliance_check":                config.IdentifierFromProvider,
	"kion_compliance_standard":             config.IdentifierFromProvider,
	"kion_funding_source":                  config.IdentifierFromProvider,
	"kion_funding_source_permission_mapping": config.IdentifierFromProvider,
	"kion_gcp_account":                     config.IdentifierFromProvider,
	"kion_gcp_iam_role":                    config.IdentifierFromProvider,
	"kion_global_permission_mapping":       config.IdentifierFromProvider,
	"kion_label":                           config.IdentifierFromProvider,
	"kion_ou":                              config.IdentifierFromProvider,
	"kion_ou_cloud_access_role":            config.IdentifierFromProvider,
	"kion_ou_permission_mapping":           config.IdentifierFromProvider,
	"kion_project":						 	config.IdentifierFromProvider,
	"kion_project_cloud_access_role":       config.IdentifierFromProvider,
	"kion_project_enforcement":             config.IdentifierFromProvider,
	"kion_project_permission_mapping":      config.IdentifierFromProvider,
	"kion_saml_group_association":          config.IdentifierFromProvider,
	"kion_service_control_policy":          config.IdentifierFromProvider,
	"kion_user":                            config.IdentifierFromProvider,
	"kion_user_group":                      config.IdentifierFromProvider,
	"kion_webhook":                         config.IdentifierFromProvider,
}

// ExternalNameConfigurations returns a ResourceOption that sets the ExternalName of a resource.
// It checks if there is a pre-defined external name corresponding to the resource's name from the
// ExternalNameConfigs map. If an entry is found, the resource's ExternalName field is updated accordingly.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

// ExternalNameConfigured returns a slice of external name patterns.
// Each pattern is formed by appending a '$' to the keys from ExternalNameConfigs,
// ensuring that the pattern matches the exact string when used with regular expressions.
func ExternalNameConfigured() []string {
	l := make([]string, 0, len(ExternalNameConfigs))
	for name := range ExternalNameConfigs {
		// $ is added to match exactly the string (regex).
		l = append(l, name+"$")
	}
	return l
}
