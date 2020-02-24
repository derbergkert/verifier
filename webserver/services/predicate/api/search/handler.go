package search

import (
	"github.com/theonlyrob/vercer/webserver/services/predicate/index"
	"github.com/theonlyrob/vercer/webserver/services/predicate/store"
	"net/http"
)

// Response is the returned value when predicates are searched.
type Response struct {
	Catalysts []*store.Catalyst `json:"predicates"`
}

// NewHandler returns a new handler that adds predicates.
func NewHandler(
	authorizer Authorizer,
	predicateS index.Searcher,
	predicateR store.Reader,
) http.Handler {
	return &handlerImpl{
		authorizer: authorizer,

		predicateS: predicateS,
		predicateR: predicateR,
	}
}
