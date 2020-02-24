package token

import (
	"time"

	"github.com/theonlyrob/vercer/webserver/pkg/identity/secret"

	"github.com/dgrijalva/jwt-go"
)

type jwtManagerImpl struct {
	secretManager secret.Manager
}

// Create a struct that will be encoded to a JWT.
type jwtClaims struct {
	UserID string `json:"userID"`
	jwt.StandardClaims
}

func (tm *jwtManagerImpl) Create(username string) (string, time.Time, error) {
	// Generate the claims for the response token.
	expirationTime := time.Now().Add(48 * time.Hour)
	claims := &jwtClaims{
		UserID: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create the token with the claims.
	var tokenString string
	var err error
	tm.secretManager.WhileRead(func(privateKey []byte) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err = token.SignedString(privateKey)
	})
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenString, expirationTime, nil
}

func (tm *jwtManagerImpl) Validate(tokenString string) (*Claims, error) {
	claims := &jwtClaims{}
	var tkn *jwt.Token
	var err error
	tm.secretManager.WhileRead(func(privateKey []byte) {
		tkn, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return privateKey, nil
		})
	})
	if err != nil {
		return nil, err
	}
	return &Claims{
		UserID:  claims.UserID,
		Valid:   tkn.Valid,
		Expires: time.Unix(claims.StandardClaims.ExpiresAt, 0),
	}, nil
}
