package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)

	if err != nil {
		return "", err
	}

	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {

	match, err := argon2id.ComparePasswordAndHash(password, hash)

	if err != nil {
		return false, err
	}

	if match == false {
		return false, nil
	}

	return true, nil

}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	mySigningKey := []byte(tokenSecret)

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", nil
	}

	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	} else if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {

		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			return uuid.Nil, errors.New("Token is expired")
		}

		userId, err := claims.GetSubject()

		if err != nil {
			return uuid.Nil, err
		}

		parsedUserId, err := uuid.Parse(userId)

		if err != nil {
			return uuid.Nil, err
		}

		return parsedUserId, nil

	} else {
		return uuid.Nil, errors.New("Invalid token")
	}

}

func GetBearerToken(headers http.Header) (string, error) {
	authString := headers.Get("Authorization")

	if authString == "" {
		return "", errors.New("No Authorization header set")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authString, prefix) {
		return "", errors.New("improper Authorization format")
	}

	tokenString := strings.TrimPrefix(authString, prefix)

	if tokenString == "" {
		return "", errors.New("No token sent")
	}

	return tokenString, nil
}
