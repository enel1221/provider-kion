package config

import "github.com/crossplane/upjet/pkg/config"

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

func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

func ExternalNameConfigured() []string {
	l := make([]string, 0, len(ExternalNameConfigs))
	for name := range ExternalNameConfigs {
		// $ is added to match exactly the string (regex).
		l = append(l, name+"$")
	}
	return l
}
