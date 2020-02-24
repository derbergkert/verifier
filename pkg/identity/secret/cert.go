package secret

var key = []byte(``)

type managerImpl struct{}

func (m *managerImpl) WhileRead(toRun func([]byte)) {
	toRun(key)
}
