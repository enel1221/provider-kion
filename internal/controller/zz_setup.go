// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/pkg/controller"

	armtemplate "github.com/enel1221/provider-kion/internal/controller/kion/armtemplate"
	awsaccount "github.com/enel1221/provider-kion/internal/controller/kion/awsaccount"
	awsiampolicy "github.com/enel1221/provider-kion/internal/controller/kion/awsiampolicy"
	azureaccount "github.com/enel1221/provider-kion/internal/controller/kion/azureaccount"
	azurepolicy "github.com/enel1221/provider-kion/internal/controller/kion/azurepolicy"
	azurerole "github.com/enel1221/provider-kion/internal/controller/kion/azurerole"
	cloudformationtemplate "github.com/enel1221/provider-kion/internal/controller/kion/cloudformationtemplate"
	cloudrule "github.com/enel1221/provider-kion/internal/controller/kion/cloudrule"
	compliancecheck "github.com/enel1221/provider-kion/internal/controller/kion/compliancecheck"
	compliancestandard "github.com/enel1221/provider-kion/internal/controller/kion/compliancestandard"
	fundingsource "github.com/enel1221/provider-kion/internal/controller/kion/fundingsource"
	gcpaccount "github.com/enel1221/provider-kion/internal/controller/kion/gcpaccount"
	gcpiamrole "github.com/enel1221/provider-kion/internal/controller/kion/gcpiamrole"
	globalpermissionmapping "github.com/enel1221/provider-kion/internal/controller/kion/globalpermissionmapping"
	group "github.com/enel1221/provider-kion/internal/controller/kion/group"
	groupassociation "github.com/enel1221/provider-kion/internal/controller/kion/groupassociation"
	label "github.com/enel1221/provider-kion/internal/controller/kion/label"
	ou "github.com/enel1221/provider-kion/internal/controller/kion/ou"
	oucloudaccessrole "github.com/enel1221/provider-kion/internal/controller/kion/oucloudaccessrole"
	oupermissionmapping "github.com/enel1221/provider-kion/internal/controller/kion/oupermissionmapping"
	project "github.com/enel1221/provider-kion/internal/controller/kion/project"
	projectcloudaccessrole "github.com/enel1221/provider-kion/internal/controller/kion/projectcloudaccessrole"
	projectenforcement "github.com/enel1221/provider-kion/internal/controller/kion/projectenforcement"
	projectpermissionmapping "github.com/enel1221/provider-kion/internal/controller/kion/projectpermissionmapping"
	servicecontrolpolicy "github.com/enel1221/provider-kion/internal/controller/kion/servicecontrolpolicy"
	sourcepermissionmapping "github.com/enel1221/provider-kion/internal/controller/kion/sourcepermissionmapping"
	user "github.com/enel1221/provider-kion/internal/controller/kion/user"
	webhook "github.com/enel1221/provider-kion/internal/controller/kion/webhook"
	providerconfig "github.com/enel1221/provider-kion/internal/controller/providerconfig"
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
