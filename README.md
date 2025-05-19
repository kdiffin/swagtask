chat gippity made the readme
im doing this to learn htmx and golang i hate react now im a htmx-grug-pilled dev, f the virtual dom
chat gippity made most of my frontends styles btw

---

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

- [x] read the other md files rq
- [x] refactor v2

  - [x] add net/http and remove echo
  - [x] add decent logging on errors
  - [x] move handler logic outside of the router
  - [x] add sqlc
  - [x] complete the rest of the app, get it back to how it was

- [x] implement actual database migrations
- [x] add the actual http/1.1 response and request in the logging middleware so i know http strings are being sent and recieved over http
- [x] host with https on local network

- [ ] add barebones auth using cookies then jwt

  - [x] finish the migration business
  - [x] implement the readme on how to add session auth
  - [x] add sql to the sessions instead of storing them in memory
  - [x] make an indicator showing if the user is logged in or not
  - [x] finish fetching the tags and shit
  - [x] change the hx-redirects to just forms with actions
  - [ ] block csrf, add session expirations, add jwt

- [ ] refactor v3 (vaults)

  - **notes**
  - note to self: every tag action based on its vault should be validated with the user id

  - **large scale objectives**
  - ![alt text](./structure.jpeg)

  - **objectives**
  - [x] replace auto incrementing ids with uuids (no random id guess attack vector)
  - [x] make it work with the uuids and vault ids now

    - [x] change all of the queries to use vault ids
    - [x] add default vault id to users, add vault type
    - [x] make the sign up create a default vault

  - [x] add author info to tags and tasks, fix the n+1s with joins, fix updatednow

    - [x] tags
    - [x] tasks
    - [x] fix the n+1 problem with author info

  - [x] fix the workflow
    - [x] make all queries work with roles and vaults
    - [x] change constraints of uniqueness to include vaults
  - [ ] add redirect to tasks/vault filter on tag id click
  - [ ] fix that one dumb down migration at 8
  - [x] [golang standards project layout](https://github.com/golang-standards/project-layout) conform to this structure [example](./code-structure.md)
  - [x] ALWAYS Update the updated_now fields so its actually true
  - [ ] order everything by created_at
  - [ ] add a vault middleware which gives the vault depending on if the route is the shared one or the default one, cuz rn its only default
  - [ ] add decent error logging

- [ ] add live editing via websockets between 2 users
- [ ] add notification system via sse
- [ ] [add docs to api](https://www.boot.dev/lessons/109e29ef-cdfd-4d5b-ad47-6a609a638896)
- [ ] add some tests bruh

### long term

- [ ] add openapi specs

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
