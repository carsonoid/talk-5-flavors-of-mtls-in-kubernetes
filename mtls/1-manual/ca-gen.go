package main

import "github.com/carsonoid/talk-all-the-mtls-in-k8s/internal/demo"

const basePath = `.`
const script = `
set -e

// START OMIT
mkdir -p certs/ca
cat > certs/ca/csr.json <<EOL
{
    "hosts": [
        "127.0.0.1",
        "localhost"
    ],
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C":  "US",
            "L":  "Utah",
            "O":  "UT Kubernetes",
            "OU": "Infra",
            "ST": "Utah"
        }
    ]
}
EOL

cfssl genkey -initca certs/ca/csr.json | cfssljson -bare certs/ca/tls
// END OMIT
`

func main() {
	demo.RunShellScript(basePath, script)
}
