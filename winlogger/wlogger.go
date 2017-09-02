package winlogger

type WinLogger struct {
	AppToken string
	Type     string
}

func NewLogger(appToken, typ string) *WinLogger {
	return &WinLogger{AppToken: appToken, Type: typ}
}

func (logger *WinLogger) Err(msg string) error {

}

func (logger *WinLogger) Info(msg string) error {

}

func (logger *WinLogger) Emerg(msg string) error {

}

func (logger *WinLogger) Crit(msg string) error {

}

func (logger *WinLogger) Warning(msg string) error {

}

func (logger *WinLogger) Notice(msg string) error {

}

func (logger *WinLogger) Debug(msg string) error {

}
