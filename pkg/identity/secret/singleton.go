package secret

import (
	"sync"
)

var (
	once          sync.Once
	secretManager Manager
)

func Singleton() Manager {
	once.Do(func() {
		secretManager = &managerImpl{}
	})
	return secretManager
}
