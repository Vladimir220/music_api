// @title Music_API
// @version 0.0.1
// @host localhost:1234
// @schemes http
// @BasePath /
package main

import "os"

func main() {
	infoLog, debugLog := InitSystem()
	daoDB, err := createDaoPostgreSQL[Track]()
	if err != nil {
		debugLog.Fatal(err)
	}
	defer daoDB.Close()

	// для проверки подключения токена
	// в гитхаб для безопасности я не буду выкладывать
	// файл окружения с токенами
	// но возможно он будет в файле на Яндекс диске
	token := os.Getenv("TOKEN_LASTFM")
	var daoEnrch DaoEnrichment[Track]
	if token == "" {
		daoEnrch = createDaoLastFm()
	} else {
		daoEnrch = createTrackEnricherDefault()
	}

	h := createHandlers(daoDB, daoEnrch, createService, infoLog, debugLog)
	r := createRouter(h, infoLog)

	r.initHandlers()
}
