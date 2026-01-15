# Provider Kion

`provider-kion` is a [Crossplane](https://crossplane.io/) provider that
is built using [Upjet](https://github.com/crossplane/upjet) code
generation tools and exposes XRM-conformant managed resources for the
Kion API.

## Getting Started

Install the provider by using the following command after changing the image tag
to the [latest release](https://marketplace.upbound.io/providers/enel1221/provider-kion):
```
up ctp provider install ghcr.io/enel1221/provider-kion:0.1.4
```

Alternatively, you can use declarative installation:
```
cat <<EOF | kubectl apply -f -
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-kion
spec:
  package: ghcr.io/enel1221/provider-kion:0.1.4
EOF
```

Notice that in this example Provider resource is referencing ControllerConfig with debug enabled.

You can see the API reference [here](https://doc.crds.dev/github.com/enel1221/provider-kion).

## Developing

Run code-generation pipeline:
```console
go run cmd/generator/main.go "$PWD"
```

Run against a Kubernetes cluster:

```console
make run
```

Build, push, and install:

```console
make all
```

Build binary:

```console
make build
```

## Creating a Release

To create a new release, follow these steps:

1. Ensure you're on the main branch with all changes committed and pushed
2. Create and push a version tag:
```bash
git checkout main
git tag -a v0.1.X -m "Release v0.1.X"
git push origin v0.1.X
```
3. Create/update the release branch from the tag:
```bash
git branch -D release-0.1  # Delete old release branch if exists
git checkout -b release-0.1 v0.1.X
git push -u origin release-0.1 -f
```

The CI workflow will automatically:
- Build the provider for linux/amd64 and linux/arm64
- Publish to `ghcr.io/enel1221/provider-kion:v0.1.X`
- Create a GitHub Release

## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please
open an [issue](https://github.com/enel1221/provider-kion/issues).
