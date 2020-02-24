package add

import (
	"encoding/json"
	"github.com/theonlyrob/vercer/webserver/services/predicate/notifier"
	"github.com/theonlyrob/vercer/webserver/services/predicate/store"
	"net/http"

	"github.com/theonlyrob/vercer/webserver/pkg/api"
)

type handlerImpl struct {
	authorizer Authorizer
	validator  Validator

	predicateW store.Writer
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

	// Fill in the predicate and upsert to DB.
	predicate := request.Catalyst
	if err := l.predicateW.Add(*predicate); err != nil {
		api.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Hit notifier.
	l.notify.Added(*predicate)

	// Return nothing except valid code.
	json.NewEncoder(w).Encode(&Response{
		CatalystID: predicate.ID,
	})
}
