package app

import (
	"os"

	"github.com/rdoorn/gohelper/logging"
	"github.com/rdoorn/gohelper/profiling"
	"github.com/rdoorn/gohelper/signaling"
)

type WebserverInterface interface {
	Start() error
	Stop() error
}

type App struct {
	Name      string
	logger    logging.SimpleLogger
	profiler  profiling.Interface
	webserver *WebserverHandler
	signals   *signaling.Handler
}

func (app *App) New(opts ...Option) error {
	logger, _ := logging.NewZap("stdout")

	app.logger = logger

	if addr, ok := os.LookupEnv("PROFILING"); ok {
		app.profiler = profiling.Default(addr)
	}

	for _, o := range opts {
		o(app)
	}

	return nil
}

func (a *App) Println(v ...interface{}) {
	a.logger.Println(v...)
}

func (a *App) Debugf(v ...interface{}) {
	a.logger.Debugf(v...)
}

func (a *App) Infof(v ...interface{}) {
	a.logger.Infof(v...)
}

func (a *App) Warnf(v ...interface{}) {
	a.logger.Warnf(v...)
}

func (a *App) Errorf(v ...interface{}) {
	a.logger.Errorf(v...)
}

func (a *App) Fatalf(v ...interface{}) {
	a.logger.Fatalf(v...)
}

func (a *App) Panicf(v ...interface{}) {
	a.logger.Panicf(v...)
}

func (a *App) Start() error {
	a.logger.Infof("application starting")
	if a.webserver != nil {
		if err := a.webserver.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) Stop() {
	a.logger.Infof("application stopping")
	if a.webserver != nil {
		a.webserver.Stop()
	}
}
