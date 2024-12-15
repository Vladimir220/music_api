// @title Music_API
// @version 0.0.1
// @host localhost:1234
// @schemes http
// @BasePath /
package main

func main() {
	infoLog, debugLog := InitSystem()

	dao, err := createDaoPostgreSQL[Track]()
	if err != nil {
		debugLog.Fatal(err)
	}
	defer dao.Close()

	e := createTrackEnricherDefault()
	s := createService(dao, e, debugLog)
	h := createHandlers(s, infoLog, debugLog)
	r := createRouter(h, infoLog)

	r.initHandlers()
}
