package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type JWTClaim struct {
	ID int
	jwt.RegisteredClaims
}

func CreateToken(id int, secret string) (string, error) {
	secretByte := []byte(secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTClaim{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	})

	tokenStr, err := token.SignedString(secretByte)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func ClaimToken(token string, secret string) (*JWTClaim, error) {
	secretByte := []byte(secret)
	jwtToken, err := jwt.ParseWithClaims(
		token,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) { return secretByte, nil },
	)

	if err != nil {
		return nil, err
	}

	if jwtToken.Method != jwt.SigningMethodHS256 {
		return nil, errors.New("Invalid token")
	}

	claim, ok := jwtToken.Claims.(*JWTClaim)
	if !ok {
		return nil, errors.New("Invalid token")
	}
	return claim, nil
}
