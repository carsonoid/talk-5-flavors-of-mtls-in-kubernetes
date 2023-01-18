package main

import "github.com/carsonoid/talk-all-the-mtls-in-k8s/internal/demo"

const basePath = `./mtls/1-manual/`
const script = `
set -e

// START OMIT
mkdir -p certs/server
cat > certs/server/csr.json <<EOL
{
    "CN": "server",
    "hosts": [
        "127.0.0.1",
        "localhost",
        "server",
        "server.default.svc.cluster.local"
    ],
    "usages": [
        "server auth",
        "client auth"
    ],
    "key": {"algo":"rsa","size":2048},
    "names": [
        {"C":"US","L":"Utah","O":"UT Kubernetes","OU":"Infra","ST":"Utah"}
    ]
}
EOL

cp certs/ca/tls.pem certs/server/ca.pem
cfssl gencert -ca certs/ca/tls.pem -ca-key certs/ca/tls-key.pem \
    certs/server/csr.json | cfssljson -bare certs/server/tls
// END OMIT
`

func main() {
	demo.RunShellScript(basePath, script)
}
