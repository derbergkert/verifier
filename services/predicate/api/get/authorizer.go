package get

import (
	"context"
	"github.com/theonlyrob/vercer/webserver/services/predicate/store"

	userStore "github.com/theonlyrob/vercer/webserver/services/user/store"
)

type Authorizer interface {
	Authorize(ctx context.Context, req *Request) error
}

func NewAuthorizer(predicateR store.Reader, userR userStore.Reader) Authorizer {
	return &authorizerImpl{
		predicateR: predicateR,
		userR: userR,
	}
}
