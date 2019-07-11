package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rdoorn/gohelper/logging"
	"github.com/rdoorn/gohelper/signaling"
)

type WebserverConfig struct {
	IP        string
	Port      int
	Handler   http.Handler
	TLSConfig *tls.Config
}

type Webserver struct {
	App
	Config     *WebserverConfig
	server     http.Server
	serverDone chan struct{}
	signal     *signaling.Handler
}

func (w *Webserver) Start() error {
	if w.logging == nil {
		panic("webserver not initialized")
	}
	w.logging.Println("Server starting")
	// router.POST("/opvragen/naw", handler.ClientRequest)
	// router.POST("/1bnr", handler.IbanRequest)

	listenAddr := fmt.Sprintf("%s:%d", w.Config.IP, w.Config.Port)
	// start https server
	w.server = http.Server{
		Addr:      listenAddr,
		Handler:   w.Config.Handler,
		TLSConfig: w.Config.TLSConfig,
	}

	w.logging.Println("Server is ready to handle requests", "addr", listenAddr)
	if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		w.logging.Fatalf("Could not start listener", "addr", listenAddr, "error", err)
		return err
	}
	<-w.serverDone
	return nil
}

func (w *Webserver) Stop() error {
	w.logging.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	w.server.SetKeepAlivesEnabled(false)
	if err := w.server.Shutdown(ctx); err != nil {
		w.logging.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	close(w.serverDone)
	w.signal.Close()

	return nil
}

func (w *Webserver) Init() (err error) {
	w.serverDone = make(chan struct{})
	w.logging, err = logging.NewZap("stdout")
	w.signal = signaling.New()
	return
}

func (w *Webserver) Reload() error {
	return nil
}

func (w *Webserver) LoadWebserverConfig(config WebserverConfig) error {
	w.Config = &config
	return nil
}

func (w *Webserver) Signal(f func(), sig ...os.Signal) {
	w.signal.Add(f, sig...)
}
