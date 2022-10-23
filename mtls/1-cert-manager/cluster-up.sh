#!/bin/bash -e

k3d cluster create mtls-cert-manager

docker build -t carsonoid/go-test-app ../../..
docker save carsonoid/go-test-app -o app.tar
k3d image  import -c mtls-cert-manager app.tar

kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.10.0/cert-manager.yaml

kubectl -n cert-manager delete secret ca-tls || true
kubectl -n cert-manager create secret generic ca-tls \
  --from-file=tls.crt=certs/ca/tls.pem \
  --from-file=tls.key=certs/ca/tls-key.pem

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
