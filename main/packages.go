package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
)

type nodePackage struct {
	Version  string   `json:"version"`
	Name     string   `json:"name"`
	Main     string   `json:"main"`
	Location string   `json:"_location"`
	Keywords []string `json:"keywords"`
}

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

func getInstalledPackages(path string) map[string]nodePackage {
	packages := make(map[string]nodePackage)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}

	for _, file := range files {
		if file.IsDir() {
			if strings.HasPrefix(file.Name(), "@") {
				for k, v := range getInstalledPackages(filepath.Join(path, file.Name())) {
					packages[k] = v
				}
			} else if file.Name() != ".bin" {
				pkgFile, err := ioutil.ReadFile(filepath.Join(path, file.Name(), "package.json"))
				if err == nil {
					var p nodePackage
					err := json.Unmarshal(pkgFile, &p)
					if err == nil && contains(p.Keywords, "signalk-webapp") {
						if strings.Contains(path, "node_modules") {
							packages[filepath.Join(path, file.Name())[13:]] = p
						} else {
							packages[filepath.Join(path, file.Name())] = p
						}
					}
				} else {
					return nil
				}
			}
		}
	}

	return packages
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
