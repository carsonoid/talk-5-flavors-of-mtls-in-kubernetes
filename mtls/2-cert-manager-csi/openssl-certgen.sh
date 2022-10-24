#!/bin/bash -e

# Generate CA key and cert

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

