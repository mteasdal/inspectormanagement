package auditing

import (
	"bytes"
	"log"
	"os"
)

// Log logs information using the standard Go logger package
func Log(logMessage string) {
	buf := bytes.Buffer{}
	logger := log.New(&buf, "AWSinspector", log.Lshortfile|log.Ldate)
	logFile, err := getLogFile()

	if err != nil {
		log.Fatalf(err.Error())
	}

	logger.SetOutput(logFile)

	logger.Println(logMessage)
	logger.SetPrefix("new logger:")
	closeLogFile(logFile)
}

// getLogFile returns a pointer to a file for logging
func getLogFile() (*os.File, error) {
	file, err := os.OpenFile("inspector.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// closeLogFile closes the memory pointer to the log file.
func closeLogFile(file *os.File) error {
	err := file.Close()
	if err != nil {
		return err
	}
	return nil
}
