# Workflows

## Resource Addition Or Upgrade
1. Update handwritten config first:
   - `config/external_name.go`
   - `config/kionprovider/config.go`
   - provider setup/tests as needed
2. Run:
   - `make generate.init`
   - `make generate`
3. Normalize module state with `go mod tidy` if needed.
4. Validate with:
   - `go test ./config/... ./internal/...`
   - `make lint`
   - `go build ./...`
   - `make build`

## Key Rotation Changes
1. Read the provider config examples and rotation annotations in `README.md`.
2. Inspect:
   - `internal/controller/keyrotation/`
   - `internal/clients/kion.go`
   - provider config handling in `config/`
3. Prefer focused controller/client tests before broader make targets.

## Review Expectations
- After generation, inspect whether both cluster and namespaced outputs changed
  coherently.
- Check generated CRDs and examples for unexpected contract drift.
- Call out any network-dependent or cluster-dependent validation that was not
  run locally.

## Release Cycle
1. Start from an up-to-date `main` and keep the release fix scoped.
2. For SafeStart or provider runtime changes, verify `package/crossplane.yaml`
   includes:
   - `spec.capabilities: [safe-start]`
   - no reliance on `spec.controller.permissionRequests`; Crossplane 2.0 grants
     CRD read permissions from the `safe-start` capability.
   - `Makefile` pins `UP_VERSION` to a package builder that preserves
     capabilities in the xpkg metadata.
   - `Makefile` keeps `RELEASE_BRANCH_FILTER` matching `v%` so tag workflows
     publish versioned GHCR package tags.
3. Run local validation before pushing:
   - `make vendor vendor.check`
   - `make lint`
   - `make test`
   - `make integration-test`
   - `make check-diff`
   - `make build`
   - `make local-deploy`
4. Confirm `make local-deploy` verifies the generated provider system
   ClusterRole includes `get`, `list`, and `watch` on
   `customresourcedefinitions.apiextensions.k8s.io`.
5. Push the commit and verify the GitHub Actions run for that exact SHA.
6. Create the release tag, for example `v1.0.4`, and wait for the tag workflow
   to publish artifacts and the GitHub release.
7. Verify the released package is pullable from GHCR, for example:
   - `oras manifest fetch ghcr.io/enel1221/provider-kion:v1.0.4`
   - `oras pull ghcr.io/enel1221/provider-kion:v1.0.4 -o /tmp/provider-kion-v1.0.4`
