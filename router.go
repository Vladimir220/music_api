package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type router struct {
	h       handlers
	infoLog *log.Logger
}

func (r router) initHandlers() {
	router := mux.NewRouter()
	host := os.Getenv("NET_HOST")

	router.HandleFunc("/all/pages", r.h.getAll).Methods(http.MethodGet)
	r.infoLog.Printf("Зарегистрирован метод Get на маршрут %s/all/pages/\n", host)

	router.HandleFunc("/track/lyrics/couplets", r.h.getTrackLyrics).Methods(http.MethodGet)
	r.infoLog.Printf("Зарегистрирован метод Get на маршрут %s/track/lyrics/couplets\n", host)

	router.HandleFunc("/track", r.h.deleteTrack).Methods(http.MethodDelete)
	r.infoLog.Printf("Зарегистрирован метод Delete на маршрут %s/track\n", host)

	router.HandleFunc("/track", r.h.updateTrack).Methods(http.MethodPatch)
	r.infoLog.Printf("Зарегистрирован метод Patch на маршрут %s/track\n", host)

	router.HandleFunc("/track", r.h.createTrack).Methods(http.MethodPost)
	r.infoLog.Printf("Зарегистрирован метод Post на маршрут %s/track\n", host)

	router.HandleFunc("/info", r.h.getInfo).Methods(http.MethodGet)
	r.infoLog.Printf("Зарегистрирован метод Get на маршрут %s/info\n", host)

	http.Handle("/", router)
	http.ListenAndServe(host, nil)
}

func createRouter(h handlers, infoLog *log.Logger) router {
	return router{h: h, infoLog: infoLog}
}
