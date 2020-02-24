package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type globalServerImpl struct{}

func (g *globalServerImpl) Run(handler http.Handler) error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Good practice to set timeouts to avoid Slowloris attacks.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}
	// Run our server in a goroutine so that it doesn't block.
	var err error
	go func() {
		log.Printf("Listening for connections on %s\n", srv.Addr)
		err = srv.ListenAndServe()
	}()

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	<-c

	// Create a deadline to wait for.
	wait := time.Second * 15
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)

	return err
}

// Helper class that wraps handler funcs.
/////////////////////////////////////////
type handlerFuncWrapper struct {
	handlerFunc http.HandlerFunc
}

func (h *handlerFuncWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handlerFunc(w, r)
}

func wrapHandlerFunc(handlerFunc http.HandlerFunc) http.Handler {
	return &handlerFuncWrapper{
		handlerFunc: handlerFunc,
	}
}
