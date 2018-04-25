# SEMATEXT-LOGGER

[![Build Status](https://travis-ci.org/Ephram84/sematext-logger.svg?branch=master)](https://travis-ci.org/Ephram84/sematext-logger)

<<<<<<< HEAD
With sematext-logger, log events can be send to [Sematext](https://sematext.com/).
=======
With sematext-logger, log events can be send to [Sematext](https://sematext.com/) and output to the console at the same time.
>>>>>>> syslog

## Install
<code>go get github.com/Ephram84/sematext-logger</code>

## Usage
First create a new logger:
```golang
<<<<<<< HEAD
logger := sematextlogger.NewLogger(appToken).WithType(typ).WithURL(url)
```
Parameter:
* appToken (string, required): Is the token of your Logsene app. You will find it on [sematext website](https://apps.sematext.com/ui/logs).
* typ (string, optional): Typ is a logical division of events and can be anything. For example, for syslog messages typ would be called "syslog". But also the name of your app would be possible. Default is syslog.
* url (string, optional): If url is not set, all messages are sent to http://logsene-receiver.sematext.com:80. With <code>WithURL(any URL)</code>, you can specify your own URL.
=======
logger := sematextlogger.NewLogger(appToken, service)
```
Parameter:
* appToken (string, required): Is the token of your Logsene app. You will find it on [sematext website](https://apps.sematext.com/ui/logs).
* service (string, required): Typ is a logical division of events and can be anything. For example, for syslog messages typ would be called "syslog". But also the name of your app would be possible. Default is syslog.
>>>>>>> syslog

If you have an enviroment variable that contains a url with an appToken, e.g., https://logsene-receiver.sematext.com/fzr64ktn-...., you can use:
```golang
logger := sematextlogger.InitLogger(envVarName, service)
```
### Send message
logger can now be used to call different methods that specify different severities.
For example, the following code
```golang
logger.Error("any ID", "An error has occurred", "-", errors.New("Example error").Error())
```
produces this output
```{"request_id":"65CMp4TVdaZbJUCAJQGLopWsZUPhaKlp","time":"2018-04-07T10:29:15+02:00","loglevel":"ERROR","service":"test","file":"testAPI.go","line":"47","message":"An error has occurred - Example error"}```

Please note that "any ID" is any alphanumeric string, e.g. the RequestID to easily find the message at Sematext. The remaining parameters form the message.

## Using Echo
If you want to implement a REST API with [Echo](https://echo.labstack.com/), process following steps.

<<<<<<< HEAD
then produces this output on sematext</br>
![](pictures/Sematext_err.PNG?raw=true)

### Additional information
If you want to send further information, apart from the message, you can add by adding several strings in the 'key:value' format of the corresponding method.
For example
```golang
logger.Err("An error has occurred", "path:/api/example/err")
```
looks like on sematext</br>
![](pictures/Sematext_err2.PNG?raw=true)
=======
1. Import sematext:
```golang
import(
    ...
    sematextlogger "github.com/Ephram84/sematext-logger"
    ...
)
```

2. Define a custom context:
```golang
type TestContext struct {
	echo.Context
	Sematextlogger *sematextlogger.Logger
}
```

3. Create a new sematext logger and a middleware to extend default context:
```golang
  logger := sematextlogger.NewLogger(os.Getenv("LOGGING_URL"), "test")

  router := echo.New()
  router.HideBanner = true

  router.Use(logger.EchoMiddlwareLogger()) // Replaces Echo's default logger with the difference that the output is also sent to Sematext.
  router.Use(middleware.RequestID())  // Automatically generates a request ID.
  router.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
      tcontext := &TestContext{c, logger}
      return h(tcontext)
    }
  })
```

4. Send a message within a HandlerFunc:
```golang
context := c.(*TestContext)
...
context.Sematextlogger.Error(context.Response().Header().Get(echo.HeaderXRequestID), "An error has occurred", " - ", errors.New("Example error").Error())
return context.JSON(http.StatusConflict, "An error has occurred")
```
>>>>>>> syslog
