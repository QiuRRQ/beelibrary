package mytoken

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtCredential struct {
	TokenSecret         string
	ExpiredToken        int
	RefreshTokenSecret  string
	ExpiredRefreshToken int
}

type CustomClaims struct {
	jwt.StandardClaims
	Email   string `json:"email"`
	Session string `json:"session"`
}

func GetToken(session, userid, id string) (string, string, error) {
	expirationTime := time.Now().Add(time.Duration(0) * time.Hour).Unix()
	unixTimeUtc := time.Unix(expirationTime, 0)
	unitTimeInRFC3339 := unixTimeUtc.UTC().Format(time.RFC3339)

	claims := &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			Id:        id,
			Issuer:    userid,
		},
		Session: session,
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := rawToken.SignedString([]byte(""))

	return token, unitTimeInRFC3339, err
}

func GetRefreshToken(session, userid, id string) (string, string, error) {
	expirationTime := time.Now().Add(time.Duration(0) * time.Hour).Unix()
	unixTimeUtc := time.Unix(expirationTime, 0)
	unitTimeInRFC3339 := unixTimeUtc.UTC().Format(time.RFC3339)

	claims := &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			Id:        id,
			Issuer:    userid,
		},
		Session: session,
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := rawToken.SignedString([]byte(""))

	return token, unitTimeInRFC3339, err
}
