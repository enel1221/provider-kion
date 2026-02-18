// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	armtemplate "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/armtemplate"
	awsaccount "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/awsaccount"
	awsiampolicy "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/awsiampolicy"
	azureaccount "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/azureaccount"
	azurepolicy "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/azurepolicy"
	azurerole "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/azurerole"
	cloudformationtemplate "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/cloudformationtemplate"
	cloudrule "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/cloudrule"
	compliancecheck "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/compliancecheck"
	compliancestandard "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/compliancestandard"
	fundingsource "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/fundingsource"
	gcpaccount "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/gcpaccount"
	gcpiamrole "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/gcpiamrole"
	globalpermissionmapping "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/globalpermissionmapping"
	group "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/group"
	groupassociation "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/groupassociation"
	label "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/label"
	ou "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/ou"
	oucloudaccessrole "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/oucloudaccessrole"
	oupermissionmapping "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/oupermissionmapping"
	project "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/project"
	projectcloudaccessrole "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/projectcloudaccessrole"
	projectenforcement "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/projectenforcement"
	projectpermissionmapping "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/projectpermissionmapping"
	servicecontrolpolicy "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/servicecontrolpolicy"
	sourcepermissionmapping "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/sourcepermissionmapping"
	user "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/user"
	webhook "github.com/enel1221/provider-kion/internal/controller/namespaced/kion/webhook"
	providerconfig "github.com/enel1221/provider-kion/internal/controller/namespaced/providerconfig"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		armtemplate.Setup,
		awsaccount.Setup,
		awsiampolicy.Setup,
		azureaccount.Setup,
		azurepolicy.Setup,
		azurerole.Setup,
		cloudformationtemplate.Setup,
		cloudrule.Setup,
		compliancecheck.Setup,
		compliancestandard.Setup,
		fundingsource.Setup,
		gcpaccount.Setup,
		gcpiamrole.Setup,
		globalpermissionmapping.Setup,
		group.Setup,
		groupassociation.Setup,
		label.Setup,
		ou.Setup,
		oucloudaccessrole.Setup,
		oupermissionmapping.Setup,
		project.Setup,
		projectcloudaccessrole.Setup,
		projectenforcement.Setup,
		projectpermissionmapping.Setup,
		servicecontrolpolicy.Setup,
		sourcepermissionmapping.Setup,
		user.Setup,
		webhook.Setup,
		providerconfig.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		armtemplate.SetupGated,
		awsaccount.SetupGated,
		awsiampolicy.SetupGated,
		azureaccount.SetupGated,
		azurepolicy.SetupGated,
		azurerole.SetupGated,
		cloudformationtemplate.SetupGated,
		cloudrule.SetupGated,
		compliancecheck.SetupGated,
		compliancestandard.SetupGated,
		fundingsource.SetupGated,
		gcpaccount.SetupGated,
		gcpiamrole.SetupGated,
		globalpermissionmapping.SetupGated,
		group.SetupGated,
		groupassociation.SetupGated,
		label.SetupGated,
		ou.SetupGated,
		oucloudaccessrole.SetupGated,
		oupermissionmapping.SetupGated,
		project.SetupGated,
		projectcloudaccessrole.SetupGated,
		projectenforcement.SetupGated,
		projectpermissionmapping.SetupGated,
		servicecontrolpolicy.SetupGated,
		sourcepermissionmapping.SetupGated,
		user.SetupGated,
		webhook.SetupGated,
		providerconfig.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
