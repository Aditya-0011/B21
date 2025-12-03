package main

import (
	"fmt"
	"net/http"
	"os"
	"proxy-server/server"
	"proxy-server/utils"
)

func main() {
	logFile := os.Getenv("LOG_FILE")
	if logFile == "" {
		logFile = "proxy.log"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}

	log := utils.Logger{LogFile: logFile}

	log.Setup()

	secret := os.Getenv("TOTP_SECRET")
	if secret == "" {
		log.Entry("[ERROR] TOTP secret not configured")
		return
	}

	controller := server.Controller{
		Log:     &log,
		LogFile: logFile,
		Secret:  secret,
	}

	http.HandleFunc("/", controller.GetData)

	http.HandleFunc("/logs", controller.GetLogs)

	log.Entry(fmt.Sprintf("Server running on: %s", port))
	http.ListenAndServe(":"+port, nil)
}
