package sematextlogger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

// Logger is a writer to sematext server.
type Logger struct {
	AppToken string
	Type     string
	Host     string
	URL      string
}

// NewLogger returns a logger
func NewLogger(appToken string) *Logger {
	host, _ := os.Hostname()
	typ := "syslog"
	url := "http://logsene-receiver.sematext.com:80"
	return &Logger{AppToken: appToken, Type: typ, Host: host, URL: url}
}

func (logger *Logger) WithType(t string) *Logger {
	logger.Type = t
	return logger
}

func (logger *Logger) WithURL(url string) *Logger {
	logger.URL = url
	return logger
}

// Err logs a message with severity "err".
func (logger *Logger) Err(msg string, additional ...string) (bool, error) {
	return logger.buildMessage("err", msg, additional)
}

// Info logs a message with severity "info".
func (logger *Logger) Info(msg string, additional ...string) (bool, error) {
	return logger.buildMessage("info", msg, additional)
}

// Emerg logs a message with severity "emerg".
func (logger *Logger) Emerg(msg string, additional ...string) (bool, error) {
	return logger.buildMessage("emerg", msg, additional)
}

// Crit logs a message with severity "crit".
func (logger *Logger) Crit(msg string, additional ...string) (bool, error) {
	return logger.buildMessage("crit", msg, additional)
}

// Warning logs a message with severity "warning".
func (logger *Logger) Warning(msg string, additional ...string) (bool, error) {
	return logger.buildMessage("warning", msg, additional)
}

// Notice logs a message with severity "notice".
func (logger *Logger) Notice(msg string, additional ...string) (bool, error) {
	return logger.buildMessage("notice", msg, additional)
}

// Debug logs a message with severity "debug".
func (logger *Logger) Debug(msg string, additional ...string) (bool, error) {
	return logger.buildMessage("debug", msg, additional)
}

func (logger *Logger) buildMessage(severity, msg string, additional []string) (bool, error) {
	var message map[string]interface{}
	message = make(map[string]interface{})

	message["Severity"] = severity
	message["Message"] = msg
	message["Host"] = logger.Host

	for _, keyvalue := range additional {
		parts := strings.Split(keyvalue, ":")
		switch len(parts) {
		case 1:
			message["?"] = parts[0]
		case 2:
			message[parts[0]] = parts[1]
		default:
			value := parts[1]
			for _, part := range parts[2:] {
				value += ":" + part
			}
			message[parts[0]] = value
		}
		// if len(parts) == 2 {
		// 	message[parts[0]] = parts[1]
		// }
	}

	jmsg, err := json.Marshal(message)
	if err != nil {
		return false, err
	}
	return logger.sendMessage(jmsg)
}

func (logger *Logger) sendMessage(msg []byte) (bool, error) {
	req, err := http.NewRequest("POST", logger.URL+"/"+logger.AppToken+"/"+logger.Type, bytes.NewBuffer(msg))
	if err != nil {
		return false, err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return false, err
	}

	return true, err
}
