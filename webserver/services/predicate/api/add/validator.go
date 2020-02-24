package add

import (
	"context"
)

type Validator interface {
	Validate(ctx context.Context, req *Request) error
}

func NewValidator() Validator {
	return &validatorImpl{}
}
