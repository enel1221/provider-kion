# provider-kion agent notes

## Output discipline

- Avoid streaming high-volume command output into VS Code. For generation, lint, build, and other noisy targets, redirect output to files under `.work/` and inspect the log files instead.
- Prefer commands like:
  - `make generate.init > .work/generate-init.log 2>&1`
  - `make generate > .work/generate.log 2>&1`
  - `make lint > .work/lint.log 2>&1`
  - `make build > .work/build.log 2>&1`

## Terraform generation trap

- If `make generate.init` fails after changing `TERRAFORM_PROVIDER_VERSION`, check for a stale lock in `.work/terraform/.terraform.lock.hcl`.
- The usual fix in this repo is to remove `.work/terraform` and rerun `make generate.init`.

## Expected workflow

- When adding or upgrading Terraform resources, update the hand-maintained config first:
  - `config/external_name.go`
  - `config/kionprovider/config.go`
  - `config/external_name_test.go`
  - `config/provider_test.go`
- Then run validation in this order:
  1. `make generate.init`
  2. `make generate`
  3. `go mod tidy`
  4. `go test ./config/... ./internal/...`
  5. `make lint`
  6. `go build ./...`
  7. `make build`

## Repo-specific expectations

- This provider maintains both cluster-scoped and namespaced APIs. Resource additions should result in generated changes under both `apis/cluster` and `apis/namespaced`, both controller trees, and both CRD groups.
- Keep the cluster-scoped ProviderConfig example for compatibility, and maintain a namespaced ProviderConfig example and e2e setup path alongside it.