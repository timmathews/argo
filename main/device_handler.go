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
	"github.com/gorilla/mux"
	"github.com/timmathews/argo/request"
	"net/http"
	"time"
)

type ManufacturerInfo struct {
	CompanyName     string `json:"companyName"`
	ProductName     string `json:"productName"`
	SerialNumber    string `json:"serialNumber"`
	FirmwareVersion string `json:"firmwareVersion"`
	HardwareVersion string `json:"hardwareVersion"`
}

type AccessRequest struct {
	Permission string    `json:"permission"`
	Token      string    `json:"token"`
	Expiration time.Time `json:"expirationTime"`
}

type ClientRegistration struct {
	ClientId       string           `json:"clientId"`
	Description    string           `json:"description"`
	SignalKVersion string           `json:"signalKVersion"`
	Manufacturer   ManufacturerInfo `json:"manufacturer"`
	request        AccessRequest
}

type Response struct {
	State         requestState.State
	Status        int32 `json:"statusCode"`
	Message       string
	Url           string `json:"href"`
	AccessRequest AccessRequest
}

func (c *ClientRegistration) GetStatus() AccessRequest {
	return c.request
}

func (c *ClientRegistration) Grant(until time.Time) {
	c.request.Expiration = until
	c.request.Permission = "APPROVED"
	c.request.Token = "SampleTokenTotesNotValid"
}

func (c *ClientRegistration) Renew(until time.Time) {
	c.request.Expiration = until
}

func (c *ClientRegistration) Revoke() {
	c.request.Expiration = time.Unix(0, 0)
	c.request.Permission = "DENIED"
	c.request.Token = ""
}

func (c *ClientRegistration) IsValid() bool {
	return c.request.Expiration.After(time.Now()) &&
		c.request.Permission == "APPROVED"
}

func registrationHandler(cmd chan ClientRegistration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var request ClientRegistration
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&request)

			if err != nil {
				log.Warning("Invalid JSON request")
				fmt.Fprintf(w, "Invalid JSON")
			}

			log.Debugf("Client: %s", request.ClientId)
			log.Debugf("Description: %s", request.Description)
			log.Debugf("Signal K Version: %s", request.SignalKVersion)
			log.Debugf("Manufacturer: %s", request.Manufacturer.CompanyName)
			log.Debugf("Product: %s", request.Manufacturer.ProductName)
			log.Debugf("Serial Number: %s", request.Manufacturer.SerialNumber)
			log.Debugf("Firmware: %s", request.Manufacturer.FirmwareVersion)
			log.Debugf("Hardware: %s", request.Manufacturer.HardwareVersion)

			cmd <- request
		}
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
}

func DeviceRegistrationServer(addr *string, cmd chan ClientRegistration) {
	r := mux.NewRouter()
	s := r.PathPrefix("/signalk/v1/access").Subrouter()
	s.HandleFunc("/requests", registrationHandler(cmd))
	s.HandleFunc("/requests/{key}", statusHandler)
	http.Handle("/signalk/v1/access/", r)
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		s, _ := route.GetPathTemplate()
		log.Notice("%v", s)

		return nil
	})
}
