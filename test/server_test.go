package test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sematextlogger "github.com/Ephram84/sematext-logger"
	"github.com/labstack/echo"
)

const (
	appToken string = "3cb2be30-05c6-45d6-bdc9-075cac545206"
	typ      string = "syslog"
)

func TestServer(t *testing.T) {
	ts := httptest.NewServer(getHandler())
	defer ts.Close()

	logger := sematextlogger.NewLogger(appToken, typ, ts.URL, "method", "path")

	ok, err := logger.Info("test message", "GET", "/api/test/info")
	check(t, ok, err)

	logger.NewKeys("time")
	ok, err = logger.Debug("test message", time.Now())
	check(t, ok, err)

	logger.NewKeys("method", "path")
	ok, err = logger.Warning("test message", "GET", "/api/test/warning")
	check(t, ok, err)

	ok, err = logger.Notice("test message", "GET", "/api/test/notice")
	check(t, ok, err)

	ok, err = logger.Crit("test message", "GET", "/api/test/crit")
	check(t, ok, err)

	ok, err = logger.Emerg("test message", "GET", "/api/test/emerg")
	check(t, ok, err)
}

func check(t *testing.T, ok bool, err error) {
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("Status code != 200")
	}
}

func getHandler() *echo.Echo {
	e := echo.New()

	e.POST("/"+appToken+"/"+typ, handleMessage)

	return e
}

func handleMessage(context echo.Context) error {
	var msg map[string]string
	msg = make(map[string]string)
	if err := context.Bind(&msg); err != nil {
		fmt.Println(err)
		return context.JSON(http.StatusInternalServerError, err)
	}

	for _, v := range []string{"Host", "Message", "Severity"} {
		value, ok := msg[v]
		if !ok || len(value) == 0 {
			return context.JSON(http.StatusBadRequest, v+" is missing")
		}
	}

	for key, value := range msg {
		fmt.Println(key, ":", value)
	}
	fmt.Println()

	return context.JSON(http.StatusOK, "OK")
}

func TestAddKey(t *testing.T) {
	logger := sematextlogger.NewLogger(appToken, typ, "")

	if len(logger.Keys) > 0 {
		t.Fail()
	}

	err := logger.AddKey("testkey")
	if err != nil {
		t.Fatal(err)
	}

	err = logger.AddKey("testkey")
	if err == nil {
		t.Fatal("logger should have found the key")
	}
}

func TestRemoveKey(t *testing.T) {
	logger := sematextlogger.NewLogger(appToken, typ, "", "testkey1", "testkey2")

	if len(logger.Keys) != 2 {
		t.Fatal("Something went wrong")
	}

	err := logger.RemoveKey("testkey2")
	if err != nil {
		t.Fatal(err)
	}

	if len(logger.Keys) != 1 {
		t.Fatal("Removing testkey2 was not successful")
	}

	err = logger.RemoveKey("testkey2")
	if err == nil {
		t.Fatal("logger should have removed testkey2")
	}
	log.Println(err)

	err = logger.RemoveKey("testkey1")
	if err != nil {
		t.Fatal(err)
	}

	err = logger.RemoveKey("testkey1")
	if err == nil {
		t.Fatal("logger should have removed testkey1")
	}
	log.Println(err)
}
