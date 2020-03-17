package ping

import (
	"context"
)

type validatorImpl struct {}

func (val *validatorImpl) Validate(ctx context.Context, req *Request) error {
	return nil
}
