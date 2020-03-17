package handler

import (
	"github.com/rs/cors"
	"net/http"
	"sync"

	globalinterceptor "github.com/theonlyrob/vercer/webserver/cmd/interceptor"
	livenessService "github.com/theonlyrob/vercer/webserver/services/liveness"

	// Add more APIs you want to register here.
)

var (
	once sync.Once

	globalHandlerInstance *globalHandlerImpl
)

func Singleton() http.Handler {
	once.Do(func() {
		mux := http.NewServeMux()

		// Serve the API
		registerAPI(mux)

		// Create global handler which will use the interceptor to add user identity.
		globalHandlerInstance = &globalHandlerImpl{
			interceptor: globalinterceptor.Singleton(),
			mux:         mux,
		}
	})

	return globalHandlerInstance
}

// Static helper functions.
///////////////////////////

func registerAPI(mux *http.ServeMux) {
	// Serve the API
	apiMux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMux))

	// Attach individual APIs.
	livenessMux := http.NewServeMux()
	apiMux.Handle("/liveness/", http.StripPrefix("/liveness", livenessMux))
	livenessService.Register(livenessMux)
}