package search

import (
	"github.com/theonlyrob/vercer/webserver/services/predicate/index"
	"github.com/theonlyrob/vercer/webserver/services/predicate/store"
	"net/http"
	"sync"

	userStore "github.com/theonlyrob/vercer/webserver/services/user/store"
)

var (
	once sync.Once

	authorizer Authorizer
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

// Initialization.
//////////////////

func initialize() {
	authorizer = NewAuthorizer(userStore.Singleton())
	handler = NewHandler(
		authorizer,
		index.Singleton(),
		store.Singleton(),
	)
}
