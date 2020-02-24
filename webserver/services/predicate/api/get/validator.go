package get

import (
	"context"
	"github.com/theonlyrob/vercer/webserver/services/predicate/store"
)

type Validator interface {
	Validate(ctx context.Context, req *Request) error
}

func NewValidator(predicateR store.Reader) Validator {
	return &validatorImpl{
		predicateR: predicateR,
	}
}
