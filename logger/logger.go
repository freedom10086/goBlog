package logger

import (
	"log"
	"os"
)

var fileLogger *log.Logger

func init() {
	log.Println("start init log")
	var filename = "logger/log.log"

	var file *os.File
	var err error

	file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("fail to create log file: %s \n", filename)
	}

	//defer file.Close()
	fileLogger = log.New(file, "", log.LstdFlags)
}

func D(format string, v ...interface{}) {
	log.Printf(format, v)
	fileLogger.SetPrefix("[D] ")
	fileLogger.Printf(format, v)
}

func I(format string, v ...interface{}) {
	log.Printf(format, v)
	fileLogger.SetPrefix("[I] ")
	fileLogger.Printf(format, v)
}

func E(format string, v ...interface{}) {
	log.Printf(format, v)
	fileLogger.SetPrefix("[E] ")
	fileLogger.Printf(format, v)
}
