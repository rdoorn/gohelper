package app

import (
	"github.com/rdoorn/gohelper/logging"
)

type AppInterface interface {
	Configure() error
	Reload() error
	Stop() error
	Start() error

	logging.SimpleLogger
}

func New(i AppInterface) (AppInterface, error) {
	err := i.Configure()
	return i, err
}

type App struct {
	logging logging.SimpleLogger
}

func (a *App) Println(v ...interface{}) {
	a.logging.Println(v...)
}

func (a *App) Debugf(v ...interface{}) {
	a.logging.Debugf(v...)
}

func (a *App) Infof(v ...interface{}) {
	a.logging.Infof(v...)
}

func (a *App) Warnf(v ...interface{}) {
	a.logging.Warnf(v...)
}

func (a *App) Errorf(v ...interface{}) {
	a.logging.Errorf(v...)
}

func (a *App) Fatalf(v ...interface{}) {
	a.logging.Fatalf(v...)
}

func (a *App) Panicf(v ...interface{}) {
	a.logging.Panicf(v...)
}
