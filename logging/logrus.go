package logging

import (
	"io/ioutil"
	"log/syslog"
	"os"

	"github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

type Logrus struct {
	Logger *logrus.Logger
}

func (f *Logrus) Println(v ...interface{}) {
	f.Logger.Println(v...)
}

func (f *Logrus) Debugf(v ...interface{}) {
	f.Logger.Debugf(v[0].(string), v[1:]...)
}

func (f *Logrus) Infof(v ...interface{}) {
	f.Logger.Infof(v[0].(string), v[1:]...)
}

func (f *Logrus) Warnf(v ...interface{}) {
	f.Logger.Warnf(v[0].(string), v[1:]...)
}

func (f *Logrus) Errorf(v ...interface{}) {
	f.Logger.Errorf(v[0].(string), v[1:]...)
}

func (f *Logrus) Fatalf(v ...interface{}) {
	f.Logger.Fatalf(v[0].(string), v[1:]...)
}

func (f *Logrus) Panicf(v ...interface{}) {
	f.Logger.Panicf(v[0].(string), v[1:]...)
}

func NewLogrus(dst ...string) (SimpleLogger, error) {
	logger := logrus.New()

	for _, d := range dst {
		switch d {
		case "stdout", "":
			logger.Out = os.Stdout
			logger.Formatter = &logrus.TextFormatter{DisableColors: false, DisableTimestamp: false, QuoteEmptyFields: true}
		case "stderr":
			logger.Out = os.Stderr
			logger.Formatter = &logrus.TextFormatter{DisableColors: false, DisableTimestamp: false, QuoteEmptyFields: true}
		case "syslog":
			logger.Formatter = &logrus.TextFormatter{DisableColors: true, DisableTimestamp: true, QuoteEmptyFields: true}
			hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_LOCAL5, "")

			if err == nil {
				logger.Hooks.Add(hook)
				logger.Out = ioutil.Discard
			}

		default:
			logger.Formatter = &logrus.TextFormatter{DisableColors: true, DisableTimestamp: false, QuoteEmptyFields: true}
			f, err := os.OpenFile(d, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
			if err != nil {
				logrus.Fatal(err)
			}
			logger.Out = f
		}
	}

	return &Logrus{
		Logger: logger,
	}, nil
}
