package token

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidToken = errors.New("invalid token has shared")
	ErrExpiredToken = errors.New("token expired")
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretkey string
}

func NewJWTMaker(secretkey string) (Maker, error) {
	if len(secretkey) < minSecretKeySize {
		return nil, errors.New("minimum length of secret key should be 32")
	}
	return &JWTMaker{secretkey}, nil
}

func (jm *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString([]byte(jm.secretkey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (jm *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(jm.secretkey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

func (jm *JWTMaker) CreateRefreshToken(username string, duration time.Duration) (string, error) {
	return "", nil
}
