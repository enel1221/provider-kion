# AGENTS.md

## Purpose
This repository is a Crossplane provider for Kion built with Upjet from Kion's
official Terraform provider. Agents working here should preserve the generated
pipeline, keep manual config changes minimal and intentional, and avoid casual
edits to generated APIs or CRDs.

## Read First
Before editing code, read:
- `README.md`
- `Makefile`
- the relevant file under `config/`
- the nearest generated or controller code for the resource being changed

## Repo Priorities
- Treat `config/` as the primary handwritten source of truth for provider
  resource behavior.
- Expect most API, controller, CRD, example, and package changes to be
  generation outputs rather than hand-edited files.
- Preserve both cluster-scoped and namespaced provider support unless the task
  explicitly changes that contract.
- Keep the Kion API key rotation controller working alongside the generated
  provider behavior.

## Key Areas
- `config/`: handwritten Upjet resource mapping, external-name behavior, and
  provider setup
- `cmd/generator/`: generation entrypoint
- `cmd/provider/`: provider runtime entrypoint
- `internal/clients/`: Kion client behavior
- `internal/controller/keyrotation/`: Kion API key rotation controller
- `internal/controller/integration/`: envtest-backed integration coverage
- `apis/`, `package/`, `examples-generated/`: generated outputs

## Safety Rules
- Safe by default: repo inspection, focused edits in handwritten code,
  generation to local files, unit/integration tests, and local builds.
- Ask first: cluster deployment targets such as `local-deploy`, `e2e`, `uptest`,
  `controlplane.up`, image pushes, and any credentialed registry or Upbound
  publication flow.
- Avoid hand-editing generated `zz_*.go`, generated CRDs, or generated examples
  unless the task is explicitly about a post-generation fix and the reason is
  documented.

## Standard Workflow
1. Change handwritten config or controller code first.
2. Run the generation pipeline if the task affects provider schema or resource
   mappings.
3. Validate with focused Go tests before broader make targets.
4. Summarize which outputs were regenerated and what deployment/runtime paths
   were not exercised.

## Release Guardrails
- Before releasing a provider package, run the package metadata tests. SafeStart
  requires `package/crossplane.yaml` to declare the `safe-start` capability;
  Crossplane's RBAC manager turns that capability into `get`, `list`, and
  `watch` on `customresourcedefinitions.apiextensions.k8s.io`.
- Keep the package builder pinned to a version that preserves package
  capabilities in the xpkg metadata. `up v0.28.0` drops capabilities and will
  recreate the SafeStart RBAC failure even when `crossplane.yaml` is correct.
- Run the local CI-equivalent checks before pushing release fixes:
  `make vendor vendor.check`, `make lint`, `make test`,
  `make integration-test`, `make check-diff`, `make build`, and
  `make local-deploy`.
- After pushing, verify the GitHub Actions run on the pushed commit before
  tagging. Release tags use the `vX.Y.Z` format and should be followed by an
  explicit package pull check from `ghcr.io/enel1221/provider-kion:vX.Y.Z`.

## Source Of Truth
Detailed repo guidance lives under `.codex/`. More specific guidance there
overrides this file for the area being changed.
