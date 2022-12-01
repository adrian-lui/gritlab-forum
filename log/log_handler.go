package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

const logFile = "log/main.log" // Set log file location

func init() {
	// Rotatet log files
	if _, err := os.Stat(logFile); err == nil {
		// Rotate last log
		e := os.Rename(logFile, "log/main01.log")
		if e != nil {
			log.Fatal(e)
		}
	}

}

// Write To Log
func WTL(logString string, toStdOut bool) {
	if toStdOut {
		fmt.Println(logString)
	}

	// Open log file
	fh, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Could not open log-file to append string, " + err.Error())
		return
	}

	defer fh.Close()

	logString = time.Now().String()[0:19] + ": " + logString + "\n"

	if _, err := fh.WriteString(logString); err != nil {
		fmt.Println(err.Error())
	}
}
