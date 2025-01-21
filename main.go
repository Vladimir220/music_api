// @title Music_API
// @version 0.0.1
// @host localhost:1234
// @schemes http
// @BasePath /
package main

import (
	"log"
	"music_api/api"
	"music_api/dao/caching"
	"music_api/dao/db"
	enrch "music_api/dao/enrichment"
	"music_api/models"
)

func main() {
	infoLog, debugLog := api.InitSystem()

	daoDB, err := db.CreateDaoPostgreSQL[models.Track]()
	if err != nil {
		log.Println(err.Error())
		debugLog.Fatal(err)
	}
	defer daoDB.Close()

	enrchDefault := enrch.CreateTrackEnricherDefault()
	enrchLyricsCom := enrch.CreateDaoLyricsCom()
	enrchLastFm := enrch.CreateDaoLastFm()
	enrchServer := enrch.CreateDaoEnrchServer()

	caching, _ := caching.CreateDaoRedis()

	enrchServer.SetNext(enrchLastFm)
	enrchLastFm.SetNext(enrchLyricsCom)
	enrchLyricsCom.SetNext(enrchDefault)

	h := api.CreateHandlers(daoDB, enrchServer, caching, api.CreateService, infoLog, debugLog)
	r := api.CreateRouter(h, infoLog)

	r.InitHandlers()
}
