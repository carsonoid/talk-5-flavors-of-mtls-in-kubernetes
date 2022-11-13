package main

import (
	"os"
	"os/exec"
	"syscall"
)

const msg = `
// START OMIT
Run it!
// END OMIT
`

const basePath = "./mtls/0-none/"

func main() {
	err := os.Chdir(basePath)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("bash", "./cluster-up.sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
