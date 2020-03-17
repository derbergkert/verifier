package ping

import (
	"github.com/theonlyrob/vercer/webserver/services/liveness/notifier"
	"github.com/theonlyrob/vercer/webserver/services/liveness/store"
	"net/http"
)

// Request is a request.
type Request struct {}

// Response is the returned.
type Response struct {}

// NewHandler returns a new handler.
func NewHandler(
	authorizer Authorizer,
	validator Validator,
) http.Handler {
	return &handlerImpl{
		authorizer: authorizer,
		validator:  validator,
	}
}
