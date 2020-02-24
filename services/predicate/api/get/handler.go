package get

import (
	"github.com/theonlyrob/vercer/webserver/services/predicate/notifier"
	"github.com/theonlyrob/vercer/webserver/services/predicate/store"
	"net/http"
)

// Request is a request to add a predicate.
type Request struct {
	CatalystIDs []string `json:"predicateId"`
}

// Response is the returned value when a predicate is added.
type Response struct {
	Catalysts []*store.Catalyst `json:"predicates"`
}

// NewHandler returns a new handler that adds predicates.
func NewHandler(
	authorizer Authorizer,
	validator Validator,
	predicateR store.Reader,
	notify notifier.Notifier,
) http.Handler {
	return &handlerImpl{
		authorizer: authorizer,
		validator:  validator,

		predicateR:  predicateR,
		notify: notify,
	}
}
