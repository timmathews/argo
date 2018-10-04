package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

type yarnMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func installPackage(name, version string) error {
	cmd := exec.Command(
		"yarn",
		"add",
		fmt.Sprintf("%s@%s", name, version),
		"--json",
		"--silent",
	)
	chn := make(chan string)

	go func(cm *exec.Cmd, ch chan string) {
		defer func() { ch <- "" }()

		stdout, _ := cmd.StdoutPipe()

		<-ch

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			m := scanner.Text()
			var res yarnMessage
			json.Unmarshal([]byte(m), &res)
			if res.Type == "step" {
				if v, ok := res.Data.(map[string]interface{}); ok {
					c, ok1 := v["current"].(float64)
					t, ok2 := v["total"].(float64)

					if ok1 && ok2 {
						log.Noticef("Progress %v/%v", c, t)
					}
				} else {
					log.Error("Could not convert 'data'")
				}
			}
		}
	}(cmd, chn)

	chn <- ""
	cmd.Start()

	return cmd.Wait()
}

func getPathForPackage(pkg string) string {
	path := filepath.Join("node_modules", pkg, "index.html")
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return filepath.Join("/apps", pkg, "public")
	}

	return filepath.Join("/apps", pkg)
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
						if p.Location == "" {
							p.Location = fmt.Sprintf("/%v", p.Name)
						}

						_, err := os.Stat(filepath.Join(path, file.Name(), "index.html"))
						if os.IsNotExist(err) {
							p.Location = filepath.Join(p.Location, "/public")
						}

						// strip node_modules from path, if it exists
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
