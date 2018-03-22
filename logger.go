package sematextlogger

import (
	"io"
	"log"

	"os"

	"github.com/Ephram84/sematext-logger/logwriter"
	echoLog "github.com/labstack/gommon/log"
)

// Logger is a writer to sematext server.
type Logger struct {
	sematextWriter *logwriter.Writer
}

// NewLogger returns a logger
func NewLogger(appToken, prefix string) (*Logger, error) {
	sematext, err := logwriter.Dial("udp", "logsene-receiver-syslog.sematext.com:514", logwriter.LOG_LOCAL0, appToken)
	if err != nil {
		return &Logger{sematextWriter: nil}, err
	}
	echoLog.SetOutput(io.MultiWriter(sematext, os.Stdout))
	echoLog.SetHeader(`{"time":"${time_rfc3339_nano}","request_id":"${id}","remote_ip":"${remote_ip}","host":"${host}",` +
		`"method":"${method}","uri":"${uri}","status":${status}, "latency":${latency},` +
		`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
		`"bytes_out":${bytes_out}, "service":"` + prefix + `"}`)

	// set log level
	if os.Getenv("ENVIRONMENT") == "dev" {
		echoLog.SetLevel(echoLog.DEBUG)
	} else {
		echoLog.SetLevel(echoLog.INFO)
	}

	// set prefix
	echoLog.SetPrefix(prefix)

	return &Logger{sematextWriter: sematext}, nil
}

// Err logs a message with severity "err".
func (logger *Logger) Err(msg string) {
	if logger.sematextWriter == nil {
		log.Println("sematextWriter is nil")
		return
	}
	echoLog.Error(msg)
}

// Info logs a message with severity "info".
func (logger *Logger) Info(msg string) {
	if logger.sematextWriter == nil {
		log.Println("sematextWriter is nil")
		return
	}
	echoLog.Info(msg)
}

// Warning logs a message with severity "warning".
func (logger *Logger) Warning(msg string) {
	if logger.sematextWriter == nil {
		log.Println("sematextWriter is nil")
		return
	}
	echoLog.Warn(msg)
}

// Debug logs a message with severity "debug".
func (logger *Logger) Debug(msg string) {
	if logger.sematextWriter == nil {
		log.Println("sematextWriter is nil")
		return
	}
	echoLog.Warn(msg)
}
