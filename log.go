package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Open or make directory and file for log
func OpenLogFile() *os.File {
	// log file name is current day
	date := time.Now().Format("2006-01-02")
	logFolderPath := "./log"
	logFilePath := fmt.Sprintf("%s/%s.log", logFolderPath, date)

	if _, err := os.Stat(logFolderPath); os.IsNotExist(err) {
		os.MkdirAll(logFolderPath, 0777)
	}

	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		os.Create(logFilePath)
	}

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return logFile
}

// Write log
func WriteLog(logFile *os.File, logData string) {
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println(logData)
}
