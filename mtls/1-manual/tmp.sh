
#!/bin/bash -e

# START OMIT
# START CLUSTER OMIT
# Create the cluster
k3d cluster create mtls-linkerd

# Build and inject our test image
docker build -t carsonoid/go-test-app ../..
docker save carsonoid/go-test-app -o app.tar
k3d image import -c mtls-linkerd app.tar
# END CLUSTER OMIT

# START LINKERD OMIT
# Install the linkerd CLI
curl --proto '=https' --tlsv1.2 -sSfL https://run.linkerd.io/install | sh

# Install linkerd
linkerd check --pre
linkerd install --crds | kubectl apply -f -
linkerd install | kubectl apply -f -

# Ensure linkerd is ready
kubectl -n linkerd wait --for=condition=Available=True \
  deploy/linkerd-identity \
  deploy/linkerd-destination \
  deploy/linkerd-proxy-injector

# END LINKERD OMIT


# Create client and server with linkerd sidecars injected
linkerd inject server-k8s.yaml | kubectl apply -f -
linkerd inject client-k8s.yaml | kubectl apply -f -

kubectl wait --for=condition=Available=True \
  deployment/test-server deployment/test-client

kubetail --follow --skip-colors


# END OMIT
