package secret

type Manager interface {
	WhileRead(func([]byte))
}
