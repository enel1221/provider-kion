package kionprovider

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

const shortGroup = "kion"

// Configure registers resource configurators for all supported Kion resources with the provided provider.
// It assigns a default short group to each resource to assist with grouping and identification within the provider.
func Configure(p *config.Provider) {
	// ── Accounts ─────────────────────────────────────────────────────

	p.AddResourceConfigurator("kion_aws_account", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "AWSAccount"
		r.UseAsync = true
		r.References["project_id"] = config.Reference{
			TerraformName: "kion_project",
		}
	})

	p.AddResourceConfigurator("kion_azure_account", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "AzureAccount"
		r.References["project_id"] = config.Reference{
			TerraformName: "kion_project",
		}
	})

	p.AddResourceConfigurator("kion_gcp_account", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "GCPAccount"
		r.References["project_id"] = config.Reference{
			TerraformName: "kion_project",
		}
	})

	// ── Templates & Policies ─────────────────────────────────────────

	p.AddResourceConfigurator("kion_aws_cloudformation_template", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})

	p.AddResourceConfigurator("kion_aws_iam_policy", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "AWSIAMPolicy"
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

	p.AddResourceConfigurator("kion_service_control_policy", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "ServiceControlPolicy"
	})

	p.AddResourceConfigurator("kion_gcp_iam_role", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "GCPIAMRole"
	})

	// ── Cloud Rule ───────────────────────────────────────────────────

	p.AddResourceConfigurator("kion_cloud_rule", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "CloudRule"

		// Compliance & template references (nested block .id fields)
		r.References["compliance_standards.id"] = config.Reference{
			TerraformName: "kion_compliance_standard",
		}
		r.References["aws_cloudformation_templates.id"] = config.Reference{
			TerraformName: "kion_aws_cloudformation_template",
		}
		r.References["aws_iam_policies.id"] = config.Reference{
			TerraformName: "kion_aws_iam_policy",
		}
		r.References["azure_arm_template_definitions.id"] = config.Reference{
			TerraformName: "kion_azure_arm_template",
		}
		r.References["azure_policy_definitions.id"] = config.Reference{
			TerraformName: "kion_azure_policy",
		}
		r.References["azure_role_definitions.id"] = config.Reference{
			TerraformName: "kion_azure_role",
		}
		r.References["service_control_policies.id"] = config.Reference{
			TerraformName: "kion_service_control_policy",
		}
		r.References["gcp_iam_roles.id"] = config.Reference{
			TerraformName: "kion_gcp_iam_role",
		}

		// Scope references
		r.References["ous.id"] = config.Reference{
			TerraformName: "kion_ou",
		}
		r.References["projects.id"] = config.Reference{
			TerraformName: "kion_project",
		}

		// Webhook references
		r.References["pre_webhook_id"] = config.Reference{
			TerraformName: "kion_webhook",
		}
		r.References["post_webhook_id"] = config.Reference{
			TerraformName: "kion_webhook",
		}

		// Owner references
		r.References["owner_user_groups.id"] = config.Reference{
			TerraformName: "kion_user_group",
		}
		r.References["owner_users.id"] = config.Reference{
			TerraformName: "kion_user",
		}
	})

	// ── Compliance ───────────────────────────────────────────────────

	p.AddResourceConfigurator("kion_compliance_check", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "ComplianceCheck"
	})

	p.AddResourceConfigurator("kion_compliance_standard", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "ComplianceStandard"
	})

	// ── Funding Source ───────────────────────────────────────────────

	p.AddResourceConfigurator("kion_funding_source", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "FundingSource"
	})

	// ── OU ───────────────────────────────────────────────────────────

	p.AddResourceConfigurator("kion_ou", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.References["parent_ou_id"] = config.Reference{
			TerraformName: "kion_ou",
		}
	})

	// ── Project ──────────────────────────────────────────────────────

	p.AddResourceConfigurator("kion_project", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.References["ou_id"] = config.Reference{
			TerraformName: "kion_ou",
		}
	})

	// ── Cloud Access Roles ───────────────────────────────────────────

	p.AddResourceConfigurator("kion_ou_cloud_access_role", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "OUCloudAccessRole"
		r.References["ou_id"] = config.Reference{
			TerraformName: "kion_ou",
		}
		r.References["azure_role_definitions.id"] = config.Reference{
			TerraformName: "kion_azure_role",
		}
		r.References["gcp_iam_roles.id"] = config.Reference{
			TerraformName: "kion_gcp_iam_role",
		}
		r.References["user_groups.id"] = config.Reference{
			TerraformName: "kion_user_group",
		}
	})

	p.AddResourceConfigurator("kion_project_cloud_access_role", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "ProjectCloudAccessRole"
		r.References["project_id"] = config.Reference{
			TerraformName: "kion_project",
		}
		r.References["azure_role_definitions.id"] = config.Reference{
			TerraformName: "kion_azure_role",
		}
		r.References["gcp_iam_roles.id"] = config.Reference{
			TerraformName: "kion_gcp_iam_role",
		}
		r.References["user_groups.id"] = config.Reference{
			TerraformName: "kion_user_group",
		}
	})

	// ── Project Enforcement ──────────────────────────────────────────

	p.AddResourceConfigurator("kion_project_enforcement", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "ProjectEnforcement"
		r.References["project_id"] = config.Reference{
			TerraformName: "kion_project",
		}
		r.References["cloud_rule_id"] = config.Reference{
			TerraformName: "kion_cloud_rule",
		}
	})

	// ── Permission Mappings ──────────────────────────────────────────

	p.AddResourceConfigurator("kion_ou_permission_mapping", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "OUPermissionMapping"
		r.References["ou_id"] = config.Reference{
			TerraformName: "kion_ou",
		}
		r.References["user_groups_ids"] = config.Reference{
			TerraformName: "kion_user_group",
		}
		r.References["user_ids"] = config.Reference{
			TerraformName: "kion_user",
		}
	})

	p.AddResourceConfigurator("kion_project_permission_mapping", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "ProjectPermissionMapping"
		r.References["project_id"] = config.Reference{
			TerraformName: "kion_project",
		}
		r.References["user_groups_ids"] = config.Reference{
			TerraformName: "kion_user_group",
		}
		r.References["user_ids"] = config.Reference{
			TerraformName: "kion_user",
		}
	})

	p.AddResourceConfigurator("kion_funding_source_permission_mapping", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "FundingSourcePermissionMapping"
		r.References["funding_source_id"] = config.Reference{
			TerraformName: "kion_funding_source",
		}
		r.References["user_groups_ids"] = config.Reference{
			TerraformName: "kion_user_group",
		}
		r.References["user_ids"] = config.Reference{
			TerraformName: "kion_user",
		}
	})

	p.AddResourceConfigurator("kion_global_permission_mapping", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "GlobalPermissionMapping"
		r.References["user_groups_ids"] = config.Reference{
			TerraformName: "kion_user_group",
		}
		r.References["user_ids"] = config.Reference{
			TerraformName: "kion_user",
		}
	})

	// ── SAML Group Association ───────────────────────────────────────

	p.AddResourceConfigurator("kion_saml_group_association", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.References["user_group_id"] = config.Reference{
			TerraformName: "kion_user_group",
		}
	})

	// ── User & Group ─────────────────────────────────────────────────

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
			EmbeddedObject: false,
		}
		r.SchemaElementOptions["users"] = &config.SchemaElementOption{
			EmbeddedObject: false,
		}
	})

	// ── Label & Webhook ──────────────────────────────────────────────

	p.AddResourceConfigurator("kion_label", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})

	p.AddResourceConfigurator("kion_webhook", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
}
