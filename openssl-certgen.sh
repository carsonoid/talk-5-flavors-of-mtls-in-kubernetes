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
            "L":  "Utah",CN
            "O":  "UT Kubernetes",
            "OU": "Infra",
            "ST": "Utah"
        }
    ]
}
EOL

cfssl genkey -initca csr.json | cfssljson -bare ca


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
    "CN": "server",CN
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

cfssl gencert -ca ca.pem -ca-key ca-key.pem server-csr.json | cfssljson -bare server

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

cfssl gencert -ca ca.pem -ca-key ca-key.pem client-csr.json | cfssljson -bare client
