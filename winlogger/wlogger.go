package winlogger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type WinLogger struct {
	AppToken string
	Type     string
	Host     string
}

func NewLogger(appToken, typ string) *WinLogger {
	host, _ := os.Hostname()
	return &WinLogger{AppToken: appToken, Type: typ, Host: host}
}

func (logger *WinLogger) Err(msg string) error {
	return logger.buildMessage("err", msg)
}

func (logger *WinLogger) Info(msg string) error {
	return logger.buildMessage("info", msg)
}

func (logger *WinLogger) Emerg(msg string) error {
	return logger.buildMessage("emerg", msg)
}

func (logger *WinLogger) Crit(msg string) error {
	return logger.buildMessage("crit", msg)
}

func (logger *WinLogger) Warning(msg string) error {
	return logger.buildMessage("warning", msg)
}

func (logger *WinLogger) Notice(msg string) error {
	return logger.buildMessage("notice", msg)
}

func (logger *WinLogger) Debug(msg string) error {
	return logger.buildMessage("debug", msg)
}

func (logger *WinLogger) buildMessage(severity, msg string) error {
	var message map[string]interface{}
	message = make(map[string]interface{})

	message["Severity"] = severity
	message["Message"] = msg
	message["Host"] = logger.Host

	jmsg, _ := json.Marshal(message)
	return logger.sendMessage(jmsg)
}

func (logger *WinLogger) sendMessage(msg []byte) error {
	req, err := http.NewRequest("POST", "http://logsene-receiver.sematext.com:80/"+logger.AppToken+"/"+logger.Type, bytes.NewBuffer(msg))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
