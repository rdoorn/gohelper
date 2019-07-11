package signaling

import (
	"os"
	"os/signal"
)

type Handler struct {
	close chan struct{}
}

func New() *Handler {
	return &Handler{
		close: make(chan struct{}),
	}
}

func (h *Handler) Add(f func(), sigs ...os.Signal) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, sigs...)

		for {
			select {
			case <-c:
				f()
				return
			case <-h.close:
				return
			}
		}
	}()
}

func (h *Handler) Close() {
	close(h.close)
}
