package sematextlogger

import (
	"errors"
	"io"
	"log"
	"strings"

	"github.com/labstack/echo"

	"os"

	"github.com/Ephram84/sematext-logger/logwriter"
	"github.com/labstack/echo/middleware"
	echoLog "github.com/labstack/gommon/log"
)

var (
	errSematextWriterIsNil = errors.New("sematextWriter is nil")
)

// Logger is a writer to sematext server.
type Logger struct {
	sematextWriter *logwriter.Writer
	service        string
	*echoLog.Logger
}

// NewLogger returns a logger
func NewLogger(appToken, service string) (*Logger, error) {
	sematext, err := logwriter.Dial("udp", "logsene-receiver-syslog.sematext.com:514", logwriter.LOG_LOCAL0, appToken)
	if err != nil {
		return &Logger{sematextWriter: nil}, err
	}

	newLogger := &Logger{sematextWriter: sematext, service: service}
	newLogger.Logger = echoLog.New(service)

	newLogger.SetOutput(io.MultiWriter(sematext, os.Stdout))

	newLogger.SetHeader(`{"time":"${time_rfc3339}","loglevel":"${level}","service":"${prefix}",` +
		`"file":"${short_file}","line":"${line}"}`)

	// set log level
	if os.Getenv("ENVIRONMENT") == "dev" {
		newLogger.SetLevel(echoLog.DEBUG)
	} else {
		newLogger.SetLevel(echoLog.INFO)
	}

	return newLogger, nil
}

// InitLogger inits a logger through an enviroment variable that contains a url,
// e.g., https://logsene-receiver.sematext.com/fzr64ktn-....
func InitLogger(envVArName, prefix string) (*Logger, error) {
	loggingURL := os.Getenv(envVArName)
	appToken := ""
	if loggingURL != "" {
		log.Println("found", envVArName, ": ", loggingURL)
		appToken = strings.Replace(loggingURL, "https://logsene-receiver.sematext.com/", "", -1)
		appToken = strings.Replace(appToken, "/", "", -1)
	} else {
		panic("Url has not been found")
	}

	return NewLogger(appToken, prefix)
}

func (logger *Logger) EchoMiddlwareLogger() echo.MiddlewareFunc {
	if logger.sematextWriter == nil {
		log.Println("sematextWriter is nil")
		return nil
	}

	return middleware.LoggerWithConfig(middleware.LoggerConfig{Output: io.MultiWriter(os.Stdout, logger.sematextWriter),
		Format: `{"time":"${time_rfc3339_nano}", "request_id":"${id}", "remote_ip":"${remote_ip}", "host":"${host}",` +
			` "method":"${method}", "uri":"${uri}", "status":${status}, "latency":${latency},` +
			` "latency_human":"${latency_human}", "bytes_in":${bytes_in},` +
			` "bytes_out":${bytes_out}, "service":"` + logger.service + `"}` + "\n"})
}

// Err logs a message with severity "err".
func (logger *Logger) Err(msg string) error {
	if logger.sematextWriter == nil {
		log.Println("sematextWriter is nil")
		return errSematextWriterIsNil
	}

	return logger.sematextWriter.Err(msg)
}

// Info logs a message with severity "info".
func (logger *Logger) Info(msg string) error {
	if logger.sematextWriter == nil {
		log.Println("sematextWriter is nil")
		return errSematextWriterIsNil
	}
	return logger.sematextWriter.Info(msg)
}

// Warning logs a message with severity "warning".
func (logger *Logger) Warning(msg string) error {
	if logger.sematextWriter == nil {
		log.Println("sematextWriter is nil")
		return errSematextWriterIsNil
	}
	return logger.sematextWriter.Warning(msg)
}

// Debug logs a message with severity "debug".
func (logger *Logger) Debug(msg string) error {
	if logger.sematextWriter == nil {
		log.Println("sematextWriter is nil")
		return errSematextWriterIsNil
	}
	return logger.sematextWriter.Debug(msg)
}
