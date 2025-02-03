// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/pkg/controller"

	account "github.com/enel1221/provider-kion/internal/controller/kion/account"
	armtemplate "github.com/enel1221/provider-kion/internal/controller/kion/armtemplate"
	check "github.com/enel1221/provider-kion/internal/controller/kion/check"
	cloudaccessrole "github.com/enel1221/provider-kion/internal/controller/kion/cloudaccessrole"
	cloudformationtemplate "github.com/enel1221/provider-kion/internal/controller/kion/cloudformationtemplate"
	controlpolicy "github.com/enel1221/provider-kion/internal/controller/kion/controlpolicy"
	enforcement "github.com/enel1221/provider-kion/internal/controller/kion/enforcement"
	group "github.com/enel1221/provider-kion/internal/controller/kion/group"
	groupassociation "github.com/enel1221/provider-kion/internal/controller/kion/groupassociation"
	iampolicy "github.com/enel1221/provider-kion/internal/controller/kion/iampolicy"
	iamrole "github.com/enel1221/provider-kion/internal/controller/kion/iamrole"
	label "github.com/enel1221/provider-kion/internal/controller/kion/label"
	ou "github.com/enel1221/provider-kion/internal/controller/kion/ou"
	permissionmapping "github.com/enel1221/provider-kion/internal/controller/kion/permissionmapping"
	policy "github.com/enel1221/provider-kion/internal/controller/kion/policy"
	project "github.com/enel1221/provider-kion/internal/controller/kion/project"
	role "github.com/enel1221/provider-kion/internal/controller/kion/role"
	rule "github.com/enel1221/provider-kion/internal/controller/kion/rule"
	source "github.com/enel1221/provider-kion/internal/controller/kion/source"
	sourcepermissionmapping "github.com/enel1221/provider-kion/internal/controller/kion/sourcepermissionmapping"
	standard "github.com/enel1221/provider-kion/internal/controller/kion/standard"
	user "github.com/enel1221/provider-kion/internal/controller/kion/user"
	webhook "github.com/enel1221/provider-kion/internal/controller/kion/webhook"
	providerconfig "github.com/enel1221/provider-kion/internal/controller/providerconfig"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		account.Setup,
		account.Setup,
		account.Setup,
		armtemplate.Setup,
		check.Setup,
		cloudaccessrole.Setup,
		cloudaccessrole.Setup,
		cloudformationtemplate.Setup,
		controlpolicy.Setup,
		enforcement.Setup,
		group.Setup,
		groupassociation.Setup,
		iampolicy.Setup,
		iamrole.Setup,
		label.Setup,
		ou.Setup,
		permissionmapping.Setup,
		permissionmapping.Setup,
		permissionmapping.Setup,
		policy.Setup,
		project.Setup,
		role.Setup,
		rule.Setup,
		source.Setup,
		sourcepermissionmapping.Setup,
		standard.Setup,
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
