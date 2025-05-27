package vault

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

// type HubManager struct {
// 	hubs map[string]*Hub
// 	mu   sync.RWMutex
// }

// var vaultHubManager = &HubManager{
// 	hubs: make(map[string]*Hub),
// }

// func (hm *HubManager) GetOrCreateHub(vaultID string) *Hub {
// 	hm.mu.Lock()
// 	defer hm.mu.Unlock()

// 	hub, exists := hm.hubs[vaultID]
// 	if !exists {
// 		hub = NewHub()
// 		hm.hubs[vaultID] = hub
// 		go hub.run() // start the goroutine
// 	}
// 	return hub
// }

// === Hub struct ===
type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan string
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan string),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (h *Hub) Run() {
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
	return func(wsConn *websocket.Conn) {
		hub.register <- wsConn
		defer func() {
			hub.unregister <- wsConn
		}()

		for {
			var msg string
			err := websocket.Message.Receive(wsConn, &msg)
			if err != nil {
				log.Println("Receive error:", err)
				break
			}

			var jsonMap map[string]interface{}
			json.Unmarshal([]byte(msg), &jsonMap)

			html := fmt.Sprintf(`
			<div id="messages" 
			  hx-swap-oob="afterbegin"
			  >	
			  <div class="bg-red-900 text-4xl">
			  %v

			  </div>
			</div>
			`, jsonMap["message"])
			hub.broadcast <- html
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
