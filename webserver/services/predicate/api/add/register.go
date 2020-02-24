package add

import (
	"net/http"
)

// Register adds the http handler to the input mux under /add.
func Register(mux *http.ServeMux) {
	mux.Handle("/add", SingletonHandler())
}
