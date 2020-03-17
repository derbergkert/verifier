package ping

import (
	"net/http"
	"sync"
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
	authorizer = NewAuthorizer(store.Singleton(), userStore.Singleton())
	validator = NewValidator()
	handler = NewHandler(
		authorizer,
		validator,
	)
}
