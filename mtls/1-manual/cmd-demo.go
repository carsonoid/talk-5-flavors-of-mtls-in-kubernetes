package main

import "github.com/carsonoid/talk-all-the-mtls-in-k8s/internal/demo"

const basePath = `./mtls/1-manual/`
const script = `
#!/bin/bash

set -e

# kill any running servers listening on 8443
kill $(lsof -t -i:8443) || true

# START OMIT
# Usage: secure-server CERTFILE KEYFILE CAFILE
go run ../../cmd/secure-server/ \
    certs/server/tls.pem \
    certs/server/tls-key.pem \
    certs/ca/tls.pem &

# Wait for the server to start
sleep 3

# Usage: secure-client CERTFILE KEYFILE CAFILE SERVERURL
go run ../../cmd/secure-client/ \
    certs/client/tls.pem \
    certs/client/tls-key.pem \
    certs/ca/tls.pem \
    https://localhost:8443
# END OMIT
`

func main() {
	demo.RunShellScript(basePath, script)
}
