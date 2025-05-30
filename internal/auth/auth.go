package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const (
	// TokenTypeAccess -
	TokenTypeAccess TokenType = "chirpy-access"
)

// ErrNoAuthHeaderIncluded -
var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

// HashPassword -
func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// CheckPasswordHash -
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// MakeJWT - creates a new JWT token for the given user ID with the specified expiration duration.
func MakeJWT(
	userID uuid.UUID,
	tokenSecret string,
	expiresIn time.Duration,
) (string, error) {
	// Ensure the tokenSecret is not empty
	signingKey := []byte(tokenSecret)
	// Create the claims for the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	// Sign the token with the provided secret
	return token.SignedString(signingKey)
}

// ValidateJWT - validates the JWT token and returns the user ID if valid.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// claimsStruct is used to parse the JWT claims.
	claimsStruct := jwt.RegisteredClaims{}
	// Parse the token with the provided claims structure and secret.
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return uuid.Nil, err
	}

	// GetSubject retrieves the subject (user ID) from the token claims.
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	// Check if the token is valid and has the correct issuer
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	// Ensure the issuer is correct
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	// Parse the user ID from the subject
	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	// Get the Authorization header
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	// Split the header into parts
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}

// MakeRefreshToken makes a random 256 bit token
// encoded in hex
func MakeRefreshToken() (string, error) {
	token := make([]byte, 32) // 256 bits = 32 bytes
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

// GetAPIKey extracts the API key from the request headers.
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	// Split the header into parts
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) != 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}
