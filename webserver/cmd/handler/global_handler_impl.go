package handler

import (
	"net/http"

	globalinterceptor "github.com/theonlyrob/vercer/webserver/cmd/interceptor"
)

type globalHandlerImpl struct {
	interceptor globalinterceptor.Interceptor
	mux         *http.ServeMux
}

func (g *globalHandlerImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Use global interceptor.
	w, r = g.interceptor.Intercept(w, r)

	// Run mux.
	g.mux.ServeHTTP(w, r)
}
