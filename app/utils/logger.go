package utils

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
)

// InitLogger sets up file + console logging.
func InitLogger() {
	logPath := os.Getenv("LOG_FILE")
	if logPath == "" {
		logPath = "storage/logs/app.log"
	}

	if err := os.MkdirAll("storage/logs", 0755); err != nil {
		log.Println("[LianGo] Warning: could not create log directory:", err)
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("[LianGo] Warning: could not open log file, using stdout only")
		file = os.Stdout
	}

	multi := io.MultiWriter(file, os.Stdout)

	flags := log.Ldate | log.Ltime | log.Lshortfile

	InfoLogger = log.New(multi, "[INFO]  ", flags)
	WarnLogger = log.New(multi, "[WARN]  ", flags)
	ErrorLogger = log.New(multi, "[ERROR] ", flags)

	InfoLogger.Println("Logger initialized")
}

// GinLogger returns a Gin middleware for structured HTTP request logging.
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		ip := c.ClientIP()

		if status >= 500 {
			ErrorLogger.Printf("%d | %s | %s | %s | %v", status, method, path, ip, latency)
		} else if status >= 400 {
			WarnLogger.Printf("%d | %s | %s | %s | %v", status, method, path, ip, latency)
		} else {
			InfoLogger.Printf("%d | %s | %s | %s | %v", status, method, path, ip, latency)
		}
	}
}
