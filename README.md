## TODOS:

- [x] finish the main page and abstract the update logic of the tags
- [x] split the task updating and tag creation/relation endpoints
- [x] make ui composable
- [x] implement the get task by tags page
- [x] implement manage tags page
  - [x] implement CRUD on tags, make it a card based ui for CRUDing tags
- [x] make the dropdowns for adding relations only show items which havent been added already as tags
- [x] add many to many relationship thing and show tasks related to tag in the card
- [x] implement tags page tag relation to task
- [x] make the footer work
- [x] IMPLEMENT ACTUAL FILTERING ON TASKS PAGE

  - [x] implement searching by name
  - [x] implement filtering by tag
  - [ ] implement filtering by completed (cba)

- [ ] read the other md files rq
- [ ] implement actual database migrations and logging
- [ ] add sqlc
  - [ ] 
- [ ] add net/http and remove echo
- [ ] refactor everything
- [ ] add barebones auth using cookies then jwt
- [ ] add live editing via websockets or polling between 2 users

### long term

- [ ] add openapi specs

---

chat gippity made the readme
im doing this to learn htmx and golang i hate react now im a based htmxpilled dev now f the virtual dom
chat gippity made my frontends styles btw

---

Aight 😤 let’s cook something that's:

- 🔥 fun to build
- 🧠 hits core web fundamentals (routing, templates, components, state)
- 💡 has room to grow into full-stack if you want later
- 🛠️ perfect for your **Go + HTMX + Tailwind** stack

---

## 🎯 Project Idea: **Task Vault**

> A keyboard-driven task manager with nested tasks, tagging, filters, and HTMX-powered inline editing.

### 🧩 Features (V1)

- 📝 Add / Edit / Delete tasksw
- 🏷️ Tag tasks with colors (e.g., `#school`, `#project`, `#bug`)
- ✅ Mark tasks complete
- 🧵 Nest sub-tasks inside parent tasks
- 🔍 HTMX-powered search and tag filter
- 💥 All inline — no page reloads

---

### 📐 Pages

| URL          | Purpose                |
| ------------ | ---------------------- |
| `/`          | Dashboard w/ all tasks |
| `/tasks/:id` | Focus view on one task |
| `/tags/:tag` | Filtered view by tag   |
| `/about`     | (optional) About page  |

---

### 🧠 What You’ll Practice

- Go HTML templates w/ layout + components ✅
- HTMX for real-time interactivity ✅
- Tailwind for styling ✅
- URL routing + query handling in Go ✅
- Basic CRUD over a slice/map (or SQLite later) ✅
- Optional: persistence w/ a local file or sqlite 🔥
- Optional: keyboard shortcuts (`js + htmx`) 🔥

---

### 🔄 Future You Could Add

- ⏱️ Due dates + calendar view
- 🔁 Recurring tasks
- 📦 Export to Markdown / JSON
- 🔒 Auth if you make it multi-user later

---
