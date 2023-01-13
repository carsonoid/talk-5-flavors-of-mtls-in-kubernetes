package main

import "github.com/carsonoid/talk-all-the-mtls-in-k8s/internal/demo"

const basePath = `./mtls/1-manual/`
const script = `
set -e

// START OMIT
find certs -type f
// END OMIT
`

func main() {
	demo.RunShellScript(basePath, script)
}
