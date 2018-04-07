# SEMATEXT-LOGGER

[![Build Status](https://travis-ci.org/Ephram84/sematext-logger.svg?branch=master)](https://travis-ci.org/Ephram84/sematext-logger)

With sematext-logger, log events can be send to [Sematext](https://sematext.com/) and output to the console at the same time.

## Install
<code>go get github.com/Ephram84/sematext-logger</code>

## Usage
First create a new logger:
```golang
logger := sematextlogger.NewLogger(appToken, service)
```
Parameter:
* appToken (string, required): Is the token of your Logsene app. You will find it on [sematext website](https://apps.sematext.com/ui/logs).
* service (string, required): Typ is a logical division of events and can be anything. For example, for syslog messages typ would be called "syslog". But also the name of your app would be possible. Default is syslog.

If you have an enviroment variable that contains a url with an appToken, e.g., https://logsene-receiver.sematext.com/fzr64ktn-...., you can use:
```golang
logger := sematextlogger.InitLogger(envVarName, service)
```
### Send message
Logger can now be used to call different methods that specify different severities.
For example, the following code
```golang
logger.Error("any ID", "An error has occurred", "-", err)
```
produces this output
```{"request_id":"65CMp4TVdaZbJUCAJQGLopWsZUPhaKlp","time":"2018-04-07T10:29:15+02:00","loglevel":"ERROR","service":"test","file":"testAPI.go","line":"47","message":"An error has occurred - Example error"}```

Please note that "any ID" is any alphanumeric string, e.g. the RequestID to easily find the message at Sematext. The remaining parameters form the message.
