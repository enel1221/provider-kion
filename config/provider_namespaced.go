package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"

	"github.com/enel1221/provider-kion/config/kionprovider"
)

// GetNamespacedProvider returns the namespaced provider configuration
func GetNamespacedProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider(
		[]byte(providerSchema),
		resourcePrefix,
		modulePath,
		[]byte(providerMetadata),

		ujconfig.WithRootGroup("m.upbound.io"),
		ujconfig.WithShortName("kion"),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		),
		ujconfig.WithBasePackages(ujconfig.BasePackages{
			APIVersion: []string{
				"v1beta1",
			},
			ControllerMap: map[string]string{
				"providerconfig": ujconfig.PackageNameConfig,
			},
		}),
	)

	// Register resource configs — shared with cluster scope
	kionprovider.Configure(pc)

	pc.ConfigureResources()

	return pc
}
