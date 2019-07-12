package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

type WebserverConfig struct {
	IP        string
	Port      int
	Handler   http.Handler
	TLSConfig *tls.Config
}

type WebserverHandler struct {
	config     *WebserverConfig
	server     http.Server
	serverDone chan struct{}
}

func NewWebserverHandler(c WebserverConfig) *WebserverHandler {
	return &WebserverHandler{
		config:     &c,
		serverDone: make(chan struct{}),
	}
}

func (w *WebserverHandler) Start() error {
	// router.POST("/opvragen/naw", handler.ClientRequest)
	// router.POST("/1bnr", handler.IbanRequest)

	listenAddr := fmt.Sprintf("%s:%d", w.config.IP, w.config.Port)
	// start https server
	w.server = http.Server{
		Addr:      listenAddr,
		Handler:   w.config.Handler,
		TLSConfig: w.config.TLSConfig,
	}

	if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	<-w.serverDone
	return nil
}

func (w *WebserverHandler) Stop() error {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	w.server.SetKeepAlivesEnabled(false)
	if err := w.server.Shutdown(ctx); err != nil {
		w.logging.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	close(w.serverDone)

	return nil
}

func (w *WebserverHandler) Reload() error {
	return nil
}

/*
func (w *WebserverHandler) LoadWebserverConfig(config WebserverConfig) error {
	w.Config = &config
	return nil
}
*/
/*
func (w *WebserverHandler) Signal(f func(), sig ...os.Signal) {
	w.signal.Add(f, sig...)
}
*/
