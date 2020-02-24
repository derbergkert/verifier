package add

import (
	"context"
	"errors"

	identityContext "github.com/theonlyrob/vercer/webserver/pkg/identity/context"
	userStore "github.com/theonlyrob/vercer/webserver/services/user/store"
)

type authorizerImpl struct {
	userR userStore.Reader
}

func (auth *authorizerImpl) Authorize(ctx context.Context, req *Request) error {
	// Get the identity from the context.
	identity := identityContext.GetIdentity(ctx)
	if identity == nil {
		return errors.New("permission denied")
	}

	// Fill in the predicate's user info to track who added it.
	req.Catalyst.UserID = identity.ID

	// Load the user if it exists, and add it to the context.
	user, err := auth.userR.Get(identity.ID)
	if err != nil {
		return err
	} else if user == nil {
		return errors.New("user missing")
	}
	return nil
}
