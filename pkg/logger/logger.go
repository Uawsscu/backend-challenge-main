package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Logger struct {
	level string
}

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

var defaultLogger *Logger

func Init(level string) {
	defaultLogger = &Logger{level: level}
}

func Info(message string, data ...map[string]interface{}) {
	if defaultLogger == nil {
		defaultLogger = &Logger{level: "info"}
	}
	defaultLogger.log("INFO", message, data...)
}

func Error(message string, data ...map[string]interface{}) {
	if defaultLogger == nil {
		defaultLogger = &Logger{level: "info"}
	}
	defaultLogger.log("ERROR", message, data...)
}

func Debug(message string, data ...map[string]interface{}) {
	if defaultLogger == nil {
		defaultLogger = &Logger{level: "info"}
	}
	if defaultLogger.level == "debug" {
		defaultLogger.log("DEBUG", message, data...)
	}
}

func (l *Logger) log(level, message string, data ...map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   message,
	}

	if len(data) > 0 {
		entry.Data = data[0]
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Failed to marshal log entry: %v", err)
		return
	}

	log.Println(string(jsonData))
}

func SetOutput(file *os.File) {
	log.SetOutput(file)
}
