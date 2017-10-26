package sematextlogger

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Logger is a writer to sematext server.
type Logger struct {
	AppToken string
	Type     string
	Host     string
	URL      string
}

// NewLogger returns a logger.
func NewLogger(appToken string) *Logger {
	host, _ := os.Hostname()
	typ := "syslog"
	url := "http://logsene-receiver.sematext.com:80"
	return &Logger{AppToken: appToken, Type: typ, Host: host, URL: url}
}

// InitLogger inits a logger throught an enviroment variable that contains a url,
// e.g., https://logsene-receiver.sematext.com/fzr64ktn-....
func InitLogger(envVArName string) *Logger {
	loggingURL := os.Getenv(envVArName)
	appToken := ""
	if loggingURL != "" {
		log.Println("found", envVArName, ": ", loggingURL)
		appToken = strings.Replace(loggingURL, "https://logsene-receiver.sematext.com/", "", -1)
		appToken = strings.Replace(appToken, "/", "", -1)
	} else {
		panic("Url has not been found")
	}

	return NewLogger(appToken)
}

// WithType sets a new type.
func (logger *Logger) WithType(t string) *Logger {
	logger.Type = t
	return logger
}

// WithType sets a new URL.
func (logger *Logger) WithURL(url string) *Logger {
	logger.URL = url
	return logger
}

// Err logs a message with severity "err".
func (logger *Logger) Err(msg string, additional ...string) error {
	return logger.buildMessage("err", msg, additional)
}

// Info logs a message with severity "info".
func (logger *Logger) Info(msg string, additional ...string) error {
	return logger.buildMessage("info", msg, additional)
}

// Emerg logs a message with severity "emerg".
func (logger *Logger) Emerg(msg string, additional ...string) error {
	return logger.buildMessage("emerg", msg, additional)
}

// Crit logs a message with severity "crit".
func (logger *Logger) Crit(msg string, additional ...string) error {
	return logger.buildMessage("crit", msg, additional)
}

// Warning logs a message with severity "warning".
func (logger *Logger) Warning(msg string, additional ...string) error {
	return logger.buildMessage("warning", msg, additional)
}

// Notice logs a message with severity "notice".
func (logger *Logger) Notice(msg string, additional ...string) error {
	return logger.buildMessage("notice", msg, additional)
}

// Debug logs a message with severity "debug".
func (logger *Logger) Debug(msg string, additional ...string) error {
	return logger.buildMessage("debug", msg, additional)
}

func (logger *Logger) buildMessage(severity, msg string, additional []string) error {
	var message map[string]interface{}
	message = make(map[string]interface{})

	unknown := make([]string, 0)

	message["Severity"] = severity
	message["Message"] = msg
	message["Host"] = logger.Host

	for _, keyvalue := range additional {
		if keyvalue == "" {
			continue
		}
		parts := strings.Split(keyvalue, ":")
		switch len(parts) {
		case 0, 1:
			unknown = append(unknown, parts[0])
		case 2:
			if parts[0] == "?" {
				unknown = append(unknown, parts[1])
			} else {
				message[parts[0]] = parts[1]
			}
		default:
			r, _ := regexp.Compile("\\d{4}-\\d{2}-\\d{2}(T| )\\d{2}")
			if r.MatchString(parts[0]) {
				unknown = append(unknown, keyvalue)
			} else {
				value := parts[1]
				for _, part := range parts[2:] {
					value += ":" + part
				}
				message[parts[0]] = value
			}
		}
	}

	if len(unknown) > 0 {
		unknownStr := unknown[0] + "\n"
		for _, str := range unknown[1:] {
			unknownStr += str + "\n"
		}
		message["?"] = unknownStr
	}

	jmsg, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return logger.sendMessage(jmsg)
}

func (logger *Logger) sendMessage(msg []byte) error {
	req, err := http.NewRequest("POST", logger.URL+"/"+logger.AppToken+"/"+logger.Type, bytes.NewBuffer(msg))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return err
}
