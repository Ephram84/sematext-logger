package sematextlogger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	isatty "github.com/mattn/go-isatty"
	"github.com/valyala/fasttemplate"

	"os"

	"github.com/Ephram84/sematext-logger/logwriter"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/color"
)

var (
	errSematextWriterIsNil = errors.New("sematextWriter is nil")
)

// Logger is a writer to sematext server.
type (
	Logger struct {
		SematextWriter *logwriter.Writer
		service        string
		level          Lvl
		output         io.Writer
		template       *fasttemplate.Template
		levels         []string
		color          *color.Color
		bufferPool     sync.Pool
		mutex          sync.Mutex
	}

	Lvl uint8

	JSON map[string]interface{}
)

const (
	DEBUG Lvl = iota + 1
	INFO
	WARN
	ERROR
	OFF
)

var (
	// global        = New("-") ???
	defaultHeader = `{"request_id":"${id}","time":"${time_rfc3339}","loglevel":"${level}","service":"${prefix}",` +
		`"file":"${short_file}","line":"${line}"}`
)

// NewLogger returns a logger
func NewLogger(appToken, service string) *Logger {
	sematext, err := logwriter.Dial("udp", "logsene-receiver-syslog.sematext.com:514", logwriter.LOG_LOCAL0, appToken)

	if service == "" {
		service = "syslog"
	}

	newLogger := &Logger{
		SematextWriter: sematext,
		service:        service,
		template:       newTemplate(defaultHeader),
		color:          color.New(),
		bufferPool: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 256))
			},
		},
	}
	newLogger.initLevels()
	if err != nil {
		log.Println("sematextWriter is nil", "-", err)
		newLogger.setOutput(os.Stdout)
	} else {
		newLogger.setOutput(io.MultiWriter(sematext, os.Stdout))
	}

	// set log level
	if os.Getenv("ENVIRONMENT") == "dev" {
		newLogger.SetLevel(DEBUG)
	} else {
		newLogger.SetLevel(INFO)
	}

	return newLogger
}

// InitLogger inits a logger through an enviroment variable that contains a url,
// e.g., https://logsene-receiver.sematext.com/fzr64ktn-....
func InitLogger(envVArName, service string) *Logger {
	loggingURL := os.Getenv(envVArName)
	appToken := ""
	if loggingURL != "" {
		log.Println("found", loggingURL)
		appToken = strings.Replace(loggingURL, "https://logsene-receiver.sematext.com/", "", -1)
		appToken = strings.Replace(appToken, "/", "", -1)
	} else {
		panic("Url has not been found")
	}

	return NewLogger(appToken, service)
}

func (l *Logger) initLevels() {
	l.levels = []string{
		"-",
		l.color.Blue("DEBUG"),
		l.color.Green("INFO"),
		l.color.Yellow("WARN"),
		l.color.Red("ERROR"),
	}
}

func (l *Logger) setOutput(w io.Writer) {
	l.output = w
	if w, ok := w.(*os.File); !ok || !isatty.IsTerminal(w.Fd()) {
		l.DisableColor()
	}
}

func (l *Logger) DisableColor() {
	l.color.Disable()
	l.initLevels()
}

func (l *Logger) SetLevel(v Lvl) {
	l.level = v
}

func newTemplate(format string) *fasttemplate.Template {
	return fasttemplate.New(format, "${", "}")
}

func (l *Logger) SetHeader(h string) {
	l.template = newTemplate(h)
}

func (logger *Logger) EchoMiddlwareLogger() echo.MiddlewareFunc {
	var output io.Writer
	if logger.SematextWriter != nil {
		output = io.MultiWriter(os.Stdout, logger.SematextWriter)
	} else {
		output = os.Stdout
	}

	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return strings.HasSuffix(c.Request().RequestURI, "/ping") || strings.HasSuffix(c.Request().RequestURI, "/ready")
		},
		Output: output,
		Format: `{"time":"${time_rfc3339}", "request_id":"${id}", "remote_ip":"${remote_ip}", "host":"${host}",` +
			` "method":"${method}", "uri":"${uri}", "status":${status}, "latency":${latency},` +
			` "latency_human":"${latency_human}", "bytes_in":${bytes_in},` +
			` "bytes_out":${bytes_out}, "service":"` + logger.service + `"}` + "\n"})
}

// func (l *Logger) Print(i ...interface{}) {
// 	l.log(0, "", i...)
// }

// func (l *Logger) Printf(format string, args ...interface{}) {
// 	l.log(0, format, args...)
// }

// func (l *Logger) Printj(j JSON) {
// 	l.log(0, "json", j)
// }

func (l *Logger) Debug(id string, i ...interface{}) {
	i = append(i, id)
	l.log(DEBUG, "id", i...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

func (l *Logger) Debugj(j JSON) {
	l.log(DEBUG, "json", j)
}

func (l *Logger) Info(id string, i ...interface{}) {
	i = append(i, id)
	l.log(INFO, "id", i...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func (l *Logger) Infoj(j JSON) {
	l.log(INFO, "json", j)
}

func (l *Logger) Warn(id string, i ...interface{}) {
	i = append(i, id)
	l.log(WARN, "id", i...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

func (l *Logger) Warnj(j JSON) {
	l.log(WARN, "json", j)
}

func (l *Logger) Error(id string, i ...interface{}) {
	i = append(i, id)
	l.log(ERROR, "id", i...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func (l *Logger) Errorj(j JSON) {
	l.log(ERROR, "json", j)
}

func (l *Logger) log(v Lvl, format string, args ...interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	buf := l.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer l.bufferPool.Put(buf)
	_, file, line, _ := runtime.Caller(2)

	if v >= l.level {
		message := ""
		id := ""
		switch format {
		case "":
			message = fmt.Sprint(args...)
		case "json":
			msg := args[0].(JSON)
			_, ok := msg[echo.HeaderXRequestID]
			if ok {
				id = msg[echo.HeaderXRequestID].(string)
				delete(msg, echo.HeaderXRequestID)
			}
			b, err := json.Marshal(msg)
			if err != nil {
				panic(err)
			}
			message = string(b)
		case "id":
			id = args[len(args)-1].(string)
			message = fmt.Sprint(args[:len(args)-1]...)
		default:
			message = fmt.Sprintf(format, args...)
		}

		_, err := l.template.ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) {
			switch tag {
			case "id":
				return w.Write([]byte(id))
			case "time_rfc3339":
				return w.Write([]byte(time.Now().Format(time.RFC3339)))
			case "time_rfc3339_nano":
				return w.Write([]byte(time.Now().Format(time.RFC3339Nano)))
			case "level":
				return w.Write([]byte(l.levels[v]))
			case "prefix":
				return w.Write([]byte(l.service))
			case "long_file":
				return w.Write([]byte(file))
			case "short_file":
				return w.Write([]byte(path.Base(file)))
			case "line":
				return w.Write([]byte(strconv.Itoa(line)))
			}
			return 0, nil
		})

		if err == nil {
			s := buf.String()
			i := buf.Len() - 1
			if s[i] == '}' {
				// JSON header
				buf.Truncate(i)
				buf.WriteByte(',')
				if format == "json" {
					buf.WriteString(message[1:])
				} else {
					buf.WriteString(`"message":`)
					buf.WriteString(strconv.Quote(message))
					buf.WriteString(`}`)
				}
			} else {
				// Text header
				buf.WriteByte(' ')
				buf.WriteString(message)
			}
			buf.WriteByte('\n')
			l.output.Write(buf.Bytes())
		}
	}
}
