```
fem-htmx-proj/
├── cmd/
│   └── taskvault/
│       └── main.go         # Entry point for your app
├── internal/
│   ├── auth/               # Private application code (auth logic)
│   ├── db/                 # Database access, migrations, models
│   ├── middleware/         # Echo/HTTP middleware
│   ├── service/            # Business logic/services
│   └── utils/              # Utility functions
├── pkg/                    # (optional) Public Go packages
├── web/
│   ├── css/
│   ├── js/
│   ├── images/
│   └── views/              # HTML templates
├── scripts/                # Helper scripts (migrate, seed, etc)
├── migrations/             # SQL migration files
├── test/                   # Additional test data and helpers
├── .env
├── .gitignore
├── go.mod
├── go.sum
├── README.md
└── docs/                   # Markdown docs, diagrams, etc
```

**How to migrate your current code:**

- Move main.go to main.go.
- Move `backend/auth/`, `backend/db/`, `backend/handlers/`, `backend/middleware/`, `backend/models/`, `backend/service/`, `backend/utils/` into `internal/` (e.g., `internal/handlers/`).
- Move static assets (`css/`, `js/`, `images/`) and `views/` into a `web/` directory.
- Place SQL migrations in `migrations/`.
- Keep documentation and planning markdown files in `docs/`.

**Example:**

```go
# Project Structure (golang-standards/project-layout)

- `cmd/taskvault/` — Main application entrypoint
- `internal/` — Private Go packages (auth, db, handlers, etc)
- `web/` — Static assets and templates
- `migrations/` — SQL migration files
- `docs/` — Documentation
```

