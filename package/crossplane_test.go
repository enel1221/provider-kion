package package_test

import (
	"os"
	"regexp"
	"strconv"
	"testing"

	"gopkg.in/yaml.v3"
)

type providerPackage struct {
	Spec providerSpec `yaml:"spec"`
}

type providerSpec struct {
	Capabilities []string `yaml:"capabilities"`
}

func TestSafeStartCapabilityIsDeclared(t *testing.T) {
	raw, err := os.ReadFile("crossplane.yaml")
	if err != nil {
		t.Fatalf("read package metadata: %v", err)
	}

	pkg := &providerPackage{}
	if err := yaml.Unmarshal(raw, pkg); err != nil {
		t.Fatalf("parse package metadata: %v", err)
	}

	if !hasString(pkg.Spec.Capabilities, "safe-start") {
		t.Fatal("package must declare safe-start so Crossplane grants CRD read permissions")
	}
}

func TestPackageBuilderPreservesCapabilities(t *testing.T) {
	raw, err := os.ReadFile("../Makefile")
	if err != nil {
		t.Fatalf("read Makefile: %v", err)
	}

	matches := regexp.MustCompile(`(?m)^UP_VERSION\s*=\s*v([0-9]+)\.([0-9]+)\.([0-9]+)$`).FindStringSubmatch(string(raw))
	if matches == nil {
		t.Fatal("Makefile must pin UP_VERSION to a validated release")
	}

	major := mustAtoi(t, matches[1])
	minor := mustAtoi(t, matches[2])
	patch := mustAtoi(t, matches[3])
	if major == 0 && (minor < 44 || minor == 44 && patch < 3) {
		t.Fatalf("UP_VERSION v%s.%s.%s drops provider capabilities from xpkg metadata; use v0.44.3 or newer", matches[1], matches[2], matches[3])
	}
}

func hasString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func mustAtoi(t *testing.T, value string) int {
	t.Helper()

	out, err := strconv.Atoi(value)
	if err != nil {
		t.Fatalf("parse version component %q: %v", value, err)
	}
	return out
}
