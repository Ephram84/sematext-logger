package sematextlogger

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

// WinLogger is a writer to sematext server.
type WinLogger struct {
	AppToken string
	Type     string
	Host     string
	URL      string
	Keys     []string
}

// NewLogger returns a logger
func NewLogger(appToken, typ, url string, keys ...string) *WinLogger {
	host, _ := os.Hostname()
	if typ == "" {
		typ = "syslog"
	}
	if url == "" {
		url = "http://logsene-receiver.sematext.com:80"
	}
	return &WinLogger{AppToken: appToken, Type: typ, Host: host, URL: url, Keys: keys}
}

// NewKeys sets or replaces all keys.
func (logger *WinLogger) NewKeys(newkeys ...string) {
	logger.Keys = newkeys
}

// AddKey adds a new key, except it already exists.
func (logger *WinLogger) AddKey(key string) error {
	pos := findPosition(logger.Keys, key)
	if pos < 0 {
		logger.Keys = append(logger.Keys, key)
	} else {
		return errors.New(key + " already exists")
	}
	return nil
}

// RemoveKey removes <key> if it exists.
func (logger *WinLogger) RemoveKey(key string) error {
	pos := findPosition(logger.Keys, key)
	if pos >= 0 {
		logger.Keys = append(logger.Keys[:pos], logger.Keys[pos+1:]...)
	} else {
		return errors.New(key + " has not been found")
	}
	return nil
}

func findPosition(keys []string, key string) int {
	for i, k := range keys {
		if k == key {
			return i
		}
	}
	return -1
}

// Err logs a message with severity "err".
func (logger *WinLogger) Err(msg string, values ...interface{}) (bool, error) {
	return logger.buildMessage("err", msg, values)
}

// Info logs a message with severity "info".
func (logger *WinLogger) Info(msg string, values ...interface{}) (bool, error) {
	return logger.buildMessage("info", msg, values)
}

// Emerg logs a message with severity "emerg".
func (logger *WinLogger) Emerg(msg string, values ...interface{}) (bool, error) {
	return logger.buildMessage("emerg", msg, values)
}

// Crit logs a message with severity "crit".
func (logger *WinLogger) Crit(msg string, values ...interface{}) (bool, error) {
	return logger.buildMessage("crit", msg, values)
}

// Warning logs a message with severity "warning".
func (logger *WinLogger) Warning(msg string, values ...interface{}) (bool, error) {
	return logger.buildMessage("warning", msg, values)
}

// Notice logs a message with severity "notice".
func (logger *WinLogger) Notice(msg string, values ...interface{}) (bool, error) {
	return logger.buildMessage("notice", msg, values)
}

// Debug logs a message with severity "debug".
func (logger *WinLogger) Debug(msg string, values ...interface{}) (bool, error) {
	return logger.buildMessage("debug", msg, values)
}

func (logger *WinLogger) buildMessage(severity, msg string, values []interface{}) (bool, error) {
	if len(logger.Keys) != len(values) {
		return false, errors.New("Size of keys and values are odd")
	}
	var message map[string]interface{}
	message = make(map[string]interface{})

	message["Severity"] = severity
	message["Message"] = msg
	message["Host"] = logger.Host

	for i, key := range logger.Keys {
		message[key] = values[i]
	}

	jmsg, _ := json.Marshal(message)
	return logger.sendMessage(jmsg)
}

func (logger *WinLogger) sendMessage(msg []byte) (bool, error) {
	req, err := http.NewRequest("POST", logger.URL+"/"+logger.AppToken+"/"+logger.Type, bytes.NewBuffer(msg))
	if err != nil {
		return false, err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return false, err
	}

	return true, err
}
