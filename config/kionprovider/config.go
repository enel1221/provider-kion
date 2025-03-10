package kionprovider

import (
	"github.com/crossplane/upjet/pkg/config"
)

const shortGroup = "kion"

// Configure registers resource configurators for all supported Kion resources with the provided provider.
// It assigns a default short group to each resource to assist with grouping and identification within the provider.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("kion_aws_account", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		// r.ExternalName = config.IdentifierFromProvider
		// r.References["some_field"] = config.Reference{ ... }

		// r.References["account"] = config.Reference{
		// 	TerraformName: "kion_aws_account",
		// }
		r.Kind = "AWSAccount"
	})

	p.AddResourceConfigurator("kion_aws_cloudformation_template", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})

	p.AddResourceConfigurator("kion_aws_iam_policy", func(r *config.Resource) {
		r.ShortGroup = shortGroup

		r.Kind = "AWSIAMPolicy"
	})

	p.AddResourceConfigurator("kion_azure_account", func(r *config.Resource) {
		r.ShortGroup = shortGroup

		r.Kind = "AzureAccount"
	})

	p.AddResourceConfigurator("kion_azure_arm_template", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})

	p.AddResourceConfigurator("kion_azure_policy", func(r *config.Resource) {
		r.ShortGroup = shortGroup

		r.Kind = "AzurePolicy"
	})

	p.AddResourceConfigurator("kion_azure_role", func(r *config.Resource) {
		r.ShortGroup = shortGroup

		r.Kind = "AzureRole"
	})

	p.AddResourceConfigurator("kion_cloud_rule", func(r *config.Resource) {
		r.ShortGroup = shortGroup

		r.Kind = "CloudRule"
	})

	p.AddResourceConfigurator("kion_compliance_check", func(r *config.Resource) {
		r.ShortGroup = shortGroup

		r.Kind = "ComplianceCheck"
	})

	p.AddResourceConfigurator("kion_compliance_standard", func(r *config.Resource) {
		r.ShortGroup = shortGroup

		r.Kind = "ComplianceStandard"
	})

	p.AddResourceConfigurator("kion_funding_source", func(r *config.Resource) {
		r.ShortGroup = shortGroup

		r.Kind = "FundingSource"
	})

	p.AddResourceConfigurator("kion_funding_source_permission_mapping", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})

	p.AddResourceConfigurator("kion_gcp_account", func(r *config.Resource) {
		r.ShortGroup = shortGroup

		r.Kind = "GCPAccount"
	})

	p.AddResourceConfigurator("kion_gcp_iam_role", func(r *config.Resource) {
		r.ShortGroup = shortGroup

		r.Kind = "GCPIAMRole"
	})

	p.AddResourceConfigurator("kion_global_permission_mapping", func(r *config.Resource) {
		r.ShortGroup = shortGroup

		r.Kind = "GlobalPermissionMapping"
	})

	p.AddResourceConfigurator("kion_label", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})

	p.AddResourceConfigurator("kion_ou", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})

	p.AddResourceConfigurator("kion_ou_cloud_access_role", func(r *config.Resource) {
		r.ShortGroup = shortGroup

		r.Kind = "OUCloudAccessRole"
	})

	p.AddResourceConfigurator("kion_ou_permission_mapping", func(r *config.Resource) {
		r.ShortGroup = shortGroup

		r.Kind = "OUPermissionMapping"
	})

	p.AddResourceConfigurator("kion_project", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})

	p.AddResourceConfigurator("kion_project_cloud_access_role", func(r *config.Resource) {
		r.ShortGroup = shortGroup

		r.Kind = "ProjectCloudAccessRole"
	})

	p.AddResourceConfigurator("kion_project_enforcement", func(r *config.Resource) {
		r.ShortGroup = shortGroup

		r.Kind = "ProjectEnforcement"
	})

	p.AddResourceConfigurator("kion_project_permission_mapping", func(r *config.Resource) {
		r.ShortGroup = shortGroup

		r.Kind = "ProjectPermissionMapping"
	})

	p.AddResourceConfigurator("kion_saml_group_association", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})

	p.AddResourceConfigurator("kion_service_control_policy", func(r *config.Resource) {
		r.ShortGroup = shortGroup

		r.Kind = "ServiceControlPolicy"
	})

	p.AddResourceConfigurator("kion_user", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})

	p.AddResourceConfigurator("kion_user_group", func(r *config.Resource) {
		r.ShortGroup = shortGroup

		// Initialize SchemaElementOptions if nil
		if r.SchemaElementOptions == nil {
			r.SchemaElementOptions = make(config.SchemaElementOptions)
		}

		// Ensure owner_users and users are lists, not embedded objects
		r.SchemaElementOptions["owner_users"] = &config.SchemaElementOption{
			EmbeddedObject: false, // Default, but explicit for clarity
		}
		r.SchemaElementOptions["users"] = &config.SchemaElementOption{
			EmbeddedObject: false,
		}
	})

	p.AddResourceConfigurator("kion_webhook", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
}
