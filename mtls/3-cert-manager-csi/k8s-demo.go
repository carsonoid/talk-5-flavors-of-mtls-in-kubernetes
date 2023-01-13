package main

import "github.com/carsonoid/talk-all-the-mtls-in-k8s/internal/demo"

const basePath = `./mtls/1-manual/`
const script = `
#!/bin/bash -e

// START CLUSTER OMIT
# Create the cluster with a shared mount where needed for the csi driver
k3d cluster create mtls-cert-manager-csi -v ~/k3d/run:/tmp/cert-manager-csi-driver:shared

# Build and inject our test image
docker build -t carsonoid/go-test-app ../..
docker save carsonoid/go-test-app -o app.tar
k3d image import -c mtls-cert-manager-csi app.tar
# END CLUSTER OMIT

# START CM OMIT
# Setup cert-manager + csi driver
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.10.0/cert-manager.yaml

# Use helm to install the csi driver
# - don't be fooled by the "upgrade" command
#   the '-i' flag means install if missing
helm upgrade -i -n cert-manager \
  cert-manager-csi-driver \
  jetstack/cert-manager-csi-driver --wait
# END CM OMIT

# START VAULT INSTALL OMIT
helm repo add hashicorp https://helm.releases.hashicorp.com
kubectl create ns vault
helm upgrade -i -n vault vault hashicorp/vault \
  --set server.dev.enabled=true \
  --set server.dev.devRootToken=insecure-root-token
# END VAULT INSTALL OMIT

# START VAULT OMIT
kubectl  --namespace=vault port-forward service/vault 8200 &

export VAULT_ADDR=http://localhost:8200
export VAULT_TOKEN=insecure-root-token
# END VAULT OMIT

# START VAULT SETUP OMIT
vault secrets enable pki

# the default lease is 768 hours, we will set it to be ~6 months
vault secrets tune -max-lease-ttl=4380h pki

# Create a CA cert and write it to a file
vault write -field=certificate pki/root/generate/internal \
     common_name="kube.local" \
     issuer_name="root-2022" \
     ttl=87600h > root_2022_ca.crt

# Create a a signer role
vault write pki/roles/kube-tls-signer \
    allow_any_name=true \
    ttl=87600h
# END VAULT SETUP OMIT

# START CM_ISSUER OMIT
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
# END CM_ISSUER OMIT

# START RUN OMIT
# Create client and server
kubectl apply -f server-k8s.yaml
kubectl apply -f client-k8s.yaml

kubectl wait --for=condition=Available=True \
  deployment/test-server deployment/test-client

kubetail --follow --skip-colors
# END RUN OMIT
`

func main() {
	demo.RunShellScript(basePath, script)
}
