package handler

import "net/http"

func webappHandler() http.Handler {
	return http.FileServer(http.Dir("../webapp/build"))
}
