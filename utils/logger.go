package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type Logger struct {
	LogFile string
}

var LogMutex sync.Mutex

func (l *Logger) Setup() {
	logFile, err := os.OpenFile(l.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func (l *Logger) Entry(entry string) {
	LogMutex.Lock()
	defer LogMutex.Unlock()

	log.Println(entry)
}
