#!/bin/bash

# Create a new cluster for flavor 0
k3d cluster create 0-none

# Build and import the image
docker build -t carsonoid/go-test-app ../..
docker save carsonoid/go-test-app -o app.tar
k3d image import -c 0-none app.tar

# Deploy the server and client
kubectl apply -f server-k8s.yaml
kubectl apply -f client-k8s.yaml

# Wait for the deployments to be ready
kubectl wait --for=condition=Available=True \
    deployment/test-server deployment/test-client

# Watch the logs
kubetail --follow --skip-colors
