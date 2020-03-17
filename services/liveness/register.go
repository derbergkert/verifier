package liveness

import (
	get2 "github.com/theonlyrob/vercer/webserver/services/liveness/ping"
	"net/http"
)

func Register(mux *http.ServeMux) {
	get2.Register(mux)
}
