package interceptor

import (
	"sync"

	identityInterceptor "github.com/theonlyrob/vercer/webserver/services/identity/interceptor"
)

var (
	once              sync.Once
	globalInterceptor Interceptor
)

func Singleton() Interceptor {
	once.Do(func() {
		globalInterceptor = &globalInterceptorImpl{
			interceptors: []Interceptor{
				identityInterceptor.Singleton(),
			},
		}
	})
	return globalInterceptor
}
