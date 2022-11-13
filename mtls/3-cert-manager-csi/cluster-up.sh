#!/bin/bash -e

# Create the cluster with a shared mount where needed for the csi driver
k3d cluster create mtls-cert-manager-csi -v ~/k3d/run:/tmp/cert-manager-csi-driver:shared

# Build and inject our test image
docker build -t carsonoid/go-test-app ../..
docker save carsonoid/go-test-app -o app.tar
k3d image import -c mtls-cert-manager-csi app.tar

# Setup cert-manager + csi driver
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.10.0/cert-manager.yaml
helm upgrade -i -n cert-manager cert-manager-csi-driver jetstack/cert-manager-csi-driver --wait

# Setup vault as issuer (self-signed not supported for CSI)
helm repo add hashicorp https://helm.releases.hashicorp.com
kubectl create ns vault
helm install -n vault vault hashicorp/vault --set server.dev.enabled=true --set server.dev.devRootToken=insecure-root-token

# Do Vault setup
#
# NOTE: As this is a dev server instance, all vault settings and data will
#       be lost when the active vault pod is deleted.
#
#       This means the ClusterIssuer will no longer be valid after the pod is deleted
#       unless the setup is done again
kubectl  --namespace=vault port-forward service/vault 8200 &

export VAULT_ADDR=http://localhost:8200
export VAULT_TOKEN=insecure-root-token

vault secrets enable pki

vault secrets tune -max-lease-ttl=87600h pki

vault write -field=certificate pki/root/generate/internal \
     common_name="kube.local" \
     issuer_name="root-2022" \
     ttl=87600h > root_2022_ca.crt

vault write pki/roles/kube-tls-signer \
    allow_any_name=true \
    ttl=87600h

CA_BASE_64=$(base64 -w0 root_2022_ca.crt)
kubectl replace -f <(cat <<EOL
---
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: cert-manager-vault-token
  namespace: cert-manager
stringData:
  token: insecure-root-token
---
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: vault-issuer
spec:
  vault:
    path: pki/sign/kube-tls-signer
    server: http://vault.vault.svc.cluster.local:8200
    caBundle: $CA_BASE_64
    auth:
      tokenSecretRef:
          name: cert-manager-vault-token
          key: token
EOL
)

# Bring up the server and client
kubectl apply -f server-k8s.yaml
kubectl apply -f client-k8s.yaml

kubectl wait --for=condition=Available=True deployment/test-server deployment/test-client

kubetail --follow --skip-colors
