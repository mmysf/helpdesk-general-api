package helpers

import (
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

func GetJWTSecretKeyCustomer() string {
	return os.Getenv("JWT_SECRET_KEY_USER")
}

func GetJWTSecretKeyAgent() string {
	return os.Getenv("JWT_SECRET_KEY_AGENT")
}

func GetJWTSecretKeySuperuser() string {
	return os.Getenv("JWT_SECRET_KEY_SUPERUSER")
}

func GetJWTSecretKeySuperadmin() string {
	return os.Getenv("JWT_SECRET_KEY_SUPERADMIN")
}

func GetJWTTTL() int {
	ttl, _ := strconv.Atoi(os.Getenv("JWT_TTL"))
	return ttl
}

func GenerateJWTTokenCustomer(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(GetJWTSecretKeyCustomer()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateJWTTokenAgent(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(GetJWTSecretKeyAgent()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateJWTTokenSuperuser(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(GetJWTSecretKeySuperuser()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateJWTTokenSuperadmin(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(GetJWTSecretKeySuperadmin()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
