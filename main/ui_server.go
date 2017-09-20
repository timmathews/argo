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
	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
	"github.com/satori/go.uuid"
	"github.com/timmathews/argo/config"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Uuid struct {
	Uuid []string `json:"uuid"`
}

type npmPackages struct {
	Objects []npmObject
}

type npmObject struct {
	Package npmPackage
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
		http.Error(w, "Could not generate UUID", 500)
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
			http.Error(w, "Could not parse data", 400)
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

			t.ExecuteTemplate(w, "layout", data.Objects)
		}
	} else {
		http.Error(w, "Method not allowed", 405)
	}
}

func appInstallHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var data pkgRequest
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Error("%v", err)
		} else {
			if err = installPackage(data.Package, data.Version); err != nil {
				log.Error("%v", err)
			}
		}
	} else {
		http.Error(w, "Method not allowed", 405)
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

func UiServer(addr *string, cmd chan CommandRequest) {
	r := mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(sysconf.Server.AssetPath))))
	r.HandleFunc("/admin/uuid", uuidHandler)
	r.HandleFunc("/admin", adminHandler)
	r.HandleFunc("/apps/install", appInstallHandler)
	r.HandleFunc("/apps", appsHandler)
	r.HandleFunc("/", indexHandler)
	http.Handle("/", r)
}
