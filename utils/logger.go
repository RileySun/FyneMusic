package utils

import(
	"os"
	"log"
	"path/filepath"
)

var logPath string

func SetLogPath(logDir string) {
	logPath = filepath.Join(logDir, "Log.txt")
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }
    
    log.SetOutput(file)
}