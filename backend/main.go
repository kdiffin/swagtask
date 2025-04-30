package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	db "swagtask/db/generated"
	"swagtask/handlers"
	"swagtask/middleware"
	"swagtask/models"
	"swagtask/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}




func main() {
	// DB INIT START
	// pgx pool starts a pool thats concurrency safe
	dbpool, errDbConn := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	// connects to db via url
	if errDbConn != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", errDbConn)
		os.Exit(1)
	}
	defer dbpool.Close()
	queries := db.New(dbpool)
	// DB INIT END
	
	
	log.SetFlags(log.LstdFlags)
	templates := models.NewTemplate()
	mux := http.NewServeMux()
	// e.Renderer = newTemplate() implement
	// .Use(middleware.Logger()) implement
	// .Static("/images", "../images")
	// e.Static("/css", "../css")
	// e.Static("/js", "../js")
	// router.Tasks(mux, queries, templates)

	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("../images"))))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../css"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("../js"))))

	// tasks
	mux.HandleFunc("GET /tasks/{$}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerGetTasks(w, r, queries, templates)
	})
	mux.HandleFunc("POST /tasks/{$}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerCreateTask(w, r, queries, templates)
	})
	mux.HandleFunc("POST /tasks/{id}/toggle-complete/{$}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, errConv := strconv.Atoi(idStr) 
		if errConv != nil {
			utils.LogError("couldnt convert to str", errConv)
			return
		}

		handlers.HandlerTaskToggleComplete(w, r, queries, templates, int32(id))
	})
	mux.HandleFunc("POST /tasks/{id}/tags/{$}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, errConv1 := strconv.Atoi(idStr) 
		tagIdStr := r.FormValue("tag_id")
		tagId, errConv := strconv.Atoi(tagIdStr) 
		if errConv != nil || errConv1 != nil  {
			utils.LogError("couldnt convert to str", errConv)
			utils.LogError("couldnt convert to str", errConv1)
			return
		}

		handlers.HandlerAddTagToTask(w,r,queries, templates, int32(id), int32(tagId))
	})
	mux.HandleFunc("DELETE /tasks/{id}/{$}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, errConv := strconv.Atoi(idStr) 
		if errConv != nil {
			utils.LogError("couldnt convert to str", errConv)
			return
		}

		handlers.HandlerDeleteTask(w,r,queries, templates, int32(id))
	})
	mux.HandleFunc("DELETE /tasks/{id}/tags/{$}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, errConv := strconv.Atoi(idStr) 
		tagIdStr := r.FormValue("tag_id")
		tagId, errConv1 := strconv.Atoi(tagIdStr) 
		if errConv != nil || errConv1 != nil {
			utils.LogError("couldnt convert to str", errConv)
			utils.LogError("couldnt convert to str 2nd", errConv1)
			return
		}


		handlers.HandlerRemoveTagFromTask(w,r,queries, templates, int32(id), int32(tagId))
	})
	mux.HandleFunc("PUT /tasks/{id}/", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("task_name")
		idea := r.FormValue("task_idea")
		idStr := r.PathValue("id")
		id, errConv := strconv.Atoi(idStr) 
		if errConv != nil  {
			utils.LogError("couldnt convert to str", errConv)
			return
		}
		handlers.HandlerUpdateTask(w,r,queries, templates, int32(id), idea, name)
	})

	// tags
	mux.HandleFunc("POST /tags/{$}", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("tag_name")
		source := r.FormValue("source")
		handlers.HandlerCreateTag(w, r, queries, templates, name, source )
	})
	mux.HandleFunc("GET /tags/{$}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerGetTags(w,r,queries,templates)
	})
	mux.HandleFunc("PUT /tags/{id}/{$}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, errConv := strconv.Atoi(idStr) 
		if errConv != nil  {
			utils.LogError("couldnt convert to str", errConv)
			return
		}
		name := r.FormValue("tag_name")

		handlers.HandlerUpdateTag(w,r,queries,templates,  name, int32(id))	
	})
	mux.HandleFunc("DELETE /tags/{id}/{$}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, errConv := strconv.Atoi(idStr) 
		if errConv != nil  {
			utils.LogError("couldnt convert to str", errConv)
			return
		}

		handlers.HandlerDeleteTag(w,r,queries,templates, int32(id))
	})
	mux.HandleFunc("POST /tags/{id}/tasks/{$}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, errConv := strconv.Atoi(idStr) 
		taskIdStr := r.FormValue("task_id")
		taskId, errConv2 := strconv.Atoi(taskIdStr)
		if errConv != nil || errConv2 != nil  {
			utils.LogError("couldnt convert to str", errConv)
			utils.LogError("couldnt convert to str2", errConv2)
			return
		}

		handlers.HandlerAddTaskToTag(w, r, queries, templates, int32(taskId), int32(id))
	})
	mux.HandleFunc("DELETE /tags/{id}/tasks/{$}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, errConv := strconv.Atoi(idStr) 
		taskIdStr := r.FormValue("task_id")	
		taskId, errConv2 := strconv.Atoi(taskIdStr)
		if errConv != nil || errConv2 != nil  {
			utils.LogError("couldnt convert to str", errConv)
			utils.LogError("couldnt convert to str2", errConv2)
			return
		}

		handlers.HandlerRemoveTaskFromTag(w, r, queries, templates, int32(taskId),int32(id))
	})
	

	server := http.Server{
		Addr: ":42069",
		Handler: middleware.Logging(mux),
	}
	fmt.Println("running server")
	log.Fatal(server.ListenAndServe())
	fmt.Println("not running server")

}
