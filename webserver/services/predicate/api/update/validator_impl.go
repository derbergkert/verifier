package update

import (
	"context"
	"errors"
	"github.com/theonlyrob/vercer/webserver/services/predicate/store"
	"time"
)

type validatorImpl struct {
	predicateR store.Reader
}

func (val *validatorImpl) Validate(ctx context.Context, req *Request) error {
	// Load the predicate if it exists.
	predicate, err := val.predicateR.Get(req.Catalyst.ID)
	if err != nil {
		return err
	} else if predicate == nil {
		return errors.New("predicate missing")
	}
	// Update the updated timestamp.
	timestamp := time.Now()
	predicate.UpdatedTime = timestamp

	return nil
}
