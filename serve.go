package twig

import (
	"context"
	"net/http"
	"os"
)

type Worker interface {
	Attacher
	Cycler
	Handler(http.Handler)
}

type Work struct {
	*http.Server
	twig *Twig
}

func NewWork() *Work {
	return &Work{
		Server: &http.Server{
			Addr:     DefaultAddress,
			ErrorLog: newLog(os.Stderr, "twig-work-"),
		},
	}
}

func (w *Work) Handler(h http.Handler) {
	w.Server.Handler = h
}

func (w *Work) Attach(twig *Twig) {
	w.twig = twig
	w.Handler(twig)
}

func (w *Work) Shutdown(ctx context.Context) error {
	return w.Server.Shutdown(ctx)
}

func (w *Work) Start() (err error) {
	go func() {
		err = w.Server.ListenAndServe()
	}()
	return
}
