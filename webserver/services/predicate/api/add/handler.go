package add

import (
	"github.com/theonlyrob/vercer/webserver/services/predicate/notifier"
	"github.com/theonlyrob/vercer/webserver/services/predicate/store"
	"net/http"
)

// Request is a request to add a predicate.
type Request struct {
	Catalyst *store.Catalyst `json:"predicate"`
}

// Response is the returned value when a predicate is added.
type Response struct {
	CatalystID string `json:"predicateId"`
}

// NewHandler returns a new handler that adds predicates.
func NewHandler(
	authorizer Authorizer,
	validator Validator,
	predicateW store.Writer,
	notify notifier.Notifier,
) http.Handler {
	return &handlerImpl{
		authorizer: authorizer,
		validator:  validator,

		predicateW:  predicateW,
		notify: notify,
	}
}
