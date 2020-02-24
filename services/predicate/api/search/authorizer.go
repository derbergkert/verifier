package search

import (
	"context"

	"github.com/theonlyrob/vercer/webserver/pkg/search"
	userStore "github.com/theonlyrob/vercer/webserver/services/user/store"
)

type Authorizer interface {
	Authorize(ctx context.Context, req *search.Request) error
}

func NewAuthorizer(userR userStore.Reader) Authorizer {
	return &authorizerImpl{
		userR: userR,
	}
}
