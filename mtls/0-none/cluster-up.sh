#!/bin/bash

k3d cluster create 0-none

docker build -t carsonoid/go-test-app ../..
docker save carsonoid/go-test-app -o app.tar
k3d image import -c 0-none app.tar

kubectl apply -f server-k8s.yaml
kubectl apply -f client-k8s.yaml

kubectl wait --for=condition=Available=True deployment/test-server deployment/test-client

kubetail --follow --skip-colors
