# provider-kion: Crossplane 2.0 / Upjet v2 Upgrade Tasks

> **Branch**: `feat-v2-upgrade`
> **Current state**: Dependencies migrated (crossplane-runtime v2, upjet v2.2.0, Go 1.24.11). Dual API structure (cluster/namespaced) generated for 28 resources. Binary builds, lint passes. 5 more commits needed to be production-ready.

## Quick Terraform Provider Update Checklist

1. Update `TERRAFORM_PROVIDER_VERSION` and `TERRAFORM_NATIVE_PROVIDER_BINARY` in `Makefile`.
2. Run `make generate.init`.
3. If `make generate.init` reuses the old provider lock, remove `.work/terraform` and rerun.
4. Run `make generate`.
5. Run `go mod tidy`.
6. Run `go test ./config/... ./internal/...`.
7. Run `make lint`.
8. Run `go build ./...`.
9. Run `make build`.

The test suite now loads Terraform version and source values from environment or the top-level `Makefile`, so routine Terraform provider bumps should not require test-file edits.

---

## Phase 1: Upgrade Terraform Provider (0.3.31 → 0.3.33)

### - [x] 1.1 Update Terraform provider version in Makefile

**File**: `Makefile`
**What**: Bump the Kion Terraform provider from v0.3.31 to v0.3.33 (latest). Releases 0.3.32 and 0.3.33 are bug-fix only (label pagination fix + AWS account retry fix) — no new resources added.

**Changes needed** (lines 15–18 of `Makefile`):
```makefile
# BEFORE:
export TERRAFORM_PROVIDER_VERSION ?= 0.3.31
export TERRAFORM_NATIVE_PROVIDER_BINARY ?= terraform-provider-kion_v0.3.31

# AFTER:
export TERRAFORM_PROVIDER_VERSION ?= 0.3.33
export TERRAFORM_NATIVE_PROVIDER_BINARY ?= terraform-provider-kion_v0.3.33
```

### - [x] 1.2 Regenerate Terraform provider schema

**Prereq**: Task 1.1 complete
**What**: Pull updated docs and regenerate the provider schema JSON with the new version.

```bash
make generate.init
```

This runs:
1. `pull-docs` — clones `terraform-provider-kion` at v0.3.33 and copies `docs/resources` into the workspace
2. Schema generation — invokes `terraform providers schema -json` with the v0.3.33 binary to produce `config/schema.json`

**Files updated**: `config/schema.json`, `config/provider-metadata.yaml`

**Verification**: Diff the new schema.json against the old. Since 0.3.32/0.3.33 are bug-fix only, the resource list should remain the same 33 resources. Check with:
```bash
python3 -c "
import json
with open('config/schema.json') as f:
    data = json.load(f)
schema = data['provider_schemas']['registry.terraform.io/kionsoftware/kion']
for r in sorted(schema['resource_schemas'].keys()):
    print(r)
"
```

---

## Phase 2: Add Missing Resources to External Name Config

### - [x] 2.1 Add 5 missing resources to external_name.go

**File**: `config/external_name.go`
**What**: The current schema.json (v0.3.31) has 33 resources, but `ExternalNameConfigs` only maps 28. Five resources exist in the TF provider schema but are excluded from generation:

| Resource | Added in TF ver | Description |
|---|---|---|
| `kion_app_config` | v0.3.28 | Global Kion application configuration |
| `kion_custom_account` | v0.3.29 | Custom (non-cloud-provider) accounts |
| `kion_custom_variable` | pre-v0.3.27 | Custom variables |
| `kion_custom_variable_override` | pre-v0.3.27 | Custom variable overrides |
| `kion_project_note` | pre-v0.3.27 | Project notes |

**Changes**: Add these 5 entries to the `ExternalNameConfigs` map in `config/external_name.go`:
```go
var ExternalNameConfigs = map[string]config.ExternalName{
    // ... existing 28 entries ...

    // New resources (previously excluded):
    "kion_app_config":               config.IdentifierFromProvider,
    "kion_custom_account":           config.IdentifierFromProvider,
    "kion_custom_variable":          config.IdentifierFromProvider,
    "kion_custom_variable_override": config.IdentifierFromProvider,
    "kion_project_note":             config.IdentifierFromProvider,
}
```

### - [x] 2.2 Add resource configurators for 5 new resources

**File**: `config/kionprovider/config.go`
**What**: Each resource needs a configurator that sets at minimum `r.ShortGroup = "kion"`. Some need explicit `Kind` overrides and cross-resource references.

