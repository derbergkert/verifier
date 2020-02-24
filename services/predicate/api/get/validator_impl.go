package get

import (
	"context"
	"github.com/theonlyrob/vercer/webserver/services/predicate/store"
)

type validatorImpl struct {
	predicateR store.Reader
}

func (val *validatorImpl) Validate(ctx context.Context, req *Request) error {
	// Load the predicate if it exists.
	if _, err := val.predicateR.GetBatch(req.CatalystIDs...); err != nil {
		return err
	}
	return nil
}
