package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"telegram-parser/db"
)

type Server struct {
	router *mux.Router
	Server *http.Server
}

func ServerInit(dbCli db.DB) *Server {
	var r = mux.NewRouter()

	h := &handler{dbCli: dbCli}

	r.HandleFunc("/best", h.GetBestInPeriod).Methods("GET").Queries("period", "{period}")

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		//WriteTimeout: 15 * time.Second,
		//ReadTimeout:  15 * time.Second,
	}

	return &Server{
		router: r,
		Server: srv,
	}
}
