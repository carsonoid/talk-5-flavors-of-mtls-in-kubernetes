#!/bin/bash

# Generate CA key and cert

cat > ca-csr.json <<EOL
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

mkdir -p certs/ca
cfssl genkey -initca ca-csr.json | cfssljson -bare certs/ca/tls

cat > server-csr.json <<EOL
{
    "hosts": [
        "127.0.0.1",
        "localhost"
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
    ]
}
EOL

mkdir -p certs/server
cp certs/ca/tls.pem certs/server/ca.pem
cfssl gencert -ca certs/ca/tls.pem -ca-key certs/ca/tls-key.pem server-csr.json | cfssljson -bare certs/server/tls

cat > client-csr.json <<EOL
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
    ]
}
EOL

mkdir -p certs/client
cp certs/ca/tls.pem certs/client/ca.pem
cfssl gencert -ca certs/ca/tls.pem -ca-key certs/ca/tls-key.pem client-csr.json | cfssljson -bare certs/client/tls
