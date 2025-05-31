package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"sso/internal/domain/entity"
	"time"
)

func NewToken(user entity.User, app entity.App, tokenTTL time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID.String()
	claims["username"] = user.Username
	claims["app_id"] = app.ID
	claims["exp"] = time.Now().Add(tokenTTL).Unix()

	tokenString, err := token.SignedString([]byte(app.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
