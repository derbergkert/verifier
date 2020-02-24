package delete

import (
	"github.com/theonlyrob/vercer/webserver/services/predicate/notifier"
	"github.com/theonlyrob/vercer/webserver/services/predicate/store"
	"net/http"

	"github.com/theonlyrob/vercer/webserver/pkg/api"
)

type handlerImpl struct {
	authorizer Authorizer
	validator  Validator

	predicates store.Store
	notify    notifier.Notifier
}

// Add a predicate for a user.
func (l *handlerImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract the request.
	var request Request
	if err := api.ExtractBody(r, &request); err != nil {
		api.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate request.
	if err := l.validator.Validate(r.Context(), &request); err != nil {
		api.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check authorizer
	if err := l.authorizer.Authorize(r.Context(), &request); err != nil {
		api.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Load the predicate if it exists.
	predicate, err := l.predicates.Get(request.CatalystID)
	if err != nil {
		api.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Delete the predicate from the DB.
	if err := l.predicates.Delete(predicate.ID); err != nil {
		api.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Hit notifier.
	l.notify.Deleted(*predicate)
}
