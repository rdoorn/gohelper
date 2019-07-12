package app

import (
	"os"

	"github.com/rdoorn/gohelper/logging"
	"github.com/rdoorn/gohelper/profiling"
	"github.com/rdoorn/gohelper/signaling"
)

type Option func(o *App)

func WithLogging(l logging.SimpleLogger) Option {
	return func(o *App) {
		o.logger = l
	}
}

func WithProfiling(p profiling.Interface) Option {
	return func(o *App) {
		o.profiler = p
	}
}

func Webserver(c WebserverConfig) Option {
	return func(o *App) {
		if o.webserver == nil {
			o.webserver = NewWebserverHandler(o.logger, c)
		}
	}
}

func Signal(f func(), sigs ...os.Signal) Option {
	return func(o *App) {
		if o.signals == nil {
			o.signals = signaling.New()
		}

		o.signals.Add(f, sigs...)
	}
}

/*
func WithWebServer(c WebserverConfig) Option {
	return func(o *App) {
		o.webserverConfig = c
	}
}
*/
