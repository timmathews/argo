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
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
	"net/http"
	"time"
)

type connection struct {
	ws *websocket.Conn

	send chan []byte
}

type hub struct {
	// Known connections
	connections map[*connection]bool

	// Messages for the connections
	broadcast chan []byte

	// Register requests from new connections
	register chan *connection

	// Unregister requests from closing connections
	unregister chan *connection
}

const (
	// Time allowed to write a message to a receiver
	writeWait = 10 * time.Second

	// Time allowed to read pong from a receiver
	pongWait = 60 * time.Second

	// Send pings to receivers
	pingPeriod = 56 * time.Second

	// Max message size
	maxMessageSize = 512
)

var websocket_hub = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

var statistics_hub = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	// TODO: Add origin check / CORS support

	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Error("Websocket:", err)
		return
	}

	if r.URL.Path == "/signalk/v1/stream" {
		c := &connection{send: make(chan []byte, 256), ws: ws}
		websocket_hub.register <- c
		go c.writePump()
	} else if r.URL.Path == "/signalk/v1/control" {
		fmt.Println("Got registration")
	}
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Error("Websocket:", err)
		return
	}

	c := &connection{send: make(chan []byte, 256), ws: ws}
	statistics_hub.register <- c
	go c.writePump()
}

func loggingHandler(handler http.Handler, log *logging.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Noticef("%v %v %v", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func WebSocketServer(addr *string, log *logging.Logger) {
	http.HandleFunc("/signalk/v1/", serveWs)
	http.HandleFunc("/ws/stats", handleStats)
	err := http.ListenAndServe(*addr, loggingHandler(http.DefaultServeMux, log))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