**Add to the `Configure` function**:
```go
// ── App Config ───────────────────────────────────────────────────

p.AddResourceConfigurator("kion_app_config", func(r *config.Resource) {
    r.ShortGroup = shortGroup
    r.Kind = "AppConfig"
})

// ── Custom Account ───────────────────────────────────────────────

p.AddResourceConfigurator("kion_custom_account", func(r *config.Resource) {
    r.ShortGroup = shortGroup
    r.Kind = "CustomAccount"
    r.References["project_id"] = config.Reference{
        TerraformName: "kion_project",
    }
})

// ── Custom Variables ─────────────────────────────────────────────

p.AddResourceConfigurator("kion_custom_variable", func(r *config.Resource) {
    r.ShortGroup = shortGroup
    r.Kind = "CustomVariable"
})

p.AddResourceConfigurator("kion_custom_variable_override", func(r *config.Resource) {
    r.ShortGroup = shortGroup
    r.Kind = "CustomVariableOverride"
    r.References["custom_variable_id"] = config.Reference{
        TerraformName: "kion_custom_variable",
    }
})

// ── Project Note ─────────────────────────────────────────────────

p.AddResourceConfigurator("kion_project_note", func(r *config.Resource) {
    r.ShortGroup = shortGroup
    r.Kind = "ProjectNote"
    r.References["project_id"] = config.Reference{
        TerraformName: "kion_project",
    }
})
```

**Reference verification**: Check `config/schema.json` for each resource to confirm which fields reference other resources (look for `project_id`, `ou_id`, `custom_variable_id`, etc.). The schema.json is a single-line JSON file — search for the resource name to find its attributes.

### - [x] 2.3 Update external_name_test.go for new resource count

**File**: `config/external_name_test.go`
**What**: Update the expected resource count from 28 to 33, and add the 5 new resources to the `AllResources` test list.

**Changes**:
```go
// Line 14: update count
const wantCount = 33  // was 28

// In the AllResources test, add to the expected list:
"kion_app_config",
"kion_custom_account",
"kion_custom_variable",
"kion_custom_variable_override",
"kion_project_note",
```

### - [x] 2.4 Update provider_test.go for new resources

**File**: `config/provider_test.go`
**What**: Add the 5 new resources to the test expectations:

1. **`expectedKinds`** map — add entries for the new explicit Kind overrides:
```go
"kion_app_config":               "AppConfig",
"kion_custom_account":           "CustomAccount",
"kion_custom_variable":          "CustomVariable",
"kion_custom_variable_override": "CustomVariableOverride",
"kion_project_note":             "ProjectNote",
```

2. **`referenceTests`** map — add reference expectations:
```go
"kion_custom_account": {
    "project_id": "kion_project",
},
"kion_custom_variable_override": {
    "custom_variable_id": "kion_custom_variable",
},
"kion_project_note": {
    "project_id": "kion_project",
},
```

---

## Phase 3: Full Code Regeneration

### - [x] 3.1 Run full code generation pipeline

**Prereq**: Tasks 1.1, 1.2, 2.1, 2.2 complete
**What**: Regenerate all Upjet-generated code (types, controllers, CRDs, deepcopy, resolvers, examples).

```bash
make generate
```

This executes `apis/generate.go` which:
1. Removes old `zz_*` files, CRDs, and controller scaffolds
2. Runs `cmd/generator/main.go` (Upjet pipeline) for both cluster and namespaced providers
3. Generates CRD manifests via `controller-gen`
4. Generates Crossplane methodsets via `angryjet`

**Expected output**: 33 resources × 2 scopes = new files in:
- `apis/cluster/kion/v1alpha1/` — ~66+ files (types + terraformed per resource)
- `apis/namespaced/kion/v1alpha1/` — ~66+ files (mirror)
- `internal/controller/cluster/kion/` — 33 controller directories
- `internal/controller/namespaced/kion/` — 33 controller directories
- `package/crds/` — CRDs for all resources in both scopes
- `examples-generated/` — example manifests for both scopes

### - [x] 3.2 Verify generated resource consistency

**Prereq**: Task 3.1 complete
**What**: Confirm all layers have the same 33 resources.

