#!/bin/bash -e

docker build -t carsonoid/go-test-app ../../..
docker save carsonoid/go-test-app -o app.tar
k3d image import app.tar


kubectl delete secret server-tls || true
kubectl create secret generic server-tls \
  --from-file=certs/server/tls.pem \
  --from-file=certs/server/tls-key.pem \
  --from-file=certs/server/ca.pem

kubectl delete secret client-tls || true
kubectl create secret generic client-tls \
  --from-file=certs/client/tls.pem \
  --from-file=certs/client/tls-key.pem \
  --from-file=certs/client/ca.pem


kubectl delete po --all

