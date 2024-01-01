package main

import (
	"log"
	"os"
)

func Print(text string, args ...interface{}) {
	// log to custom file
	LogFile := "/tmp/go-load-balancer-access-log.txt"
	// open log file
	logFile, err := os.OpenFile(LogFile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	// Set log out put and enjoy :)
	log.SetOutput(logFile)
	// optional: log date-time, filename, and line number
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println(text, args)
}
