package main

import "github.com/carsonoid/talk-all-the-mtls-in-k8s/internal/demo"

const msg = `
// START OMIT
Run it!
// END OMIT
`

const basePath = `.`
const script = `./cmd/demo/simple.sh`

func main() {
	demo.RunShellScript(basePath, script)
}
