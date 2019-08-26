package profiling

import (
	"context"
	"log"
	"net/http"
	"net/http/pprof"
	"runtime"
	"time"
)

type DefaultHandler struct {
	addr string
	srv  *http.Server
}

func Default(addr string) Interface {
	return &DefaultHandler{
		addr: addr,
	}
}

func (h *DefaultHandler) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", pprof.Index)

	h.srv = &http.Server{
		Addr:    h.addr,
		Handler: mux,
	}

	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	//http.ListenAndServe(addr, mux)

	go func() {
		// returns ErrServerClosed on graceful close
		if err := h.srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Profiler ListenAndServe(): %s", err)
		}
	}()
}

func (h *DefaultHandler) Stop() {
	if h.srv == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	h.srv.Shutdown(ctx)
	h.srv = nil
}
