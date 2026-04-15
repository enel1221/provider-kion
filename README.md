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

## API Key Rotation

Kion API keys expire (typically every 14 days). This provider includes a
built-in controller that automatically rotates the key before it expires by
calling the Kion `POST /api/v3/app-api-key/rotate` endpoint and updating the
Secret in-place.

> **Reference:** The full Kion API documentation is available at
> [https://kion.ccmo.socom.mil/swagger/](https://kion.ccmo.socom.mil/swagger/)

### Enabling Rotation

Add annotations to the credentials Secret referenced by your ProviderConfig:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: kion-creds
  namespace: crossplane-system
  annotations:
    kion.upbound.io/rotate: "true"
    kion.upbound.io/rotation-interval: "168h"  # optional, default 7 days
type: Opaque
stringData:
  credentials: |
    {
      "apikey": "YOUR_KION_API_KEY",
      "url": "https://kion.ccmo.socom.mil",
      "skipsslvalidation": "false"
    }
```

### Annotations

| Annotation | Set By | Description |
|---|---|---|
| `kion.upbound.io/rotate` | User | Set to `"true"` to enable automatic rotation |
| `kion.upbound.io/rotation-interval` | User | Go duration string (default `168h` = 7 days) |
| `kion.upbound.io/credentials-key` | User | Override the Secret data key (default `credentials`) |
| `kion.upbound.io/last-rotation` | Controller | RFC3339 timestamp of the last successful rotation |
| `kion.upbound.io/rotation-status` | Controller | `success` or `error: <message>` |

### How It Works

1. The controller watches Secrets with `kion.upbound.io/rotate: "true"`.
2. When the rotation interval elapses, it calls the Kion rotate API with the
   current key.
3. The new key is written back into the Secret, and `last-rotation` /
   `rotation-status` annotations are updated.
4. All existing ProviderConfigs referencing that Secret pick up the new key
   automatically on their next reconcile.

A complete example is available at
[`examples/providerconfig/secret-with-rotation.yaml`](examples/providerconfig/secret-with-rotation.yaml).

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
