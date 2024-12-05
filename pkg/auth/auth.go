package auth

import (
	"final/common/consts"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Role int

var JWTSecretKey string = os.Getenv(consts.JWTSecretKey)

const (
	RoleAdmin Role = 1 << iota
	RoleUser
)
const RoleNone Role = 0

var rolePermissions = map[Role][]string{
	RoleAdmin: {"create", "read", "update", "delete"},
	RoleUser:  {"read"},
}

func CreateJWT(userID string, role Role, expirationTime time.Duration) (string, error) {
	claims := jwt.MapClaims{
		consts.ClaimsUserID:    userID,
		consts.ClaimsRole:      role,
		consts.ClaimsExpiresAt: time.Now().Add(expirationTime).Unix(),
		consts.ClaimsIssuedAt:  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(JWTSecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(JWTSecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}
