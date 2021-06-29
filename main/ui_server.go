/*
 * Copyright (C) 2016 Tim Mathews <tim@signalk.org>
 *
 * This file is part of Argo.
 *
 * Argo is free software: you can redistribute it and/or modify it under the
 * terms of the GNU General Public License as published by the Free Software
 * Foundation, either version 3 of the License, or (at your option) any later
 * version.
 *
 * Argo is distributed in the hope that it will be useful, but WITHOUT ANY
 * WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
 * FOR A PARTICULAR PURPOSE. See the GNU General Public License for more
 * details.
 *
 * You should have received a copy of the GNU General Public License along with
 * this program. If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
	uuid "github.com/satori/go.uuid"
	"github.com/timmathews/argo/config"
)

type Uuid struct {
	Uuid []string `json:"uuid"`
}

type npmPackages struct {
	Objects []npmObject
}

type npmObject struct {
	Package npmPackage
	Path    string
	Version string
}

type npmLinks struct {
	Npm        string
	Homepage   string
	Repository string
}

type npmAuthor struct {
	Name  string
	Email string
}

type npmPackage struct {
	Name        string
	Version     string
	Description string
	Links       npmLinks
	Author      npmAuthor
}

type pkgRequest struct {
	Package string
	Version string
}

var pages *template.Template

func uuidHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	uuid := Uuid{
		Uuid: strings.Split(strings.ToUpper(uuid.NewV4().String()), "-"),
	}

	b, err := json.Marshal(uuid)

	if err != nil {
		log.Error("JSON Error:", err)
		http.Error(w, "Could not generate UUID", http.StatusInternalServerError)
	} else {
		io.WriteString(w, string(b))
	}
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var data config.TomlConfig

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Error("Parsing %v", err)
			http.Error(w, "Could not parse data", http.StatusBadRequest)
		} else {
			log.Notice("\n%v", data)
			err = mergo.Merge(&sysconf, data)
			if err == nil {
				tmp, _ := ioutil.ReadFile("argo.conf")
				ioutil.WriteFile("argo.conf~", tmp, 0644)
				config.WriteConfig("argo.conf", sysconf)
			} else {
				log.Error("%v", err)
			}
		}
	}
}

func appsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("templates/layout.gtpl", "templates/apps.gtpl")
		if err != nil {
			log.Error("%v", err)
		} else {
			resp, err := http.Get("http://registry.npmjs.org/-/v1/search?size=250&text=keywords:signalk-webapp")
			if err != nil {
				log.Error("%v", err)
				return
			}
			defer resp.Body.Close()

			var data npmPackages
			err = json.NewDecoder(resp.Body).Decode(&data)
			if err != nil {
				log.Error("%v", err)
				return
			}

			installedPackages := getInstalledPackages("./node_modules")

			for k, v := range data.Objects {
				log.Notice("%v", v.Package.Name)
				if val, ok := installedPackages[v.Package.Name]; ok {
					data.Objects[k].Path = "/apps" + val.Location
					data.Objects[k].Version = val.Version
				}
			}

			t.ExecuteTemplate(w, "layout", data.Objects)
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func appInstallHandlerFactory(m *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var data pkgRequest
			err := json.NewDecoder(r.Body).Decode(&data)
			w.Header().Set("Content-Type", "text/json")
			if err != nil {
				log.Error("%v", err)
			} else {
				if err = installPackage(data.Package, data.Version); err != nil {
					log.Error("%v", err)
					io.WriteString(
						w,
						fmt.Sprintf(`{"result":"error","message":"%v"}`, err),
					)
				} else {
					log.Noticef("%v@%v installed", data.Package, data.Version)
					dir := http.Dir("./node_modules")
					path := getPathForPackage(data.Package)
					m.PathPrefix(path).Handler(
						http.StripPrefix("/apps/", http.FileServer(dir)),
					)
					io.WriteString(
						w,
						fmt.Sprintf(`{"result":"ok","url":"%v","dir":"%v"}`, path, dir),
					)
				}
			}

		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("templates/layout.gtpl", "templates/index.gtpl")
		if err != nil {
			log.Error("%v", err)
		} else {
			t.ExecuteTemplate(w, "layout", sysconf)
		}
	}
}

func n2kHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("templates/layout.gtpl", "templates/n2k.gtpl")
		if err != nil {
			log.Error("%v", err)
		} else {
			t.ExecuteTemplate(w, "layout", nil)
		}
	}
}

func addApps(r *mux.Router) {
	pkgs := getInstalledPackages("./node_modules")
	for k, _ := range pkgs {
		log.Notice("%v", k)
		dir := http.Dir("./node_modules")
		r.PathPrefix("/apps/" + k).Handler(http.StripPrefix("/apps/", http.FileServer(dir)))
	}
}

func UiServer(addr *string, cmd chan CommandRequest) {
	r := mux.NewRouter()
	addApps(r)
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(sysconf.Server.AssetPath))))
	r.HandleFunc("/admin/uuid", uuidHandler)
	r.HandleFunc("/admin", adminHandler)
	r.HandleFunc("/apps/install", appInstallHandlerFactory(r))
	r.HandleFunc("/apps", appsHandler)
	r.HandleFunc("/n2k", n2kHandler)
	r.HandleFunc("/", indexHandler)
	http.Handle("/", r)

	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		s, _ := route.GetPathTemplate()
		log.Notice("%v", s)

		return nil
	})
}
