package kionprovider

import (
	"github.com/crossplane/upjet/pkg/config"
)

func Configure(p *config.Provider) {
	p.AddResourceConfigurator("kion_aws_account", func(r *config.Resource) {
		r.ShortGroup = "kion"
		// r.ExternalName = config.IdentifierFromProvider
		// r.References["some_field"] = config.Reference{ ... }
	})

	p.AddResourceConfigurator("kion_aws_cloudformation_template", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_aws_iam_policy", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_azure_account", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_azure_arm_template", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_azure_policy", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_azure_role", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_cloud_rule", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_compliance_check", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_compliance_standard", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_funding_source", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_funding_source_permission_mapping", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_gcp_account", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_gcp_iam_role", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_global_permission_mapping", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_label", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_ou", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_ou_cloud_access_role", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_ou_permission_mapping", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_project", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_project_cloud_access_role", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_project_enforcement", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_project_permission_mapping", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_saml_group_association", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_service_control_policy", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_user", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_user_group", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

	p.AddResourceConfigurator("kion_webhook", func(r *config.Resource) {
		r.ShortGroup = "kion"
	})

}
