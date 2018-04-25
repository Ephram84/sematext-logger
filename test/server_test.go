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
<<<<<<< HEAD
	appToken string = "430fv34u-05c6-45d6-bdc9-634dgez"
	typ      string = "syslog"
=======
	appToken string = "3cb2be30-05c6-45d6-bdc9-075cac545206"
>>>>>>> syslog
)

func TestServer(t *testing.T) {
	ts := httptest.NewServer(GetAPI())
	defer ts.Close()

<<<<<<< HEAD
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
=======
	answer, _, err := SendRequest("GET", ts.URL+"/testRoute?isError=true", nil, nil)
>>>>>>> syslog
	if err != nil {
		t.Fatal(err)
	}
<<<<<<< HEAD
}

//////////////////////////////////////
//				Dummy				//
//////////////////////////////////////
func getHandler() *echo.Echo {
	e := echo.New()
=======

	println(string(answer))
}
>>>>>>> syslog

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
