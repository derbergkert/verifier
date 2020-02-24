package token

import (
	"sync"

	"github.com/theonlyrob/vercer/webserver/pkg/identity/secret"
)

var (
	once         sync.Once
	tokenManager Manager
)

func Singleton() Manager {
	once.Do(func() {
		tokenManager = &jwtManagerImpl{
			secretManager: secret.Singleton(),
		}
	})
	return tokenManager
}
