# Swagtask

this project started off as a fork of a lesson given by ThePrimeagen on Frontend Masters which then went to be this project.

Everything here was hand-rolled (even the pub-sub) and follows **_HATEOAS_**.

no frameworks, no build step, just Go stdlib and htmx
pub-sub, websockets, and all the glue code is mine
no dependencies, no npm, no yarn, no package.json (If I havent made it clear yet, I'm tired of JS build step/dependency hell.)
This app is a heavy WIP.

I made this website to test the limits of a no build-step, low-js, zero dependency system.

I stumbled upon the Go programming language, the Go standard library and HTMX, Which I fell in love with after developing this app. I am quite impressed with how far I was able to take this considering it was my first experience with all of these technologies.

Because of this project being so no-dependency focused I actually managed to learn many topics I delegated to SaaS and libraries previously.

Most notably being http, networking and how the browser actually works -> this was delegated to nextjs previously and I didnt understand the network tab in my browser, with htmx the network tab becomes a first class citizen in development which greatly enhanced my understanding of REST principles and http/s.

You become enlightened to browser native caching techniques and actual performance optimizations instead of pointless v-dom rerender juggling with state and memoization attemps, HTMX is lean and mean, low-level-esque and imperative, and a joy for backend developers.

Really enjoyed that freedom this project.

Getting exposed to an event-driven system leveled me up and now I really know when I need smth reactive (lets say a chat app thats dynamic in terms of client side rendering) or when i need smth event driven (most things). This realization has piqued my interest in templ and alpine, which I believe could be a good replacement to the missing part of my stack, being islands of reactivity and componentization

My only drawbacks here have been the lack of UI componentization and typesafety (in go's http/template). In my future projects I'll use templ for sure. Template based dev ain't for me, more of a components guy.

And for the realtime thing I think a reactive way of doings things would be better if the system got more complex. (this is just a hunch though, after I do the alpine stuff I'll have an opinion on this for real)

---

the things that's left to implement
[tasks](tasks.md)

---

## gippity made the bottom part

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
