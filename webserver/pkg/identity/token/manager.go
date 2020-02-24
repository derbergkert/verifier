package token

import "time"

type Claims struct {
	UserID string
	Valid    bool
	Expires  time.Time
}

type Manager interface {
	Create(username string) (token string, expTime time.Time, err error)
	Validate(token string) (claims *Claims, err error)
}
