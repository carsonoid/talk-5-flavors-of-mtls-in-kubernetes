package main

import "github.com/carsonoid/talk-all-the-mtls-in-k8s/internal/demo"

const basePath = `.`
const script = `
set -e

// START OMIT
ls -l certs/*
// END OMIT
`

func main() {
	demo.RunShellScript(basePath, script)
}
