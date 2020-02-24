package delete

import (
	"context"
	"errors"
	"github.com/theonlyrob/vercer/webserver/services/predicate/store"

	identityContext "github.com/theonlyrob/vercer/webserver/pkg/identity/context"
	userStore "github.com/theonlyrob/vercer/webserver/services/user/store"
)

type authorizerImpl struct {
	userR     userStore.Reader
	predicateR store.Reader
}

func (auth *authorizerImpl) Authorize(ctx context.Context, req *Request) error {
	// Get the identity from the context.
	identity := identityContext.GetIdentity(ctx)
	if identity == nil {
		return errors.New("permission denied")
	}

	// Load the user if it exists.
	user, err := auth.userR.Get(identity.ID)
	if err != nil {
		return err
	} else if user == nil {
		return errors.New("user missing")
	}

	// Load the predicate if it exists.
	predicate, err := auth.predicateR.Get(req.CatalystID)
	if err != nil {
		return err
	} else if predicate == nil {
		return errors.New("predicate missing")
	}

	// Fill in the predicate's user info to track who added it.
	if predicate.UserID != identity.ID {
		return errors.New("cannot impoersonate user")
	}

	return nil
}
