// patch-xvalidation restores XValidation "required parameter" CEL markers
// on generated Upjet type structs whose Terraform-required fields lost their
// CRD validation when Upjet reference generation made them schema-optional.
//
// This program is invoked as a go:generate step between the Upjet generator
// and controller-gen so that the markers are present when CRD manifests are
// emitted. It is idempotent: running it twice produces the same output.
//
// Usage:
//
//	go run ./hack/patch-xvalidation <apis-root>
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type rule struct {
	field       string
	includeInit bool
	// refSuffix is "Ref" for singular reference fields or "Refs" for
	// list-typed reference fields (e.g. userGroupsIdsRefs).
	refSuffix string
}

type patch struct {
	file       string
	structName string
	rules      []rule
}

var patches = []patch{
	{
		file: "zz_projectenforcement_types.go", structName: "ProjectEnforcement",
		rules: []rule{{field: "projectId", includeInit: true, refSuffix: "Ref"}},
	},
	{
		file: "zz_oupermissionmapping_types.go", structName: "OUPermissionMapping",
		rules: []rule{
			{field: "ouId", includeInit: true, refSuffix: "Ref"},
			{field: "userGroupsIds", includeInit: true, refSuffix: "Refs"},
			{field: "userIds", includeInit: true, refSuffix: "Refs"},
		},
	},
	{
		file: "zz_projectpermissionmapping_types.go", structName: "ProjectPermissionMapping",
		rules: []rule{
			{field: "projectId", includeInit: true, refSuffix: "Ref"},
			{field: "userGroupsIds", includeInit: true, refSuffix: "Refs"},
			{field: "userIds", includeInit: true, refSuffix: "Refs"},
		},
	},
	{
		file: "zz_fundingsourcepermissionmapping_types.go", structName: "FundingSourcePermissionMapping",
		rules: []rule{
			{field: "fundingSourceId", includeInit: true, refSuffix: "Ref"},
			{field: "userGroupsIds", includeInit: true, refSuffix: "Refs"},
			{field: "userIds", includeInit: true, refSuffix: "Refs"},
		},
	},
	{
		file: "zz_globalpermissionmapping_types.go", structName: "GlobalPermissionMapping",
		rules: []rule{
			{field: "userGroupsIds", includeInit: true, refSuffix: "Refs"},
			{field: "userIds", includeInit: true, refSuffix: "Refs"},
		},
	},
	{
		file: "zz_customvariableoverride_types.go", structName: "CustomVariableOverride",
		rules: []rule{{field: "customVariableId", includeInit: true, refSuffix: "Ref"}},
	},
	{
		file: "zz_projectnote_types.go", structName: "ProjectNote",
		rules: []rule{{field: "projectId", includeInit: true, refSuffix: "Ref"}},
	},
	{
		file: "zz_oucloudaccessrole_types.go", structName: "OUCloudAccessRole",
		rules: []rule{{field: "ouId", includeInit: true, refSuffix: "Ref"}},
	},
	{
		file: "zz_projectcloudaccessrole_types.go", structName: "ProjectCloudAccessRole",
		rules: []rule{{field: "projectId", includeInit: true, refSuffix: "Ref"}},
	},
	{
		file: "zz_ou_types.go", structName: "Ou",
		rules: []rule{{field: "parentOuId", includeInit: true, refSuffix: "Ref"}},
	},
	{
		file: "zz_project_types.go", structName: "Project",
		rules: []rule{{field: "ouId", includeInit: true, refSuffix: "Ref"}},
	},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: patch-xvalidation <apis-root>\n")
		os.Exit(1)
	}
	root := os.Args[1]
	scopes := []string{"cluster", "namespaced"}
	var errors []string

	for _, scope := range scopes {
		dir := filepath.Join(root, scope, "kion", "v1alpha1")
		for _, p := range patches {
			path := filepath.Join(dir, p.file)
			if err := applyPatch(path, p); err != nil {
				errors = append(errors, fmt.Sprintf("%s/%s: %v", scope, p.file, err))
			}
		}
	}

	if len(errors) > 0 {
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", e)
		}
		os.Exit(1)
	}
}

func markerLine(r rule) string {
	refField := r.field + r.refSuffix
	selField := r.field + "Selector"
	if r.includeInit {
		return fmt.Sprintf(
			"\t// +kubebuilder:validation:XValidation:rule=\"!('*' in self.managementPolicies || 'Create' in self.managementPolicies || 'Update' in self.managementPolicies) || has(self.forProvider.%s) || has(self.forProvider.%s) || has(self.forProvider.%s) || (has(self.initProvider) && has(self.initProvider.%s))\",message=\"spec.forProvider.%s is a required parameter\"",
			r.field, refField, selField, r.field, r.field,
		)
	}
	return fmt.Sprintf(
		"\t// +kubebuilder:validation:XValidation:rule=\"!('*' in self.managementPolicies || 'Create' in self.managementPolicies || 'Update' in self.managementPolicies) || has(self.forProvider.%s) || has(self.forProvider.%s) || has(self.forProvider.%s)\",message=\"spec.forProvider.%s is a required parameter\"",
		r.field, refField, selField, r.field,
	)
}

func applyPatch(path string, p patch) error {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	lines := strings.Split(string(data), "\n")

	specLineIdx := findSpecLine(lines, p.structName)
	if specLineIdx < 0 {
		return fmt.Errorf("could not find Spec field in struct %s", p.structName)
	}

	var toInsert []string
	for _, r := range p.rules {
		marker := markerLine(r)
		if !hasMarker(lines, specLineIdx, marker) {
			toInsert = append(toInsert, marker)
		}
	}

	if len(toInsert) == 0 {
		return nil
	}

	newLines := make([]string, 0, len(lines)+len(toInsert))
	newLines = append(newLines, lines[:specLineIdx]...)
	newLines = append(newLines, toInsert...)
	newLines = append(newLines, lines[specLineIdx:]...)

	return os.WriteFile(filepath.Clean(path), []byte(strings.Join(newLines, "\n")), 0600)
}

func findSpecLine(lines []string, structName string) int {
	structDecl := fmt.Sprintf("type %s struct {", structName)
	inStruct := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == structDecl {
			inStruct = true
			continue
		}
		if inStruct {
			if strings.HasSuffix(trimmed, "`json:\"spec\"`") {
				return i
			}
			if trimmed == "}" {
				return -1
			}
		}
	}
	return -1
}

func hasMarker(lines []string, specLineIdx int, marker string) bool {
	for i := specLineIdx - 1; i >= 0 && i >= specLineIdx-30; i-- {
		if strings.TrimSpace(lines[i]) == strings.TrimSpace(marker) {
			return true
		}
	}
	return false
}
