package vault

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

type HubManager struct {
	hubs map[string]*Hub
	mu   sync.RWMutex
}

var vaultHubManager = &HubManager{
	hubs: make(map[string]*Hub),
}

func (hm *HubManager) GetOrCreateHub(vaultID string) *Hub {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	hub, exists := hm.hubs[vaultID]
	if !exists {
		hub = newHub()
		hm.hubs[vaultID] = hub
		go hub.run() // start the goroutine
	}
	return hub
}

// === Hub struct ===
type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan string
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan string),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (h *Hub) run() {
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			h.clients[conn] = true
			h.mu.Unlock()
		case conn := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				conn.Close()
			}
			h.mu.Unlock()
		case msg := <-h.broadcast:
			h.mu.Lock()
			for conn := range h.clients {
				err := websocket.Message.Send(conn, msg)
				if err != nil {
					log.Println("Send error:", err)
					conn.Close()
					delete(h.clients, conn)
				}
			}
			h.mu.Unlock()
		}
	}
}

// === WebSocket connection handler ===
func WsHandler(hub *Hub) func(*websocket.Conn) {
	return func(ws *websocket.Conn) {
		hub.register <- ws
		defer func() {
			hub.unregister <- ws
		}()

		for {
			var msg string
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				log.Println("Receive error:", err)
				break
			}
			hub.broadcast <- msg
		}
	}
}

func DebugHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hub.mu.Lock()
		defer hub.mu.Unlock()

		clientCount := len(hub.clients)
		fmt.Fprintf(w, "Connected clients: %d\n", clientCount)

		for conn := range hub.clients {
			fmt.Fprintf(w, "- Client: %p\n", conn)
		}
	}
}

func ws() {
	hub := newHub()
	go hub.run()

	http.Handle("/ws", websocket.Handler(WsHandler(hub)))
	http.HandleFunc("/debug", DebugHandler(hub)) // <<== here

	fmt.Println("🔥 Running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
