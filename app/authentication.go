package app

import (
	"fmt"
	"github.com/alpha-omega-corp/bunapp-api/app/httputils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/uptrace/bunrouter"
	"net/http"
	"os"
)

func ParseToken(ts string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		signingMethod, isValid := token.Method.(*jwt.SigningMethodHMAC)

		if isValid {
			return nil, fmt.Errorf("unexpected signing method: %v", signingMethod)
		}

		return []byte(secret), nil
	})
}

func GetValidTokenFromReq(w http.ResponseWriter, req bunrouter.Request) (*jwt.Token, error) {
	ts := req.Header.Get(tokenHeader)
	token, err := ParseToken(ts)
	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, httputils.From(err, true)
	}

	return token, nil
}
