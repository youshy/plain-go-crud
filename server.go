package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type App struct {
	Db *sql.DB
}

func (a *App) Initialize() {
	PgUsername := os.Getenv("PG_USERNAME")
	PgPassword := os.Getenv("PG_PASSWORD")
	PgDbName := os.Getenv("PG_DB_NAME")
	PgDbHost := os.Getenv("PG_DB_HOST")

	connect := fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=disable", PgDbName, PgUsername, PgPassword, PgDbHost)
	psqlDb, err := sql.Open("postgres", connect)
	if err != nil {
		log.Fatal(err)
	}
	a.Db = psqlDb

	prefix := "/api"

	http.HandleFunc(prefix+"/post", a.Handler)
}

func (a *App) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		postId, _ := r.URL.Query()["post"]
		if postId == "" {
			posts, err := a.GetAllPost()
			if err != nil {
				JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
				return
			}
			JSONResponse(w, http.StatusOK, posts)
			return
		}
		post, err := a.GetSinglePost(postId)
		if err != nil {
			JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
			return
		}

	case "POST":
		w.WriteHeader(http.StatusCreated)

	case "PUT":
		w.WriteHeader(http.StatusAccepted)
		postId, ok := r.URL.Query()["post"]
		if !ok {
			JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "no post specified"})
			return
		}

	case "DELETE":
		w.WriteHeader(http.StatusOK)
		postId, ok := r.URL.Query()["post"]
		if !ok {
			JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "no post specified"})
			return
		}

	default:
		w.WriteHeader(http.StatusNotFound)

	}
}

func (a *App) Run(addr string) {
	log.Printf("Server is listening on %v", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func JSONResponse(w http.ResponseWriter, code int, output interface{}) {
	response, _ := json.Marshal(output)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
