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

	daoEnrch := createTrackEnricherDefault()
	h := createHandlers(daoDB, daoEnrch, createService, infoLog, debugLog)
	//s := createService(dao, e, debugLog)
	r := createRouter(h, infoLog)

	r.initHandlers()
}
