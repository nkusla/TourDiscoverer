package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"
)

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
	DEBUG LogLevel = "DEBUG"
)

type LogEntry struct {
	Timestamp  time.Time `json:"timestamp"`
	Level      LogLevel  `json:"level"`
	Service    string    `json:"service"`
	Message    string    `json:"message"`
	Method     string    `json:"method,omitempty"`
	URL        string    `json:"url,omitempty"`
	StatusCode int       `json:"status_code,omitempty"`
	Duration   int64     `json:"duration_ms,omitempty"`
	UserID     string    `json:"user_id,omitempty"`
	RequestID  string    `json:"request_id,omitempty"`
	Error      string    `json:"error,omitempty"`
	FileName   string    `json:"file,omitempty"`
	LineNumber int       `json:"line,omitempty"`
}

type AppLogger struct {
	serviceName  string
	logstashConn net.Conn
	localLogger  *log.Logger
}

func NewAppLogger(serviceName string) *AppLogger {
	logger := &AppLogger{
		serviceName: serviceName,
		localLogger: log.New(os.Stdout, fmt.Sprintf("[%s] ", serviceName), log.LstdFlags),
	}

	// Čitaj Logstash port iz environment varijable
	logstashPort := os.Getenv("LOGSTASH_PORT")
	if logstashPort == "" {
		logstashPort = "5000" // default fallback
	}

	logstashAddr := fmt.Sprintf("logstash:%s", logstashPort)
	conn, err := net.Dial("tcp", logstashAddr)
	if err != nil {
		logger.localLogger.Printf("Failed to connect to Logstash at %s: %v. Using local logging only.", logstashAddr, err)
	} else {
		logger.logstashConn = conn
		logger.Info("Connected to Logstash successfully at " + logstashAddr)
	}

	return logger
}

func (l *AppLogger) Close() {
	if l.logstashConn != nil {
		l.logstashConn.Close()
	}
}

func (l *AppLogger) logEntry(entry LogEntry) {
	entry.Timestamp = time.Now()
	entry.Service = l.serviceName

	// Dodaj informacije o fajlu i liniji
	if _, file, line, ok := runtime.Caller(2); ok {
		entry.FileName = file
		entry.LineNumber = line
	}

	// Strukturirani JSON log u stdout
	jsonLog, _ := json.Marshal(entry)
	l.localLogger.Printf("%s [%s] %s", entry.Timestamp.Format(time.RFC3339), entry.Level, entry.Message)

	// Pošalji u Logstash ako je konekcija aktivna
	if l.logstashConn != nil {
		jsonLog = append(jsonLog, '\n')
		if _, err := l.logstashConn.Write(jsonLog); err != nil {
			l.localLogger.Printf("Failed to send log to Logstash: %v", err)
		}
	}
}

func (l *AppLogger) Info(message string) {
	l.logEntry(LogEntry{
		Level:   INFO,
		Message: message,
	})
}

func (l *AppLogger) InfoWithFields(message string, fields map[string]interface{}) {
	entry := LogEntry{
		Level:   INFO,
		Message: message,
	}

	if userID, ok := fields["user_id"].(string); ok {
		entry.UserID = userID
	}
	if requestID, ok := fields["request_id"].(string); ok {
		entry.RequestID = requestID
	}
	if method, ok := fields["method"].(string); ok {
		entry.Method = method
	}
	if url, ok := fields["url"].(string); ok {
		entry.URL = url
	}
	if statusCode, ok := fields["status_code"].(int); ok {
		entry.StatusCode = statusCode
	}
	if duration, ok := fields["duration_ms"].(int64); ok {
		entry.Duration = duration
	}

	l.logEntry(entry)
}

func (l *AppLogger) Warn(message string) {
	l.logEntry(LogEntry{
		Level:   WARN,
		Message: message,
	})
}

func (l *AppLogger) Error(message string, err error) {
	entry := LogEntry{
		Level:   ERROR,
		Message: message,
	}

	if err != nil {
		entry.Error = err.Error()
	}

	l.logEntry(entry)
}

func (l *AppLogger) Debug(message string) {
	l.logEntry(LogEntry{
		Level:   DEBUG,
		Message: message,
	})
}

// HTTP Middleware
func (l *AppLogger) HTTPMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, statusCode: 200}

		// Log incoming request
		l.InfoWithFields("Incoming HTTP request", map[string]interface{}{
			"method":     r.Method,
			"url":        r.URL.String(),
			"user_agent": r.Header.Get("User-Agent"),
			"remote_ip":  r.RemoteAddr,
		})

		next(rw, r)

		duration := time.Since(start).Milliseconds()

		// Log completed request
		l.InfoWithFields("HTTP request completed", map[string]interface{}{
			"method":      r.Method,
			"url":         r.URL.String(),
			"status_code": rw.statusCode,
			"duration_ms": duration,
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
