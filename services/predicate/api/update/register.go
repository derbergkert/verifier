package update

import (
	"net/http"
)

// Register adds the http handler to the input mux under /add.
func Register(mux *http.ServeMux) {
	mux.Handle("/update", SingletonHandler())
}