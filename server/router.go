package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"telegram-parser/db"
	"time"
)

type Server struct {
	router *mux.Router
	Server *http.Server
}

func Init(dbCli db.DB) *Server {
	var r = mux.NewRouter()
	var h = handlerInit(dbCli)

	// run the update handler in database "top_3_hour"
	trackTopMsgsIn3Hours(h)

	r.HandleFunc("/ws", h.UpgradeToWs).Methods("GET")
	r.HandleFunc("/best", h.getBestInPeriod).Methods("GET").Queries("period", "{period}")
	r.HandleFunc("/best/3hour", h.getTopMsgsIn3Hours).Methods("GET") // DEPRECATE
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui/")))
	//r.Use(h.sessionMiddleware)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return &Server{
		router: r,
		Server: srv,
	}
}
