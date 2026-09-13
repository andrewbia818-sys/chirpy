package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// Hash the password using the argon2id.CreateHash function.
func HashPassword(password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}
func CheckPasswordHash(password, hash string) (bool, error) {
	matches, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("failed to check password hash: %w", err)
	}

	return matches, nil
}

// create a JWT token for the given user ID with the specified expiration time.
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now().UTC()

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tokenSecret))
}

// ValidateJWT validates the given JWT token string using the provided secret and
// returns the user ID if valid.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return uuid.Nil, fmt.Errorf("failed to parse user ID from token: %w", err)
		}
		return userID, nil
	}

	return uuid.Nil, fmt.Errorf("invalid token")
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid authorization header format")
	}
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

// Make refresh token using rand.Read to generate 32 bytes (256 bits)then use
// hex.EncodeToString to convert to a hex string
func MakeRefreshToken() string {
	refresh_token_string := make([]byte, 32)
	_, err := rand.Read(refresh_token_string)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(refresh_token_string)
}

// Extract the api key from the Authorization header, which is expected to be in this format:
// Authorization: ApiKey THE_KEY_HERE
//
//	strip out the ApiKey part and the whitespace and return just the key.
/*
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	const prefix = "ApiKey "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", fmt.Errorf("invalid authorization header format")
	}

	key := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
	if key == "" {
		return "", fmt.Errorf("api key is empty")
	}

	return key, nil
}
*/
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	const prefix = "ApiKey "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", fmt.Errorf("invalid authorization header format")
	}

	key := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
	if key == "" {
		return "", fmt.Errorf("api key is empty")
	}

	return key, nil
}

/* OLD Below
func GetAPIKey(headers http.Header) (string, error) {
	// Boot.dev Polka webhook format
	polkaSig := headers.Get("X-Polka-Signature")
	if polkaSig != "" {
		return strings.TrimSpace(polkaSig), nil
	}

	// Authorization: ApiKey <key>
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	const prefix = "ApiKey "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", fmt.Errorf("invalid authorization header format")
	}

	key := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
	if key == "" {
		return "", fmt.Errorf("api key is empty")
	}

	return key, nil
}
*/
