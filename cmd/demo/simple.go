package main

import "github.com/carsonoid/talk-all-the-mtls-in-k8s/internal/demo"

const basePath = `.`
const script = `
#!/bin/bash

set -e

# kill any running servers listening on 8080
kill $(lsof -t -i:8080) || true


# START OMIT
#!/bin/bash

go run ./cmd/insecure-server/ &

# Wait for the server to start
sleep 3

go run ./cmd/insecure-client/ http://localhost:8080
# END OMIT
`

func main() {
	demo.RunShellScript(basePath, script)
}
