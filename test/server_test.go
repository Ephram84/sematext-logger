package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	sematextlogger "github.com/Ephram84/sematext-logger"
	"github.com/labstack/echo"
)

const (
	appToken string = "430fv34u-05c6-45d6-bdc9-634dgez"
	typ      string = "syslog"
)

func TestServer(t *testing.T) {
	ts := httptest.NewServer(getHandler())
	defer ts.Close()

	logger := sematextlogger.NewLogger(appToken).WithType(typ).WithURL(ts.URL)

	err := logger.Err("An error has occurred", "path:/api/example/err")
	check(t, err)

	err = logger.Info("test message", "Methode:GET", "uri:/api/test/info")
	check(t, err)

	err = logger.Info("test message", "GET", "XYZ")
	check(t, err)

	err = logger.Debug("test message", "@timestamp:"+time.Now().String())
	check(t, err)

	err = logger.Debug("test message", time.Now().String())
	check(t, err)

	err = logger.Warning("test message", "Methode:GET", "uri:/api/test/info")
	check(t, err)

	err = logger.Notice("test message", "Methode:GET", "uri:/api/test/info")
	check(t, err)

	err = logger.Crit("test message", "Methode:GET", "uri:/api/test/info")
	check(t, err)

	err = logger.Emerg("test message", "Methode:GET", "uri:/api/test/info")
	check(t, err)
}

func TestServerWithEnv(t *testing.T) {
	//set enviroment variable for testing
	os.Setenv("LOGGING_URL", "https://logsene-receiver.sematext.com/430fv34u-05c6-45d6-bdc9-634dgez/")

	ts := httptest.NewServer(getHandler())
	defer ts.Close()

	logger := sematextlogger.InitLogger("LOGGING_URL").WithType(typ).WithURL(ts.URL)

	err := logger.Err("An error has occurred", "path:/api/example/err")
	check(t, err)
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

//////////////////////////////////////
//				Dummy				//
//////////////////////////////////////
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