For more details, see the [project-layout repo](https://github.com/golang-standards/project-layout). This structure will make your project idiomatic, maintainable, and ready for growth.

kdiffin: i think i should go from handlers/services to feature-based packages like and task should have handler.go inside of it

GitHub Copilot: That’s a great idea! A **feature-based package structure** is very popular in modern Go web apps. It improves modularity and makes it easier to scale as your app grows.

**Recommended structure:**

```
fem-htmx-proj/
├── cmd/
│   └── taskvault/
│       └── main.go
├── internal/
│   ├── task/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── model.go
│   │   └── repository.go
│   ├── tag/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── model.go
│   ├── user/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── model.go
│   ├── middleware/
│   └── utils/
├── web/
│   ├── css/
│   ├── js/
│   ├── images/
│   └── views/
├── migrations/
├── scripts/
├── test/
├── go.mod
├── go.sum
└── README.md
```

**Benefits:**

- Each feature (task, tag, user, etc.) is self-contained.
- Handlers, services, models, and repositories for a feature live together.
- Easier to maintain and onboard new contributors.

**Tip:**  
Keep shared code (middleware, utils, etc.) in their own packages under `internal/`.

This approach is widely used in Go, especially for web APIs and larger projects.

kdiffin: what would i put inside of scripts

GitHub Copilot: The `scripts` directory is for helper scripts that automate development, build, deployment, or maintenance tasks. These scripts are usually **not** part of your application’s runtime, but help you and your team work with the project.

**Common examples for `scripts/`:**

- `migrate.sh` — Run database migrations
- `seed.sh` — Seed the database with test data
- `build.sh` — Custom build commands
- `test.sh` — Run tests with custom flags or setup
- `lint.sh` — Run linters or formatters
- `dev.sh` — Start the development server with environment variables
- `generate.sh` — Run code generation tools (e.g., `go generate`)
- `docker-build.sh` — Build Docker images

Scripts can be Bash, Python, or any language your team prefers.  
Keep them simple, well-documented, and cross-platform if possible.

**Example:**

```
scripts/
├── migrate.sh
├── seed.sh
├── build.sh
└── lint.sh
```

Add a short comment at the top of each script to explain its purpose.

kdiffin: Standard Go Project Layout
Translations:

한국어 문서
简体中文
正體中文
简体中文 - ???
Français
日本語
Português
Español
Română
Русский
Türkçe
Italiano
Tiếng Việt
Українська
Indonesian
हिन्दी
فارسی
Беларуская
Overview
This is a basic layout for Go application projects. Note that it's basic in terms of content because it's focusing only on the general layout and not what you have inside. It's also basic because it's very high level and it doesn't go into great details in terms of how you can structure your project even further. For example, it doesn't try to cover the project structure you'd have with something like Clean Architecture.

This is NOT an official standard defined by the core Go dev team. This is a set of common historical and emerging project layout patterns in the Go ecosystem. Some of these patterns are more popular than others. It also has a number of small enhancements along with several supporting directories common to any large enough real world application. Note that the core Go team provides a great set of general guidelines about structuring Go projects and what it means for your project when it's imported and when it's installed. See the Organizing a Go module page in the official Go docs for more details. It includes the internal and cmd directory patterns (described below) and other useful information.

If you are trying to learn Go or if you are building a PoC or a simple project for yourself this project layout is an overkill. Start with something really simple instead (a single main.gofile andgo.mod is more than enough). As your project grows keep in mind that it'll be important to make sure your code is well structured otherwise you'll end up with a messy code with lots of hidden dependencies and global state. When you have more people working on the project you'll need even more structure. That's when it's important to introduce a common way to manage packages/libraries. When you have an open source project or when you know other projects import the code from your project repository that's when it's important to have private (aka internal) packages and code. Clone the repository, keep what you need and delete everything else! Just because it's there it doesn't mean you have to use it all. None of these patterns are used in every single project. Even the vendor pattern is not universal.

With Go 1.14 Go Modules are finally ready for production. Use Go Modules unless you have a specific reason not to use them and if you do then you don’t need to worry about $GOPATH and where you put your project. The basic go.mod file in the repo assumes your project is hosted on GitHub, but it's not a requirement. The module path can be anything though the first module path component should have a dot in its name (the current version of Go doesn't enforce it anymore, but if you are using slightly older versions don't be surprised if your builds fail without it). See Issues 37554 and 32819 if you want to know more about it.

This project layout is intentionally generic and it doesn't try to impose a specific Go package structure.

This is a community effort. Open an issue if you see a new pattern or if you think one of the existing patterns needs to be updated.

If you need help with naming, formatting and style start by running gofmt and staticcheck. The previous standard linter, golint, is now deprecated and not maintained; use of a maintained linter such as staticcheck is recommended. Also make sure to read these Go code style guidelines and recommendations:

https://talks.golang.org/2014/names.slide
https://golang.org/doc/effective_go.html#names
https://blog.golang.org/package-names
https://go.dev/wiki/CodeReviewComments
Style guideline for Go packages (rakyll/JBD)
See Go Project Layout for additional background information.

More about naming and organizing packages as well as other code structure recommendations:

GopherCon EU 2018: Peter Bourgon - Best Practices for Industrial Programming
GopherCon Russia 2018: Ashley McNamara + Brian Ketelsen - Go best practices.
GopherCon 2017: Edward Muller - Go Anti-Patterns
GopherCon 2018: Kat Zien - How Do You Structure Your Go Apps
A Chinese post about Package-Oriented-Design guidelines and Architecture layer

面向包的设计和架构分层
Go Directories

Main applications for this project.

The directory name for each application should match the name of the executable you want to have (e.g., /myapp).

Don't put a lot of code in the application directory. If you think the code can be imported and used in other projects, then it should live in the directory. If the code is not reusable or if you don't want others to reuse it, put that code in the directory. You'll be surprised what others will do, so be explicit about your intentions!

It's common to have a small main function that imports and invokes the code from the and directories and nothing else.

See the directory for examples.

Private application and library code. This is the code you don't want others importing in their applications or libraries. Note that this layout pattern is enforced by the Go compiler itself. See the Go 1.4 release notes for more details. Note that you are not limited to the top level internal directory. You can have more than one internal directory at any level of your project tree.

You can optionally add a bit of extra structure to your internal packages to separate your shared and non-shared internal code. It's not required (especially for smaller projects), but it's nice to have visual clues showing the intended package use. Your actual application code can go in the /app directory (e.g., /app/myapp) and the code shared by those apps in the /pkg directory (e.g., /pkg/myprivlib).

You use internal directories to make packages private. If you put a package inside an internal directory, then other packages can’t import it unless they share a common ancestor. And it’s the only directory named in Go’s documentation and has special compiler treatment.

Library code that's ok to use by external applications (e.g., /mypubliclib). Other projects will import these libraries expecting them to work, so think twice before you put something here :-) Note that the internal directory is a better way to ensure your private packages are not importable because it's enforced by Go. The directory is still a good way to explicitly communicate that the code in that directory is safe for use by others. The I'll take pkg over internal blog post by Travis Jeffery provides a good overview of the pkg and internal directories and when it might make sense to use them.

It's also a way to group Go code in one place when your root directory contains lots of non-Go components and directories making it easier to run various Go tools (as mentioned in these talks: Best Practices for Industrial Programming from GopherCon EU 2018, GopherCon 2018: Kat Zien - How Do You Structure Your Go Apps and GoLab 2018 - Massimiliano Pippi - Project layout patterns in Go).

See the directory if you want to see which popular Go repos use this project layout pattern. This is a common layout pattern, but it's not universally accepted and some in the Go community don't recommend it.

It's ok not to use it if your app project is really small and where an extra level of nesting doesn't add much value (unless you really want to :-)). Think about it when it's getting big enough and your root directory gets pretty busy (especially if you have a lot of non-Go app components).

