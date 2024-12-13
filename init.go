package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitSystem() (infoLog *log.Logger, debugLog *log.Logger) {

	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	infoLogFile, err := os.OpenFile(os.Getenv("log_info_file_path"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err.Error())
	}
	debugLogLogFile, err := os.OpenFile(os.Getenv("log_debug_file_path"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err.Error())
	}

	infoLog = log.New(infoLogFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix)
	debugLog = log.New(debugLogLogFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix)
	return
}
