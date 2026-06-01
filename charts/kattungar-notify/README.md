# Kattungar Notify Helm chart

This chart deploys the Kattungar Notify server with a persistent sqlite volume,
a ClusterIP Service, a Gateway API HTTPRoute, and secret-backed APNs and Google
Calendar credentials.

## Build and publish the image

Build and push the server image to GitHub Container Registry:

```sh
make push-server
```

This pushes `ghcr.io/ddeville/kattungar-notify:<chart version>`. The chart uses
that same image repository and chart-version tag by default. Set `image.digest`
to pin the image by digest.

## Create the runtime secret

The chart only supports using a pre-existing Secret. By default it looks for a
Secret named after the Helm release. For a release named `kattungar-notify`,
create it like this:

```sh
kubectl create namespace kattungar-notify
kubectl -n kattungar-notify create secret generic kattungar-notify \
  --from-file=server-api-keys.json=./server-api-keys.json \
  --from-file=apns-key.p8=./AuthKey_KEYID.p8 \
  --from-file=google-client-credentials.json=./client_credentials.json \
  --from-literal=apns-key-id=KEYID \
  --from-literal=google-refresh-token=REFRESH_TOKEN \
  --from-literal=google-calendar-id=CALENDAR_ID
```

`server-api-keys.json` must be a JSON array of admin API keys, matching what the
server expects today.

If you use a differently named Secret, set `secret.name`. The chart references
the Secret but does not create or validate it; create it directly or with an
`ExternalSecret` before the pod needs to start.

## Persistence

By default the chart uses an ephemeral `emptyDir` for sqlite data. This lets the
app run without cluster storage, but registered devices and notification history
are lost when the pod is replaced.

Enable a `1Gi` `ReadWriteOnce` PVC when you want data to survive restarts:

```sh
helm upgrade --install kattungar-notify ./charts/kattungar-notify \
  --namespace kattungar-notify \
  --create-namespace \
  --set persistence.enabled=true \
  --set persistence.storageClass=STORAGE_CLASS
```

When creating a PVC, leave `persistence.storageClass` empty to use the cluster
default StorageClass, or set it to `-` to render `storageClassName: ""`. Set
`persistence.existingClaim=CLAIM_NAME` to use a PVC managed outside this chart.

## Install

```sh
helm upgrade --install kattungar-notify ./charts/kattungar-notify \
  --namespace kattungar-notify \
  --create-namespace \
  --set gateway.name=GATEWAY_NAME \
  --set gateway.namespace=GATEWAY_NAMESPACE \
  --set gateway.listenerName=LISTENER_NAME \
  --set gateway.hostnames[0]=notify.example.com
```

Set `gateway.enabled=false` if you only want the in-cluster Service.
