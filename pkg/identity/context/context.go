package context

import (
	goContext "context"
	"github.com/theonlyrob/vercer/webserver/services/identity/store"
)

func WithIdentity(ctx goContext.Context, ident password.Identity) goContext.Context {
	allocatedIdent := ident
	return goContext.WithValue(ctx, identityContextKey{}, &allocatedIdent)
}

func GetIdentity(ctx goContext.Context) *password.Identity {
	ident, ok := ctx.Value(identityContextKey{}).(*password.Identity)
	if ok {
		return ident
	}
	return nil
}

type identityContextKey struct{}
