package server

import "net/http"

type GlobalServer interface {
	Run(http.Handler) error
}
