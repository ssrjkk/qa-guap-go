package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var (
	loggerLevel = INFO
	logger      = log.New(os.Stdout, "", 0)
)

func SetLevel(level Level) {
	loggerLevel = level
}

func Debug(format string, v ...interface{}) {
	if loggerLevel <= DEBUG {
		logger.Printf("[DEBUG] "+format, v...)
	}
}

func Info(format string, v ...interface{}) {
	if loggerLevel <= INFO {
		logger.Printf("[INFO] "+format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	if loggerLevel <= WARN {
		logger.Printf("[WARN] "+format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if loggerLevel <= ERROR {
		logger.Printf("[ERROR] "+format, v...)
	}
}

type RequestLog struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
	Time    time.Time
}

type ResponseLog struct {
	StatusCode int
	Headers    map[string]string
	Body       string
	Duration   time.Duration
}

func LogRequest(req *RequestLog) {
	Info("REQUEST: %s %s", req.Method, req.URL)
	if len(req.Headers) > 0 {
		Debug("Request Headers: %v", req.Headers)
	}
	if req.Body != "" {
		Debug("Request Body: %s", req.Body)
	}
}

func LogResponse(resp *ResponseLog, url string) {
	Info("RESPONSE [%s]: %d (%v)", url, resp.StatusCode, resp.Duration)
	if resp.StatusCode >= 400 {
		Error("Response Body: %s", resp.Body)
	}
}
