package add

import (
	"github.com/theonlyrob/vercer/webserver/services/predicate/notifier"
	"github.com/theonlyrob/vercer/webserver/services/predicate/store"
	"net/http"
	"sync"

	userStore "github.com/theonlyrob/vercer/webserver/services/user/store"
)

var (
	once sync.Once

	authorizer Authorizer
	validator  Validator
	handler    http.Handler
)

// Singletons.
//////////////

// SingletonHandler returns the singleton instance of the http.Handler.
func SingletonHandler() http.Handler {
	once.Do(initialize)
	return handler
}

// SingletonAuthorizer returns the singleton instance of the Authorizer.
func SingletonAuthorizer() Authorizer {
	once.Do(initialize)
	return authorizer
}

// SingletonValidator returns the singleton instance of the Validator.
func SingletonValidator() Validator {
	once.Do(initialize)
	return validator
}

// Initialization.
//////////////////

func initialize() {
	authorizer = NewAuthorizer(userStore.Singleton())
	validator = NewValidator()
	handler = NewHandler(
		authorizer,
		validator,
		store.Singleton(),
		notifier.Singleton(),
	)
}
