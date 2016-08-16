package main

import (
	"github.com/op/go-logging"
	"net/http"
)

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

var websocket_hub = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
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

func loggingHandler(handler http.Handler, log *logging.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("%v %v %v", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func WebSocketServer(addr *string, log *logging.Logger) {
	http.HandleFunc("/ws/v1/", serveWs)
	err := http.ListenAndServe(*addr, loggingHandler(http.DefaultServeMux, log))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
