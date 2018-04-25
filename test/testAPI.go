package test

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	sematextlogger "github.com/Ephram84/sematext-logger"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type TestContext struct {
	echo.Context
	Sematextlogger *sematextlogger.Logger
}

func GetAPI() *echo.Echo {

	logger := sematextlogger.InitLogger("LOGGING_URL", "test")

	router := echo.New()
	router.HideBanner = true

	router.Use(logger.EchoMiddlwareLogger())
	router.Use(middleware.RequestID())
	router.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tcontext := &TestContext{c, logger}
			return h(tcontext)
		}
	})

	router.GET("testRoute", handleMessage)

	return router
}

func handleMessage(c echo.Context) error {
	context := c.(*TestContext)

	context.Sematextlogger.Info(context.Response().Header().Get(echo.HeaderXRequestID), "handleMessage")

	isError := context.QueryParam("isError")
	if isError == "true" {
		context.Sematextlogger.Error(context.Response().Header().Get(echo.HeaderXRequestID), "An error has occurred", " - ", errors.New("Example error").Error())
		return context.JSON(http.StatusConflict, "An error has occurred")
	}

	return context.JSON(http.StatusOK, "")
}

func SendRequest(method, url string, body io.Reader, header map[string]string) ([]byte, int, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, err
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return bytes, resp.StatusCode, nil
}
