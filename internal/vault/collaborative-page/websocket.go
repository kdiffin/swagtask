package vault

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/tag"
	"swagtask/internal/task"
	"swagtask/internal/template"
	"swagtask/internal/utils"
	"sync"

	"golang.org/x/net/websocket"
)

type HubManager struct {
	hubs map[string]*Hub
	mu   sync.RWMutex
}

func (hm *HubManager) GetOrCreateHub(vaultID string) *Hub {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	hub, exists := hm.hubs[vaultID]
	if !exists {
		hub = NewHub()
		hm.hubs[vaultID] = hub
		go hub.Run() // start the goroutine
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

var vaultHubManager = &HubManager{
	hubs: make(map[string]*Hub),
}

type Payload struct {
	Action string            `json:"action"`
	Path   string            `json:"path"`
	Data   map[string]string `json:"data"` // keep this flexible if `data` varies
}

// === WebSocket connection handler ===
func WsHandlerPubSub(queries *db.Queries, templates *template.Template, w http.ResponseWriter, r *http.Request) func(*websocket.Conn) {
	vaultId, _ := middleware.VaultIDFromContext(r.Context())
	user, _ := middleware.UserFromContext(r.Context())

	return func(wsConn *websocket.Conn) {
		hub := vaultHubManager.GetOrCreateHub(vaultId)
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

			var payload Payload
			errJson := json.Unmarshal([]byte(msg), &payload)
			if errJson != nil {
				log.Println("errJsonor unmarshalling:", errJson)
				return
			}

			if payload.Action == "create_task" && payload.Path == fmt.Sprintf("/vaults/%v/tasks", vaultId) {
				task, errTask := task.CreateTask(queries, payload.Data["task_name"], payload.Data["task_idea"], utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
				if errTask != nil {
					utils.LogError("error at websocket", errTask)
					return
				}

				html, errRender := templates.ReturnString("collaborative-task", addVaultIdToTask(vaultId, *task))
				if errRender != nil {
					utils.LogError("error at websocket", errRender)
					return
				}
				realHtml := wrapWithAttributesDiv(*html, `id="collaborative-tasks" hx-swap-oob="afterbegin"`)

				hub.broadcast <- realHtml
			}

			if payload.Action == "delete_task" {
				errTask := task.DeleteTask(queries, utils.PgUUID(payload.Data["task_id"]), utils.PgUUID(vaultId), utils.PgUUID(user.ID), r.Context())
				if errTask != nil {
					utils.LogError("error at websocket", errTask)
					return
				}

				realHtml := wrapWithAttributesDiv("", fmt.Sprintf(`id="task-%v" hx-swap-oob="outerHTML"`, payload.Data["task_id"]))

				hub.broadcast <- realHtml
			}

			if payload.Action == "update_task" {
				task, errTask := task.UpdateTask(
					queries,
					utils.PgUUID(vaultId),
					utils.PgUUID(payload.Data["task_id"]),
					utils.PgUUID(user.ID),
					payload.Data["task_name"],
					payload.Data["task_idea"],

					r.Context())
				if errTask != nil {
					utils.LogError("error at websocket", errTask)
					return
				}

				html, errRender := templates.ReturnString("collaborative-task", addVaultIdToTask(vaultId, *task))
				if errRender != nil {
					utils.LogError("error at websocket", errRender)
					return
				}
				realHtml := wrapWithAttributesDiv(*html, fmt.Sprintf(`id="task-%v" hx-swap-oob="outerHTML"`, payload.Data["task_id"]))

				hub.broadcast <- realHtml
			}
			if payload.Action == "update_task_completion" {
				task, errTask := task.UpdateTaskCompletion(
					queries,
					utils.PgUUID(user.ID),
					utils.PgUUID(vaultId),
					utils.PgUUID(payload.Data["task_id"]),
					r.Context())
				if errTask != nil {
					utils.LogError("error at websocket", errTask)
					return
				}

				html, errRender := templates.ReturnString("collaborative-task", addVaultIdToTask(vaultId, *task))
				if errRender != nil {
					utils.LogError("error at websocket", errRender)
					return
				}
				realHtml := wrapWithAttributesDiv(*html, fmt.Sprintf(`id="task-%v" hx-swap-oob="outerHTML"`, payload.Data["task_id"]))

				hub.broadcast <- realHtml
			}

			if payload.Action == "move_cursor" {
				var cursor struct {
					Type     string
					X        string `json:"x"`
					Y        string `json:"y"`
					Username string `json:"username"`
				}

				cursor.Username = user.Username
				cursor.X = payload.Data["x"]
				cursor.Y = payload.Data["y"]
				cursor.Type = "cursor_info"

				jsonCursor, errJsonMarsh := json.Marshal(cursor)
				if errJsonMarsh != nil {
					utils.LogError("", errJsonMarsh)
					return
				}

				hub.broadcast <- string(jsonCursor)
			}

			if payload.Action == "create_tag" {
				// MAKE THIS WORK BY SENDING ONE TAG NOT EVERYTHING
				fmt.Println(user.DefaultVaultID)
				err := queries.CreateTag(r.Context(), db.CreateTagParams{
					Name:    payload.Data["tag_name"],
					UserID:  utils.PgUUID(user.ID),
					VaultID: utils.PgUUID(user.DefaultVaultID),
				})
				if utils.CheckError(w, r, err) {
					fmt.Println("error was here1")
					return
				}

				tagsWithTasks, errTags := tag.GetTagsWithTasks(queries, utils.PgUUID(user.ID), utils.PgUUID(user.DefaultVaultID), r.Context())
				if utils.CheckError(w, r, errTags) {
					fmt.Println("error was here")
					return
				}
				html, errRender := templates.ReturnString("collaborative-tags-list-container", tagsWithTasks)
				if errRender != nil {
					utils.LogError("error at websocket", errRender)
					return
				}
				realHtml := wrapWithAttributesDiv(*html, `id="collaborative-tags-list-container" hx-swap-oob="outerHTML"`)

				hub.broadcast <- realHtml
			}

		}
	}
}
func addVaultIdToTask(vaultId string, t task.TaskWithTags) task.TaskWithTags {
	tasksReal := task.TaskWithTags{
		TaskUI: task.TaskUI{
			ID:        t.ID,
			Name:      t.Name,
			Author:    t.Author,
			Idea:      t.Idea,
			Completed: t.Completed,
		},
		VaultID:       vaultId,
		RelatedTags:   t.RelatedTags,
		AvailableTags: t.AvailableTags,
	}
	return tasksReal
}
func wrapWithAttributesDiv(html string, attrs string) string {
	s := fmt.Sprintf(`<div %v>`, attrs) + html + "</div>"

	return s
}

func DebugHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Connected hubs (vault websockets): %d\n", len(vaultHubManager.hubs))
		for vaultId, hub := range vaultHubManager.hubs {

			clientCount := len(hub.clients)
			fmt.Fprintf(w, "Connected clients: (%v) %d\n", vaultId, clientCount)
			for conn := range hub.clients {
				fmt.Fprintf(w, "- Client: %p\n", conn)
			}
		}
	}
}
