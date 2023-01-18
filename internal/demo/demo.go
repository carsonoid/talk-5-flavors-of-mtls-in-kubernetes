package demo

import (
	"os"
	"os/exec"
	"syscall"
)

// RunShellScript runs a shell script in a given directory
func RunShellScript(basePath string, script string) error {
	cwd := os.Getenv("PWD")
	defer os.Chdir(cwd)

	err := os.Chdir(basePath)
	if err != nil {
		return err
	}

	tmpFile, err := os.Create("tmp.sh")
	if err != nil {
		return err
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(script)
	if err != nil {
		return err
	}

	cmd := exec.Command("bash", tmpFile.Name())
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		return err
	}

	defer cmd.Process.Signal(syscall.SIGTERM)
	// send group signal to all processes in the process group
	defer syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}
