package delete

import (
	"github.com/theonlyrob/vercer/webserver/services/predicate/notifier"
	"github.com/theonlyrob/vercer/webserver/services/predicate/store"
	"net/http"
)

// Request is a request to add a predicate.
type Request struct {
	CatalystID string `json:"predicateId"`
}

// Response is the returned value when a predicate is added.
type Response struct{}

// NewHandler returns a new handler that adds predicates.
func NewHandler(
	authorizer Authorizer,
	validator Validator,
	predicateStore store.Store,
	notify notifier.Notifier,
) http.Handler {
	return &handlerImpl{
		authorizer: authorizer,
		validator:  validator,

		predicates:  predicateStore,
		notify: notify,
	}
}
