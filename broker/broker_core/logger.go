package main

import (
	"log"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {

	InfoLogger = log.Default()
	WarningLogger = log.Default()
	ErrorLogger = log.Default()

	InfoLogger.SetPrefix("INFO: ")
	InfoLogger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	WarningLogger.SetPrefix("INFO: ")
	WarningLogger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	ErrorLogger.SetPrefix("INFO: ")
	ErrorLogger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	//InfoLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	//WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	//ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