```bash
# Count generated type files per scope
ls apis/cluster/kion/v1alpha1/zz_*_types.go | wc -l        # expect ~33
ls apis/namespaced/kion/v1alpha1/zz_*_types.go | wc -l     # expect ~33

# Count controller directories per scope
ls -d internal/controller/cluster/kion/*/ | wc -l           # expect 33
ls -d internal/controller/namespaced/kion/*/ | wc -l        # expect 33

# Count CRDs (should cover both scopes)
ls package/crds/ | wc -l

# Verify the 5 new resources exist
for r in appconfig customaccount customvariable customvariableoverride projectnote; do
    ls apis/cluster/kion/v1alpha1/zz_${r}_types.go 2>/dev/null && echo "OK: $r" || echo "MISSING: $r"
done
```

---

## Phase 4: Update Examples & Test Fixtures for Crossplane 2.0

### - [x] 4.1 Add namespaced ProviderConfig example

**File**: `examples/providerconfig/providerconfig.yaml`
**What**: Currently uses the old cluster-scoped API `kion.upbound.io/v1beta1`. Add a namespaced variant or update to demonstrate both patterns.

**Current content**:
```yaml
apiVersion: kion.upbound.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      name: example-creds
      namespace: crossplane-system
      key: credentials
```

**New content** — create a separate file `examples/providerconfig/providerconfig-namespaced.yaml`:
```yaml
apiVersion: kion.m.upbound.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
  namespace: crossplane-system
spec:
  credentials:
    source: Secret
    secretRef:
      name: example-creds
      key: credentials
```

**Note**: The cluster-scoped file should remain as-is for backward compatibility. The key difference is:
- Cluster: `kion.upbound.io/v1beta1` — no namespace on metadata, secretRef has explicit namespace
- Namespaced: `kion.m.upbound.io/v1beta1` — has namespace on metadata, secretRef resolves in same namespace

### - [x] 4.2 Update e2e test setup for namespaced ProviderConfig

**File**: `cluster/test/setup.sh`
**What**: Currently creates only the old cluster-scoped ProviderConfig. Should also create the namespaced variant for e2e testing.

**Current ProviderConfig block (line 17–26)**:
```yaml
apiVersion: kion.upbound.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      name: provider-secret
      namespace: upbound-system
      key: credentials
```

**Add a second block after the existing one**:
```bash
echo "Creating a namespaced provider config..."
cat <<EOF | ${KUBECTL} apply -f -
apiVersion: kion.m.upbound.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
  namespace: upbound-system
spec:
  credentials:
    source: Secret
    secretRef:
      name: provider-secret
      key: credentials
EOF
```

### - [x] 4.3 Update install.yaml package reference

**File**: `examples/install.yaml`
**What**: Currently references `ghcr.io/enel1221/provider-kion:0.1.4`. Should be updated for the v2 release.

**Current**:
```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-kion
spec:
  package: ghcr.io/enel1221/provider-kion:0.1.4
```

**Updated** (version TBD — placeholder for now):
```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-kion
spec:
  package: ghcr.io/enel1221/provider-kion:v0.2.0
```

### - [x] 4.4 Update package/crossplane.yaml with Crossplane version constraint

**File**: `package/crossplane.yaml`
**What**: Currently has no version constraint. Should declare Crossplane 2.0 requirement since this provider uses v2 APIs.

**Current**:
```yaml
apiVersion: meta.pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-kion
```

**Updated**:
```yaml
apiVersion: meta.pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-kion
spec:
  crossplane:
    version: ">=v2.0.0"
```

**Reference**: See `provider-upjet-azure` at `/Users/inelson/github/provider-upjet-azure/package/crossplane.yaml` for the pattern used by production Upjet providers.

---

## Phase 5: Build & Test Verification

### - [x] 5.1 Run go mod tidy

**What**: Clean up Go module dependencies. The `go.mod` was flagged for `github.com/google/go-cmp v0.7.0` and `k8s.io/api v0.33.0` being indirect but possibly needed as direct.

```bash
go mod tidy
```

**Verify** no unexpected dependency changes by reviewing the diff.

### - [x] 5.2 Verify clean compilation

**Prereq**: Tasks 3.1, 5.1 complete

```bash
go build ./...
```

