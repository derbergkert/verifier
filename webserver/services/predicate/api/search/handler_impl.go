package search

import (
	"encoding/json"
	"github.com/theonlyrob/vercer/webserver/services/predicate/index"
	"github.com/theonlyrob/vercer/webserver/services/predicate/store"
	"net/http"

	"github.com/theonlyrob/vercer/webserver/pkg/api"
	pkgSearch "github.com/theonlyrob/vercer/webserver/pkg/search"
)

type handlerImpl struct {
	authorizer Authorizer

	predicateS index.Searcher
	predicateR store.Reader
}

// Delete a predicate for a user.
func (l *handlerImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract the request.
	var request pkgSearch.Request
	err := api.ExtractBody(r, &request)
	if err != nil {
		api.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check authorizer
	if err := l.authorizer.Authorize(r.Context(), &request); err != nil {
		api.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Search the predicates.
	predicateIDs, err := l.predicateS.Search(&request)
	if err != nil {
		api.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the predicates.
	predicates, err := l.predicateR.GetBatch(predicateIDs...)
	if err != nil {
		api.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return nothing except valid code.
	json.NewEncoder(w).Encode(&Response{
		Catalysts: predicates,
	})
}
