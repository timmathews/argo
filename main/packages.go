package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func installPackage(name, version string) error {
	cmd := exec.Command("npm", "install", fmt.Sprintf("%s@%s", name, version), "--prefix ./vendor")

	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return err
	}

	processErrors, _ := ioutil.ReadAll(stderr)
	processOutput, _ := ioutil.ReadAll(stdout)

	if len(processErrors) > 0 {
		log.Error("Process Errors: %v", string(processErrors))
	}

	if len(processOutput) > 0 {
		log.Notice("%v", string(processOutput))
	}

	return cmd.Wait()
}
