package server

import (
	"sync"
)

var (
	once                 sync.Once
	globalServerInstance GlobalServer
)

func Singleton() GlobalServer {
	once.Do(func() {
		globalServerInstance = &globalServerImpl{}
	})
	return globalServerInstance
}
