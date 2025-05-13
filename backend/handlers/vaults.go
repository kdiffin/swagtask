package handlers

import (
	"log"
	"net/http"
	db "swagtask/db/generated"
	"swagtask/models"

	"golang.org/x/net/websocket"
)



func WsHandler() func(*websocket.Conn) {
	return func(ws *websocket.Conn) {
		for {
			var msg string
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				log.Println("Receive error:", err)
				break
			}

			errSend := websocket.Message.Send(ws, msg)
			if errSend != nil {
				log.Println("Receive error:", err)
				break
			}

			
		}
	}
}

func HandlerGetVaults(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *models.Template,) {
	
	templates.Render(w, "vaults-page", nil)
}