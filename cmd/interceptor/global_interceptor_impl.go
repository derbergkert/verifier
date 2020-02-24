package interceptor

import (
	"net/http"
)

type globalInterceptorImpl struct {
	interceptors []Interceptor
}

func (i *globalInterceptorImpl) Intercept(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
	for _, ceptor := range i.interceptors {
		w, r = ceptor.Intercept(w, r)
	}
	return w, r
}
