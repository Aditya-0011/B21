package main

import (
	"fmt"
	"net/http"
	"proxy-server/server"
	"proxy-server/utils"
)

const (
	LOG_FILE = "proxy.log"
	PORT     = "6000"
)

func main() {
	log := utils.Logger{LogFile: LOG_FILE}

	log.Setup()

	controller := server.Controller{
		Log:     &log,
		LogFile: LOG_FILE,
	}

	http.HandleFunc("/", controller.GetData)

	http.HandleFunc("/logs", controller.GetLogs)

	log.Entry(fmt.Sprintf("Server running on: %s", PORT))
	http.ListenAndServe(":"+PORT, nil)
}
