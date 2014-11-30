package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

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

type connection struct {
	ws *websocket.Conn

	send chan []byte
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
		log.Println(err)
		return
	}

	if r.URL.Path == "/ws/v1/data" {
		c := &connection{send: make(chan []byte, 256), ws: ws}
		h.register <- c
		go c.writePump()
	} else if r.URL.Path == "/ws/v1/control" {
		fmt.Println("Got registration")
	}
}
