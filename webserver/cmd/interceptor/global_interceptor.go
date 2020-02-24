package interceptor

import (
	"net/http"
)

type Interceptor interface {
	Intercept(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request)
}
