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
				log.Fatal("errJsonor unmarshalling:", errJson)
				return
			}

			if payload.Action == "create_task" && payload.Path == fmt.Sprintf("/vaults/%v/tasks", vaultId) {
				fmt.Println(payload.Data["taskIdea"], "idea")
				fmt.Println(payload.Data["taskName"], "name")

				task, errTask := task.CreateTask(queries, payload.Data["task_name"], payload.Data["task_idea"], utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
				if errTask != nil {
					utils.CheckErrorWebsocket(hub.broadcast, "Duplicate idea in the same vault", errTask)
					return
				}

				html, errRender := templates.ReturnString("collaborative-task", addVaultIdToTask(vaultId, *task))
				if errRender != nil {
					utils.CheckErrorWebsocket(hub.broadcast, "Error rendering component", errRender)
					return
				}
				realHtml := wrapWithAttributesDiv(*html, `id="collaborative-tasks" hx-swap-oob="afterbegin"`)

				hub.broadcast <- realHtml
			}

			if payload.Action == "delete_task" {
				errTask := task.DeleteTask(queries, utils.PgUUID(payload.Data["task_id"]), utils.PgUUID(vaultId), utils.PgUUID(user.ID), r.Context())
				if errTask != nil {
					utils.CheckErrorWebsocket(hub.broadcast, "Internal server error", errTask)
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
					utils.CheckErrorWebsocket(hub.broadcast, "Can't have two tasks with the same name", errTask)
					return
				}

				html, errRender := templates.ReturnString("collaborative-task", addVaultIdToTask(vaultId, *task))
				if errRender != nil {
					utils.CheckErrorWebsocket(hub.broadcast, "Error rendering component", errRender)
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
					utils.CheckErrorWebsocket(hub.broadcast, "Internal server error", errTask)
					return
				}

				html, errRender := templates.ReturnString("collaborative-task", addVaultIdToTask(vaultId, *task))
				if errRender != nil {
					utils.CheckErrorWebsocket(hub.broadcast, "Error rendering component", errRender)
					return
				}
				realHtml := wrapWithAttributesDiv(*html, fmt.Sprintf(`id="task-%v" hx-swap-oob="outerHTML"`, payload.Data["task_id"]))

				hub.broadcast <- realHtml
			}

			if payload.Action == "add_tag_to_task" {
				task, errTask := task.AddTagToTask(
					queries,
					utils.PgUUID(payload.Data["tag_id"]),
					utils.PgUUID(user.ID),
					utils.PgUUID(payload.Data["task_id"]),
					utils.PgUUID(vaultId),

					r.Context())
				if errTask != nil {
					utils.CheckErrorWebsocket(hub.broadcast, "Can't have two tasks with the same name", errTask)
					return
				}

				html, errRender := templates.ReturnString("collaborative-task", addVaultIdToTask(vaultId, *task))
				if errRender != nil {
					utils.CheckErrorWebsocket(hub.broadcast, "Error rendering component", errRender)
					return
				}
				realHtml := wrapWithAttributesDiv(*html, fmt.Sprintf(`id="task-%v" hx-swap-oob="outerHTML"`, payload.Data["task_id"]))

				hub.broadcast <- realHtml
			}

			if payload.Action == "remove_tag_from_task" {
				task, errTask := task.DeleteTagRelationFromTask(
					queries,
					utils.PgUUID(payload.Data["tag_id"]),
					utils.PgUUID(user.ID),
					utils.PgUUID(vaultId),
					utils.PgUUID(payload.Data["task_id"]),

					r.Context())
				if errTask != nil {
					utils.CheckErrorWebsocket(hub.broadcast, "Can't have two tasks with the same name", errTask)
					return
				}

				html, errRender := templates.ReturnString("collaborative-task", addVaultIdToTask(vaultId, *task))
				if errRender != nil {
					utils.CheckErrorWebsocket(hub.broadcast, "Error rendering component", errRender)
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
					utils.CheckErrorWebsocket(hub.broadcast, "Internal server error", errJsonMarsh)
					return
				}

				hub.broadcast <- string(jsonCursor)
			}

			if payload.Action == "create_tag" {
				// MAKE THIS WORK BY SENDING sONE TAG NOT EVERYTHING
				if i, ok := payload.Data["source"]; ok && i == "/tags" {
					tagWithTasks, errTag := tag.CreateTag(queries,
						utils.PgUUID(user.ID),
						utils.PgUUID(vaultId),
						payload.Data["tag_name"],
						r.Context())
					if errTag != nil {
						utils.CheckErrorWebsocket(hub.broadcast, "Can't have two tags with the same name", errTag)
						return
					}
					html, errRender := templates.ReturnString("collaborative-tag", tagWithTasks)
					if errRender != nil {
						utils.CheckErrorWebsocket(hub.broadcast, "Error rendering component", errRender)
						return
					}
					realHtml := wrapWithAttributesDiv(*html, `id="collaborative-tags" hx-swap-oob="afterbegin"`)

					hub.broadcast <- realHtml
				} else {
					_, errTag := tag.CreateTag(queries,
						utils.PgUUID(user.ID),
						utils.PgUUID(vaultId),
						payload.Data["tag_name"],
						r.Context())
					if errTag != nil {
						utils.CheckErrorWebsocket(hub.broadcast, "Can't have two tags with the same name", errTag)
						return
					}

					filters := task.FilterParams(r)
					tasks, errTasks := task.GetFilteredTasksWithTags(queries,
						filters,
						utils.PgUUID(user.ID),
						utils.PgUUID(vaultId),
						r.Context(),
					)
					if errTasks != nil {
						utils.CheckErrorWebsocket(hub.broadcast, "Can't have two tags with the same name", errTasks)
						return
					}

					html, errRender := templates.ReturnString("collaborative-tasks-container", tasks)

					if errRender != nil {
						utils.CheckErrorWebsocket(hub.broadcast, "Error rendering component", errRender)
						return
					}
					realHtml := wrapWithAttributesDiv(*html, `id="collaborative-tasks-container" hx-swap-oob="outerHTML"`)

					hub.broadcast <- realHtml

				}

			}

			if payload.Action == "update_tag" {
				tag, errTag := tag.UpdateTag(
					queries,
					utils.PgUUID(payload.Data["tag_id"]),
					utils.PgUUID(user.ID),
					utils.PgUUID(vaultId),
					payload.Data["tag_name"],
					r.Context())

				if errTag != nil {
					utils.CheckErrorWebsocket(hub.broadcast, "Can't have two tasks with the same name", errTag)
					return
				}

				html, errRender := templates.ReturnString("collaborative-tag", tag)
				if errRender != nil {
					utils.CheckErrorWebsocket(hub.broadcast, "Error rendering component", errRender)
					return
				}
				realHtml := wrapWithAttributesDiv(*html, fmt.Sprintf(`id="tag-%v" hx-swap-oob="outerHTML"`, payload.Data["tag_id"]))

				hub.broadcast <- realHtml
			}
			if payload.Action == "delete_tag" {
				errTag := tag.DeleteTag(
					queries,
					utils.PgUUID(payload.Data["tag_id"]),
					utils.PgUUID(user.ID),
					utils.PgUUID(vaultId),
					r.Context())

				if errTag != nil {
					utils.CheckErrorWebsocket(hub.broadcast, "Can't have two tasks with the same name", errTag)
					return
				}

				realHtml := wrapWithAttributesDiv("", fmt.Sprintf(`id="tag-%v" hx-swap-oob="outerHTML"`, payload.Data["tag_id"]))

				hub.broadcast <- realHtml
			}

			if payload.Action == "add_task_to_tag" {
				tag, errTag := tag.AddTaskToTag(
					queries,
					utils.PgUUID(payload.Data["tag_id"]),
					utils.PgUUID(user.ID),
					utils.PgUUID(payload.Data["task_id"]),

					utils.PgUUID(vaultId),
					r.Context())

				if errTag != nil {
					utils.CheckErrorWebsocket(hub.broadcast, "Can't have two tasks with the same name", errTag)
					return
				}

				html, errRender := templates.ReturnString("collaborative-tag", tag)
				if errRender != nil {
					utils.CheckErrorWebsocket(hub.broadcast, "Error rendering component", errRender)
					return
				}
				realHtml := wrapWithAttributesDiv(*html, fmt.Sprintf(`id="tag-%v" hx-swap-oob="outerHTML"`, payload.Data["tag_id"]))

				hub.broadcast <- realHtml
			}
			if payload.Action == "remove_task_from_tag" {
				tag, errTag := tag.DeleteTaskRelationFromTag(
					queries,
					utils.PgUUID(payload.Data["tag_id"]),
					utils.PgUUID(payload.Data["task_id"]),
					utils.PgUUID(user.ID),

					utils.PgUUID(vaultId),
					r.Context())

				if errTag != nil {
					utils.CheckErrorWebsocket(hub.broadcast, "Can't have two tasks with the same name", errTag)
					return
				}

				html, errRender := templates.ReturnString("collaborative-tag", tag)
				if errRender != nil {
					utils.CheckErrorWebsocket(hub.broadcast, "Error rendering component", errRender)
					return
				}
				realHtml := wrapWithAttributesDiv(*html, fmt.Sprintf(`id="tag-%v" hx-swap-oob="outerHTML"`, payload.Data["tag_id"]))

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
