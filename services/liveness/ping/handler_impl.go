package ping

import (
	"encoding/json"
	"net/http"

	"github.com/theonlyrob/vercer/webserver/pkg/api"
)

type handlerImpl struct {
	authorizer Authorizer
	validator  Validator
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

	// Return nothing except valid code.
	json.NewEncoder(w).Encode(&Response{})
}
