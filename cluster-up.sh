#!/bin/bash

docker build -t carsonoid/go-test-app .
docker save carsonoid/go-test-app -o app.tar
k3d image import app.tar


kubectl delete secret generic server-tls
kubectl create secret generic server-tls \
  --from-file=mtls/0-manually/certs/server/tls.pem \
  --from-file=mtls/0-manually/certs/server/tls-key.pem \
  --from-file=mtls/0-manually/certs/server/ca.pem
kubectl delete po --all

