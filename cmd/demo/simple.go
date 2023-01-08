package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

const msg = `
// START OMIT
Run it!
// END OMIT
`

const basePath = "."

func main() {
	err := os.Chdir(basePath)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("bash", "./cmd/demo/simple.sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	defer cmd.Process.Signal(syscall.SIGTERM)
	// send group signal to all processes in the process group
	defer syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
	fmt.Println("Done!")
}
