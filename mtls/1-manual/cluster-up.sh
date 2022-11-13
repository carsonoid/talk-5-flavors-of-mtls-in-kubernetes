#!/bin/bash -e

k3d cluster create 1-manual

docker build -t carsonoid/go-test-app ../..
docker save carsonoid/go-test-app -o app.tar
k3d image import -c 1-manual app.tar

# Generate certs
bash ./openssl-certgen.sh

# Inject server tls secret
kubectl delete secret server-tls || true
kubectl create secret generic server-tls \
  --from-file=certs/server/tls.pem \
  --from-file=certs/server/tls-key.pem \
  --from-file=certs/server/ca.pem

# Inject client tls secret
kubectl delete secret client-tls || true
kubectl create secret generic client-tls \
  --from-file=certs/client/tls.pem \
  --from-file=certs/client/tls-key.pem \
  --from-file=certs/client/ca.pem

# Create client and server
kubectl apply -f server-k8s.yaml
kubectl apply -f client-k8s.yaml
