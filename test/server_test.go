package test

import (
	"io"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	sematextlogger "github.com/Ephram84/sematext-logger"
	"github.com/Ephram84/sematext-logger/logwriter"
)

const (
	appToken string = "3cb2be30-05c6-45d6-bdc9-075cac545206"
)

func TestServer(t *testing.T) {
	ts := httptest.NewServer(GetAPI())
	defer ts.Close()

	answer, _, err := SendRequest("GET", ts.URL+"/testRoute?isError=true", nil, nil)

	if err != nil {
		t.Fatal(err)
	}


	println(string(answer))
}

func TestDialSematext(t *testing.T) {
	logger := sematextlogger.NewLogger(appToken, "test")

	// logger.Err("An error has occurred")
	logger.Info("Info")
}

func TestLogger(t *testing.T) {
	sematext, _ := logwriter.Dial("udp", "logsene-receiver-syslog.sematext.com:514", logwriter.LOG_LOCAL0, appToken)

	multi := io.MultiWriter(sematext, os.Stdout)

	info := log.New(multi, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	info.Println("An error has occurred")
}

func TestEnv(t *testing.T) {
	println(os.Getenv("LOGGING_URL"))
}
