#!/bin/bash -e

# Create the cluster
k3d cluster create mtls-linkerd

# Build and inject our test image
docker build -t carsonoid/go-test-app ../..
docker save carsonoid/go-test-app -o app.tar
k3d image import -c mtls-linkerd app.tar

# Install the linkerd CLI
curl --proto '=https' --tlsv1.2 -sSfL https://run.linkerd.io/install | sh

# Install linkerd
linkerd check --pre
linkerd install --crds | kubectl apply -f -
linkerd install | kubectl apply -f -

# Setup client and server
linkerd inject server-k8s.yaml | kubectl apply -f -
linkerd inject client-k8s.yaml | kubectl apply -f -

kubectl wait --for=condition=Available=True deployment/test-server deployment/test-client

kubetail --follow --skip-colors
