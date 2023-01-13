
#!/bin/bash

set -e

trap "pkill -9 -f ^/tmp/go-build.+/.+" EXIT



# Usage: secure-server CERTFILE KEYFILE CAFILE
go run ./cmd/secure-server/ \
    mtls/1-manual/certs/server/tls.pem \
    mtls/1-manual/certs/server/tls-key.pem \
    mtls/1-manual/certs/ca/tls.pem &

# Wait for the server to start
sleep 3

# Usage: secure-client CERTFILE KEYFILE CAFILE SERVERURL
go run ./cmd/secure-client/ \
    mtls/1-manual/certs/client/tls.pem \
    mtls/1-manual/certs/client/tls-key.pem \
    mtls/1-manual/certs/ca/tls.pem \
    https://localhost:8443



wait
