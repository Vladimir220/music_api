// @title Music_API
// @version 0.0.1
// @host localhost:1234
// @schemes http
// @BasePath /
package main

func main() {
	infoLog, debugLog := InitSystem()
	daoDB, err := createDaoPostgreSQL[Track]()
	if err != nil {
		debugLog.Fatal(err)
	}
	defer daoDB.Close()

	enrchDefault := createTrackEnricherDefault()
	enrchLyricsCom := createDaoLyricsCom()
	enrchLastFm := createDaoLastFm()

	enrchLastFm.SetNext(enrchLyricsCom)
	enrchLyricsCom.SetNext(enrchDefault)

	h := createHandlers(daoDB, enrchLastFm, createService, infoLog, debugLog)
	r := createRouter(h, infoLog)

	r.initHandlers()
}