Must exit 0 with no errors. Common issues after generation:
- Import path mismatches (should all use `github.com/enel1221/provider-kion`)
- Missing cross-resource reference targets (if a referenced TerraformName doesn't exist)
- Scheme registration gaps

### - [x] 5.3 Run all unit tests

**Prereq**: Task 5.2 passes

```bash
go test ./config/... ./internal/...
```

**Test files that must pass**:
- `config/external_name_test.go` — validates 33 resources (updated count), all external name patterns
- `config/provider_test.go` — validates Kind overrides (23 now), ShortGroup, references, SchemaElementOptions
- `internal/clients/kion_test.go` — validates TerraformSetupBuilder and ProviderConfig resolution

### - [x] 5.4 Run linter

```bash
make lint
```

Uses golangci-lint v1.64.5. Must pass with no errors. The lint binary is cached at `_output/lint/`.

### - [x] 5.5 Verify clean build (no dirty version)

```bash
make build
cat _output/version
```

The version should NOT contain `-dirty`. If it does, commit all changes first, then rebuild.

---

## Phase 6: Production Readiness

### - [x] 6.1 Verify ClusterProviderConfig CRD exists

**What**: Crossplane 2.0 namespaced providers need a `ClusterProviderConfig` CRD for cluster-wide credential configuration.

```bash
ls package/crds/ | grep -i clusterproviderconfig
# Should find: kion.m.upbound.io_clusterproviderconfigs.yaml
```

Verify the CRD:
```bash
head -20 package/crds/kion.m.upbound.io_clusterproviderconfigs.yaml
```

Should show `scope: Cluster` and `group: kion.m.upbound.io`.

### - [x] 6.2 Verify feature flags in features.go

**File**: `internal/features/features.go`
**What**: Confirm Crossplane 2.0 feature flags are properly defined.

**Current content** (confirmed correct):
```go
const (
    EnableBetaManagementPolicies xpfeature.Flag = xpfeature.EnableBetaManagementPolicies
    EnableAlphaWatches xpfeature.Flag = "EnableAlphaWatches"
)
```

**Verify in cmd/provider/main.go**: Management Policies should be enabled by default. Search for `EnableBetaManagementPolicies` in `cmd/provider/main.go` and confirm it defaults to `true`.

### - [x] 6.3 Review CI workflow compatibility

**File**: `.github/workflows/ci.yml`
**What**: Confirm CI pipeline works with the v2 structure.

**Check**:
- Go version matches `1.24.11` (see `GO_REQUIRED_VERSION` in Makefile)
- Build targets still work (`go build ./...`, `make lint`, etc.)
- Test commands cover new paths (`go test ./config/... ./internal/...`)
- Generation steps if any (some CIs run `make generate` and diff)

### - [ ] 6.4 Run e2e tests (if Kion environment available)

**Prereq**: All previous tasks complete
**What**: If a Kion test environment is accessible, run the end-to-end tests.

```bash
export UPTEST_CLOUD_CREDENTIALS='{"url":"https://kion.example.com","apikey":"...","skipsslvalidation":true}'
make e2e
```

This uses `uptest` to:
1. Deploy the provider to a local Kind cluster
2. Run `cluster/test/setup.sh` to create credentials and ProviderConfig
3. Apply example resources and validate they reconcile

---

## Summary

| Phase | Tasks | Status |
|---|---|---|
| 1. TF Provider Upgrade | 1.1, 1.2 | Complete |
| 2. Add Missing Resources | 2.1, 2.2, 2.3, 2.4 | Complete |
| 3. Code Regeneration | 3.1, 3.2 | Complete |
| 4. Examples & Fixtures | 4.1, 4.2, 4.3, 4.4 | Complete |
| 5. Build & Test | 5.1, 5.2, 5.3, 5.4, 5.5 | Complete |
| 6. Production Readiness | 6.1, 6.2, 6.3, 6.4 | 6.1-6.3 complete, 6.4 pending live env |

### Key context for subagents

- **Project module path**: `github.com/enel1221/provider-kion`
- **Upjet v2 imports**: `github.com/crossplane/upjet/v2/pkg/config`
- **Crossplane runtime v2**: `github.com/crossplane/crossplane-runtime/v2`
- **Root groups**: cluster=`upbound.io`, namespaced=`m.upbound.io`
- **All resources use ShortGroup**: `"kion"`
- **ExternalName strategy**: All use `config.IdentifierFromProvider`
- **Reference pattern for provider-upjet-azure** (production v2 Upjet provider): `/Users/inelson/github/provider-upjet-azure`
- **`kion_user_group`** generates as Kind `Group` (auto-derived by Upjet stripping `kion_` prefix and using last segment). It IS working — the generated files are `zz_group_types.go` and controller is in `group/`. No action needed.
- **TF provider schema** is a single-line JSON at `config/schema.json`. Search for resource names to inspect their attributes.
- **Terraform version locked** to 1.5.7 (BSL license compliance — cannot use >= 1.6)