The pkg directory origins: The old Go source code used to use pkg for its packages and then various Go projects in the community started copying the pattern (see this Brad Fitzpatrick's tweet for more context).

Application dependencies (managed manually or by your favorite dependency management tool like the new built-in Go Modules feature). The go mod vendor command will create the directory for you. Note that you might need to add the -mod=vendor flag to your go build command if you are not using Go 1.14 where it's on by default.

Don't commit your application dependencies if you are building a library.

Note that since 1.13 Go also enabled the module proxy feature (using https://proxy.golang.org as their module proxy server by default). Read more about it here to see if it fits all of your requirements and constraints. If it does, then you won't need the vendor directory at all.

Service Application Directories

OpenAPI/Swagger specs, JSON schema files, protocol definition files.

See the directory for examples.

Web Application Directories

Web application specific components: static web assets, server side templates and SPAs.

Common Application Directories

Configuration file templates or default configs.

Put your confd or consul-template template files here.

System init (systemd, upstart, sysv) and process manager/supervisor (runit, supervisord) configs.

Scripts to perform various build, install, analysis, etc operations.

These scripts keep the root level Makefile small and simple (e.g., https://github.com/hashicorp/terraform/blob/main/Makefile).

See the directory for examples.

Packaging and Continuous Integration.

Put your cloud (AMI), container (Docker), OS (deb, rpm, pkg) package configurations and scripts in the /package directory.

Put your CI (travis, circle, drone) configurations and scripts in the /ci directory. Note that some of the CI tools (e.g., Travis CI) are very picky about the location of their config files. Try putting the config files in the /ci directory linking them to the location where the CI tools expect them (when possible).

IaaS, PaaS, system and container orchestration deployment configurations and templates (docker-compose, kubernetes/helm, terraform). Note that in some repos (especially apps deployed with kubernetes) this directory is called

Additional external test apps and test data. Feel free to structure the directory anyway you want. For bigger projects it makes sense to have a data subdirectory. For example, you can have /data or /testdata if you need Go to ignore what's in that directory. Note that Go will also ignore directories or files that begin with "." or "\_", so you have more flexibility in terms of how you name your test data directory.

See the directory for examples.

Other Directories

Design and user documents (in addition to your godoc generated documentation).

See the directory for examples.

Supporting tools for this project. Note that these tools can import code from the and directories.

See the directory for examples.

Examples for your applications and/or public libraries.

See the directory for examples.

External helper tools, forked code and other 3rd party utilities (e.g., Swagger UI).

Git hooks.

Other assets to go along with your repository (images, logos, etc).

This is the place to put your project's website data if you are not using GitHub pages.

See the directory for examples.

Directories You Shouldn't Have

Some Go projects do have a src folder, but it usually happens when the devs came from the Java world where it's a common pattern. If you can help yourself try not to adopt this Java pattern. You really don't want your Go code or Go projects to look like Java :-)

Don't confuse the project level directory with the directory Go uses for its workspaces as described in How to Write Go Code. The $GOPATH environment variable points to your (current) workspace (by default it points to $HOME/go on non-windows systems). This workspace includes the top level , and directories. Your actual project ends up being a sub-directory under , so if you have the directory in your project the project path will look like this: /path/to/workspace/src/your_project/src/your_code.go. Note that with Go 1.11 it's possible to have your project outside of your GOPATH, but it still doesn't mean it's a good idea to use this layout pattern.

Badges
Go Report Card - It will scan your code with gofmt, go vet, gocyclo, golint, ineffassign, license and misspell. Replace github.com/golang-standards/project-layout with your project reference.

Go Report Card

GoDoc - It will provide online version of your GoDoc generated documentation. Change the link to point to your project.

Go Doc

Pkg.go.dev - Pkg.go.dev is a new destination for Go discovery & docs. You can create a badge using the badge generation tool.

PkgGoDev

Release - It will show the latest release number for your project. Change the github link to point to your project.

Release

Notes
A more opinionated project template with sample/reusable configs, scripts and code is a WIP.
anything else from here thatd be useful? or is this good enough considering my project aint that big, gimme the final project layout in this response pls

GitHub Copilot: Given your project isn’t huge, you don’t need every directory from the full project-layout. You can **skip** things like `/pkg`, `/deployments`, `/third_party`, `/examples`, `/tools`, `/website`, etc.  
You **should** keep the basics: `cmd/`, `internal/` (with feature-based subpackages), `web/`, `migrations/`, `scripts/`, and optionally `test/` and `docs/`.

Here’s a **final, practical layout** for your project:

```
fem-htmx-proj/
├── cmd/
│   └── taskvault/
│       └── main.go
├── internal/
│   ├── task/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── model.go
│   │   └── repository.go
│   ├── tag/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── model.go
│   ├── user/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── model.go
│   ├── middleware/
│   └── utils/
├── web/
│   ├── css/
│   ├── js/
│   ├── images/
│   └── views/
├── migrations/
├── scripts/
│   ├── migrate.sh
│   ├── seed.sh
│   └── lint.sh
├── test/           # (optional, for test helpers or testdata)
├── docs/           # (optional, for markdown docs)
├── .env
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

**Summary:**

- `cmd/` — app entrypoint
- `internal/` — feature-based code (handlers, services, models, etc)
- `web/` — static assets and templates
- `migrations/` — DB migrations
- `scripts/` — helper scripts
- `test/` and `docs/` — optional, add if you need them

This is clean, idiomatic, and easy to grow.  
You can always add more structure later if your project gets bigger!
