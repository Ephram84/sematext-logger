package test

import (
	"net/http/httptest"
	"os"
	"testing"

	sematextlogger "github.com/Ephram84/sematext-logger"
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
	logger := sematextlogger.NewLogger(os.Getenv("appToken"), "test")

	// logger.Err("An error has occurred")
	logger.Info("Info", "It's just a info message")
}
