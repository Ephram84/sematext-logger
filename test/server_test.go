package test

import (
	"fmt"
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

	logger := sematextlogger.NewLogger(appToken).WithType(typ).WithURL(ts.URL)

	ok, err := logger.Err("An error has occurred", "path:/api/example/err")
	check(t, ok, err)

	ok, err = logger.Info("test message", "Methode:GET", "uri:/api/test/info")
	check(t, ok, err)

	ok, err = logger.Debug("test message", "@timestamp:"+time.Now().String())
	check(t, ok, err)

	ok, err = logger.Warning("test message", "Methode:GET", "uri:/api/test/info")
	check(t, ok, err)

	ok, err = logger.Notice("test message", "Methode:GET", "uri:/api/test/info")
	check(t, ok, err)

	ok, err = logger.Crit("test message", "Methode:GET", "uri:/api/test/info")
	check(t, ok, err)

	ok, err = logger.Emerg("test message", "Methode:GET", "uri:/api/test/info")
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
