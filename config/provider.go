package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	ujconfig "github.com/crossplane/upjet/pkg/config"

	"github.com/enel1221/provider-kion/config/kionprovider"
)

const (
	resourcePrefix = "kion"
	modulePath     = "github.com/enel1221/provider-kion"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns provider configuration
func GetProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider(
		[]byte(providerSchema), // provider's schema from Terraform
		resourcePrefix,         // e.g. "kion"
		modulePath,             // your module path
		[]byte(providerMetadata),

		// Optionally adjust these as desired:
		ujconfig.WithRootGroup("upbound.io"), // or "kion.crossplane.io"
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		),
	)

	// Register your single-file config, which configures *all* Kion resources/data sources
	kionprovider.Configure(pc)

	// Tells Upjet to finalize the resource configuration
	pc.ConfigureResources()

	return pc
}
