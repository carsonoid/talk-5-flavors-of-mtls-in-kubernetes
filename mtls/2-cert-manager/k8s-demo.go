package main

import "github.com/carsonoid/talk-all-the-mtls-in-k8s/internal/demo"

const basePath = `./mtls/1-manual/`
const script = `
#!/bin/bash -e

// START OMIT
// START CLUSTER OMIT
k3d cluster create 2-cert-manager

docker build -t carsonoid/go-test-app ../..
docker save carsonoid/go-test-app -o app.tar
k3d image  import -c 2-cert-manager app.tar
// END CLUSTER OMIT

// START CM OMIT
# Install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.10.0/cert-manager.yamlm \
--from-file=certs/client/ca.pem
// END CM OMIT

// START CA OMIT
# Create CA
mkdir -p certs/ca
cat > certs/ca/csr.json <<EOL
{
    "hosts": [
        "127.0.0.1",
        "localhost"
    ],
    "key": {"algo": "rsa","size": 2048},
    "names": [
        {"C":"US","L":"Utah","O":"UT Kubernetes","OU":"Infra","ST":"Utah"        }
    ]
}
EOL
cfssl genkey -initca certs/ca/csr.json | cfssljson -bare certs/ca/tls

# Add ca secret
kubectl -n cert-manager delete secret ca-tls || true
kubectl -n cert-manager create secret generic ca-tls \
--from-file=tls.crt=certs/ca/tls.pem \
--from-file=tls.key=certs/ca/tls-key.pem
// END CA OMIT

// START ISSUER OMIT
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
// END ISSUER OMIT

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
