package router

// import (
// 	"context"
// 	"net/http"
// 	"strconv"
// 	db "swagtask/db/generated"
// 	"swagtask/models"
// )

// func Test(mux *http.ServeMux, queries *db.Queries, renderer *models.Template)  {
// 	mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
// 		allTags, errGetAllTags := queries.GetAllTagsDesc(context.Background())

// 		if errGetAllTags != nil {
// 			http.Error(w, "nolar apar meni", http.StatusInternalServerError)
// 			return
// 		}

// 		var response string
// 		for _,tag := range allTags{
// 			response +=  tag.Name +  strconv.Itoa(int(tag.ID)) + "\n"
// 		}

// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte(response))
// 		return
// 	})
// }