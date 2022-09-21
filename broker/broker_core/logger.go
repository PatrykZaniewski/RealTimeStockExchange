package main

import (
	"bytes"
	"log"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	var buff bytes.Buffer

	//InfoLogger = log.Default()
	//InfoLogger.SetPrefix("INFO: ")
	//
	InfoLogger = log.New(&buff, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)

	//WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	//ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
