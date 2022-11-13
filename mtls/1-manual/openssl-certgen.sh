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

mkdir -p certs/server
cat > certs/server/csr.json <<EOL
{
    "hosts": [
        "127.0.0.1",
        "localhost",
        "server"
    ],
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "CN": "server",
    "names": [
        {
            "C":  "US",
            "L":  "Utah",
            "O":  "UT Kubernetes",
            "OU": "Infra",
            "ST": "Utah"
        }
    ],
    "usages": [
        "server auth",
        "client auth"
    ]
}
EOL
cp certs/ca/tls.pem certs/server/ca.pem
cfssl gencert -ca certs/ca/tls.pem -ca-key certs/ca/tls-key.pem certs/server/csr.json | cfssljson -bare certs/server/tls

mkdir -p certs/client
cat > certs/client/csr.json <<EOL
{
    "hosts": [
        "127.0.0.1",
        "localhost",
        "client"
    ],
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "CN": "client",
    "names": [
        {
            "C":  "US",
            "L":  "Utah",
            "O":  "UT Kubernetes",
            "OU": "Infra",
            "ST": "Utah"
        }
    ],
    "usages": [
        "server auth",
        "client auth"
    ]
}
EOL
cp certs/ca/tls.pem certs/client/ca.pem
cfssl gencert -ca certs/ca/tls.pem -ca-key certs/ca/tls-key.pem certs/client/csr.json | cfssljson -bare certs/client/tls
