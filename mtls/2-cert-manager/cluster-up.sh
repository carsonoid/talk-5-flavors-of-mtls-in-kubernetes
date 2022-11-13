#!/bin/bash -e

k3d cluster create 2-cert-manager

docker build -t carsonoid/go-test-app ../..
docker save carsonoid/go-test-app -o app.tar
k3d image  import -c 2-cert-manager app.tar

# Install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.10.0/cert-manager.yaml

# Create issuer cert
bash ./openssl-certgen.sh

# Add ca secret
kubectl -n cert-manager delete secret ca-tls || true
kubectl -n cert-manager create secret generic ca-tls \
  --from-file=tls.crt=certs/ca/tls.pem \
  --from-file=tls.key=certs/ca/tls-key.pem

# Add issuer using ca secret
kubectl apply -f <(cat <<EOL
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: ca-issuer
  namespace: cert-manager
spec:
  ca:
    secretName: ca-tls
EOL
)

# Create client and server
kubectl apply -f server-k8s.yaml
kubectl apply -f client-k8s.yaml

kubectl wait --for=condition=Available=True deployment/test-server deployment/test-client

kubetail --follow --skip-colors
