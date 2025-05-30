package vault

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

var vaultHubManagerTest = &HubManager{
	hubs: make(map[string]*Hub),
}

// === WebSocket connection handler ===
func WsHandlerPubSubTest(vaultId string) func(*websocket.Conn) {
	return func(wsConn *websocket.Conn) {
		hub := vaultHubManagerTest.GetOrCreateHub(vaultId)
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
			errJson := json.Unmarshal([]byte(msg), &jsonMap)
			if errJson != nil {
				log.Println("errJsonor unmarshalling:", errJson)
				return
			}

			// Safely get the message string
			var message string
			if data, ok := jsonMap["data"].(map[string]interface{}); ok {
				if msgVal, ok := data["message"].(string); ok {
					message = msgVal
				}
			}

			html := fmt.Sprintf(`
			<div id="messages" 
			hx-swap-oob="afterbegin">	
			<div id="message" class="bg-red-900 text-4xl">
			%v
			</div>
			</div>
			`, message)

			hub.broadcast <- html
		}
	}
}

func DebugHandlerTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Connected hubs (vault websockets): %d\n", len(vaultHubManagerTest.hubs))
		for vaultId, hub := range vaultHubManagerTest.hubs {

			clientCount := len(hub.clients)
			fmt.Fprintf(w, "Connected clients: (%v) %d\n", vaultId, clientCount)
			for conn := range hub.clients {
				fmt.Fprintf(w, "- Client: %p\n", conn)
			}
		}
	}
}
