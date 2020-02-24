package search

import (
	"context"
	"errors"
	"fmt"
	"time"

	identityContext "github.com/theonlyrob/vercer/webserver/pkg/identity/context"
	"github.com/theonlyrob/vercer/webserver/pkg/search"
	userStore "github.com/theonlyrob/vercer/webserver/services/user/store"
)

var (
	blackoutDuration = 20 * time.Minute
)

type authorizerImpl struct {
	userR userStore.Reader
}

func (auth *authorizerImpl) Authorize(ctx context.Context, req *search.Request) error {
	// Get the identity from the context.
	identity := identityContext.GetIdentity(ctx)

	// Load the user if it exists.
	var user *userStore.User
	if identity != nil {
		var err error
		user, err = auth.userR.Get(identity.ID)
		if err != nil {
			return err
		} else if user == nil {
			return errors.New("user missing for login identity")
		}
	}

	// If the user is not subscribed or logged in, add a filter for the blackout time.
	var allowedFilter *search.Filter
	if identity == nil {
		reqTime := time.Now()
		allowedFilter = &search.Filter{
			Base: fmt.Sprintf("createdTime:<%s", reqTime.Add(-1*blackoutDuration).Format("2006-01-02T15:04:05.000Z")),
		}
	}

	// Update the request filter with the allowed values.
	if req.Filter != nil && allowedFilter != nil {
		req.Filter = &search.Filter{
			And: []*search.Filter{
				req.Filter,
				allowedFilter,
			},
		}
	} else if allowedFilter != nil {
		req.Filter = allowedFilter
	}
	return nil
}
