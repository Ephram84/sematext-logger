# SEMATEXT-LOGGER

With sematext-logger, log events can be send to [Sematext](https://sematext.com/).

## Install
<code>go get github.com/Ephram84/sematext-logger</code>

## Usage
First create a new logger:
```golang
logger := sematextlogger.NewLogger(appToken, typ, url)
```
Parameter:
* appToken (string, required): Is the token of your Logsene app. You will find it on [sematext website](https://apps.sematext.com/ui/logs).
* typ (string, optional): Typ is a logical division of events and can be anything. For example, for syslog messages typ would be called "syslog". But also the name of your app would be possible. Default is syslog.
* url (string, optional): If url is an empty string, all messages are sent to http://logsene-receiver.sematext.com:80.

### Send message
Logger can now be used to call different methods that specify different severities.
For example, the following code
```golang
logger.Err("An error has occurred")
```

