package delete

import (
	"context"
	"errors"
	"github.com/theonlyrob/vercer/webserver/services/predicate/store"
)

type validatorImpl struct {
	predicateR store.Reader
}

func (val *validatorImpl) Validate(ctx context.Context, req *Request) error {
	// Load the predicate if it exists.
	predicate, err := val.predicateR.Get(req.CatalystID)
	if err != nil {
		return err
	} else if predicate == nil {
		return errors.New("predicate missing")
	}
	return nil
}
