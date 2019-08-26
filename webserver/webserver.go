package webserver

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/rdoorn/gohelper/logging"
	"github.com/rdoorn/gohelper/signaling"
)

type Server struct {
	addr       string
	tlsConfig  *tls.Config
	handler    http.Handler
	server     http.Server
	serverDone chan struct{}
	logger     logging.SimpleLogger
	signals    *signaling.Handler
}

type Config struct {
	Addr      string
	TLSConfig *tls.Config
}

func New(config Config) *Server {
	logger, _ := logging.NewZap("stdout")
	srv := Server{
		addr:       config.Addr,
		tlsConfig:  config.TLSConfig,
		logger:     logger,
		serverDone: make(chan struct{}),
		signals:    signaling.New(),
	}
	srv.signals.Add(srv.stop, os.Interrupt, syscall.SIGTERM)
	return &srv
}

func (s *Server) WithHandler(handler http.Handler) {
	s.handler = handler
}

func (s *Server) Start() error {

	// start https server
	s.server = http.Server{
		Addr:      s.addr,
		Handler:   s.handler,
		TLSConfig: s.tlsConfig,
	}

	s.logger.Infof("Webservice started", "listener", s.addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	<-s.serverDone
	return nil
}

func (s *Server) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.server.SetKeepAlivesEnabled(false)
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	close(s.serverDone)
	s.logger.Infof("Webservice shutdown")
}

func (a *Server) Println(v ...interface{}) {
	a.logger.Println(v...)
}

func (a *Server) Debugf(v ...interface{}) {
	a.logger.Debugf(v...)
}

func (a *Server) Infof(v ...interface{}) {
	a.logger.Infof(v...)
}

func (a *Server) Warnf(v ...interface{}) {
	a.logger.Warnf(v...)
}

func (a *Server) Errorf(v ...interface{}) {
	a.logger.Errorf(v...)
}

func (a *Server) Fatalf(v ...interface{}) {
	a.logger.Fatalf(v...)
}

func (a *Server) Panicf(v ...interface{}) {
	a.logger.Panicf(v...)
}
