package main

import "github.com/carsonoid/talk-all-the-mtls-in-k8s/internal/demo"

const basePath = `./mtls/1-manual/`
const script = `
#!/bin/bash -e

// START OMIT
// START CLUSTER OMIT
k3d cluster create 1-manual

docker build -t carsonoid/go-test-app ../..
docker save carsonoid/go-test-app -o app.tar
k3d image import -c 1-manual app.tar
// END CLUSTER OMIT

// START CERTS OMIT
# Inject server tls secret, delete if it exists to make sure it's up to date
kubectl delete secret server-tls || true
kubectl create secret generic server-tls \
  --from-file=certs/server/tls.pem \
  --from-file=certs/server/tls-key.pem \
  --from-file=certs/server/ca.pem

# Inject client tls secret, delete if it exists to make sure it's up to date
kubectl delete secret client-tls || true
kubectl create secret generic client-tls \
  --from-file=certs/client/tls.pem \
  --from-file=certs/client/tls-key.pem \
  --from-file=certs/client/ca.pem
// END CERTS OMIT

// START RUN OMIT
# Create client and server
kubectl apply -f server-k8s.yaml
kubectl apply -f client-k8s.yaml

kubectl wait --for=condition=Available=True \
  deployment/test-server deployment/test-client

kubetail --follow --skip-colors
// END RUN OMIT
// END OMIT
`

func main() {
	demo.RunShellScript(basePath, script)
}
