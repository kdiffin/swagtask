REFACTOR TEST WITH CLAUDE 

# Refactor Instructions: swagtask Go Codebase
# Handler DI + Struct Receiver Pattern

## Context

Codebase: Go backend, HTMX frontend, sqlc for DB queries
Current problem: handlers are standalone functions with deps passed as args every call
Goal: convert to struct receiver pattern so handlers own their deps via DI

Current bad pattern:
```go
func HandlerCreateTask(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template)
```

Target pattern:
```go
type TaskHandler struct {
    queries   *db.Queries
    templates *template.Template
}
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request)
```

---

## Project Structure

```
internal/
  auth/
    handler.go
    service.go
    models.go
    utils.go
  task/
    handler.go
    service.go
    models.go
  tag/
    handler.go
    service.go
    models.go
  vault/
    handler.go
    handler-tasks.go
    handler-tags.go
    websocket.go
    common/
  middleware/
  db/generated/   ← sqlc generated, DO NOT MODIFY
  template/
  utils/
router/
cmd/swagtask/
```

Each noun package (auth, task, tag, vault) has its own handler file. Vault is special — it imports from task and tag packages, acting as an orchestration layer.

---

## Phase 1: Audit

Before touching any code, do the following:

1. Read every handler file in: auth, task, tag, vault
2. For each package, list:
   - All handler function signatures
   - Which deps each handler uses (queries, templates, anything else)
   - Any handler that uses deps from OTHER packages (e.g. vault calling task.GetFilteredTasksWithTags)
3. Read router/ to understand how handlers are currently registered
4. Note any handlers that have special error handling logic (vault has multi-render patterns)
5. Check if any handler uses deps beyond queries + templates (e.g. env vars, external clients)

Do NOT make changes in Phase 1. Audit only.

---

## Phase 2: Per-Package Handler Structs

Do one package at a time. Order: task → tag → auth → vault (vault last, it depends on others)

For each package:

### Step 1 — Define the struct

```go
type TaskHandler struct {
    queries   *db.Queries
    templates *template.Template
}

func NewTaskHandler(queries *db.Queries, templates *template.Template) *TaskHandler {
    return &TaskHandler{queries: queries, templates: templates}
}
```

Place this at the top of handler.go in that package.

### Step 2 — Convert each function

Before:
```go
func HandlerCreateTask(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
    // uses queries and templates
}
```

After:
```go
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
    // replace queries → h.queries
    // replace templates → h.templates
}
```

Naming convention for methods: drop the "Handler" prefix and package name
- HandlerCreateTask → Create
- HandlerGetTask → Get (or GetByID if ambiguous)
- HandlerDeleteTask → Delete
- HandlerTaskToggleComplete → ToggleComplete
- HandlerUpdateTask → Update

### Step 3 — Verify method signatures

Every method must match http.HandlerFunc exactly:
```go
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request)
```
No extra params. This is the whole point.

---

## Phase 3: Vault Package (Special Case)

Vault is an orchestration layer — its handlers call service functions from task and tag packages directly. Its struct needs to hold those deps too.

```go
type VaultHandler struct {
    queries   *db.Queries
    templates *template.Template
}
```

Vault handlers that call task/tag functions (e.g. task.GetFilteredTasksWithTags) should continue doing so — this is fine and correct. Vault owns the cross-cutting logic. Do not inline or move those calls.

Vault has multiple handler files (handler.go, handler-tags.go, handler-tasks.go). All methods go on the same VaultHandler struct. Define the struct once in handler.go, methods can live in any file in the vault package.

Websocket handlers: handle separately. Check if they need the same deps or different ones. If different, they can be methods on VaultHandler too, or a separate WebsocketHandler struct — use judgment based on what deps they actually use.

---

## Phase 4: Router Wiring

After all packages are converted, update the router.

Current pattern (guessing based on architecture):
```go
mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
    task.HandlerGetTasks(w, r, queries, templates)
})
```

Target pattern:
```go
taskHandler := task.NewTaskHandler(queries, templates)
tagHandler := tag.NewTagHandler(queries, templates)
authHandler := auth.NewAuthHandler(queries, templates)
vaultHandler := vault.NewVaultHandler(queries, templates)

mux.HandleFunc("GET /tasks", taskHandler.GetAll)
mux.HandleFunc("POST /tasks", taskHandler.Create)
// etc
```

Instantiate all handlers once at startup, wire into router. Deps flow in once at the top, not on every request.

---

## Phase 5: Cleanup

After router is updated and app compiles:

1. Remove any fmt.Println debug logs found during audit (vault handlers had some)
2. Look for duplicated logic across handlers in the same package — extract to private helper if found more than twice
3. The tasksReal mapping loop appears duplicated in vault — extract to a helper function in vault package
4. The collaborator role loop also duplicated in vault — extract similarly

---

## Rules and Constraints

- DO NOT modify anything in internal/db/generated/ — sqlc generated, hands off
- DO NOT change service.go files — only handler.go files are in scope
- DO NOT introduce interfaces for *db.Queries in this refactor — that's a separate future task. Keep *db.Queries as concrete type for now
- DO NOT change template rendering calls — templates.Render(w, "name", data) stays as-is
- KEEP error handling patterns identical — utils.CheckError pattern stays
- KEEP middleware patterns identical — middleware.UserFromContext, middleware.VaultIDFromContext stay
- ONE package at a time — compile and verify each package before moving to next
- If something is ambiguous, err on the side of doing less and noting what needs human review

---

## Verification Checklist Per Package

After each package conversion:
- [ ] Struct defined with correct fields
- [ ] NewXHandler constructor exists
- [ ] All handler functions converted to methods
- [ ] No handler function still has queries/templates as params
- [ ] All method signatures are (w http.ResponseWriter, r *http.Request)
- [ ] Package compiles: `go build ./internal/packagename/...`

After router update:
- [ ] All handlers instantiated once at startup
- [ ] All routes updated to use method syntax
- [ ] Full build passes: `go build ./...`
- [ ] No unused imports

---

## Notes on This Codebase

- Uses sqlc so all DB calls go through *db.Queries — no raw SQL in handlers
- HTMX frontend means handlers often render partial templates, not full pages — preserve this
- Vault intentionally depends on task and tag packages — this is correct architecture, not a smell
- The codebase has some known tech debt (template replication comments) — do not fix those, out of scope
