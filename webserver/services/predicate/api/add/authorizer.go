package add

import (
	"context"

	userStore "github.com/theonlyrob/vercer/webserver/services/user/store"
)

type Authorizer interface {
	Authorize(ctx context.Context, req *Request) error
}

func NewAuthorizer(userR userStore.Reader) Authorizer {
	return &authorizerImpl{
		userR: userR,
	}
}
