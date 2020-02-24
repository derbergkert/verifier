package add

import (
	"context"
	"errors"
	"time"

	"github.com/theonlyrob/vercer/webserver/services/predicate/store"
	"github.com/google/uuid"
)

type validatorImpl struct{}

func (val *validatorImpl) Validate(ctx context.Context, req *Request) error {
	predicate := req.Catalyst

	// Allocate a UUID for the predicate.
	lid, err := uuid.NewRandom()
	if err != nil {
		return errors.New("unable to generate id")
	}
	predicate.ID = lid.String()

	// Update the updated timestamp.
	timestamp := time.Now()
	predicate.CreatedTime = timestamp
	predicate.UpdatedTime = timestamp

	return store.Validate(req.Catalyst)
}
