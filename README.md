# SEMATEXT-LOGGER

[![Build Status](https://travis-ci.org/Ephram84/sematext-logger.svg?branch=master)](https://travis-ci.org/Ephram84/sematext-logger)

With sematext-logger, log events can be send to [Sematext](https://sematext.com/).

## Install
<code>go get github.com/Ephram84/sematext-logger</code>

## Usage
First create a new logger:
```golang
logger := sematextlogger.NewLogger(appToken).WithType(typ).WithURL(url)
```
Parameter:
* appToken (string, required): Is the token of your Logsene app. You will find it on [sematext website](https://apps.sematext.com/ui/logs).
* typ (string, optional): Typ is a logical division of events and can be anything. For example, for syslog messages typ would be called "syslog". But also the name of your app would be possible. Default is syslog.
* url (string, optional): If url is not set, all messages are sent to http://logsene-receiver.sematext.com:80. With <code>WithURL(any URL)</code>, you can specify your own URL.

### Send message
logger can now be used to call different methods that specify different severities.
For example, the following code
```golang
logger.Err("An error has occurred")
```

then produces this output on sematext</br>
![](pictures/Sematext_err.PNG?raw=true)

### Additional information
If you want to send further information, apart from the message, you can add by adding strings in the 'key:value' format of the corresponding method.
For example
```golang
logger.Err("An error has occurred", "path:/api/example/err")
```
looks like on sematext</br>
![](pictures/Sematext_err2.PNG?raw=true)
