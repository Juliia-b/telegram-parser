package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"telegram-parser/db"
	"time"
)

type Router struct {
	router *mux.Router
	server *http.Server
}

func RouterInit(dbCli db.DB) *Router {
	var r = mux.NewRouter()

	h := &handler{dbCli: dbCli}

	r.HandleFunc("/max/{period}", h.GetMaxInPeriod).Methods("GET")
	r.HandleFunc("/min/{period}", h.GetMinInPeriod).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return &Router{
		router: r,
		server: srv,
	}
}

func (r *Router) ListenAndServe() error {
	return r.server.ListenAndServe()
}
