package auth_test

import (
	"testing"
	"time"

	"chirpy/auth"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "supersecret"

	hashed, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hashed == "" {
		t.Fatalf("HashPassword returned empty hash")
	}

	if hashed == password {
		t.Fatalf("HashPassword returned unhashed password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mypassword"

	hashed, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	ok, err := auth.CheckPasswordHash(password, hashed)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}

	if !ok {
		t.Fatalf("CheckPasswordHash should return true for correct password")
	}

	// Wrong password should fail
	ok, err = auth.CheckPasswordHash("wrongpassword", hashed)
	if err == nil && ok {
		t.Fatalf("CheckPasswordHash should fail for incorrect password")
	}
}

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "testsecret"
	expiresIn := time.Hour

	tokenString, err := auth.MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	if tokenString == "" {
		t.Fatalf("MakeJWT returned empty token")
	}

	// Parse token to inspect claims
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		t.Fatalf("Token is invalid or claims type mismatch")
	}

	if claims.Issuer != "chirpy-access" {
		t.Fatalf("Issuer mismatch: got %s", claims.Issuer)
	}

	if claims.Subject != userID.String() {
		t.Fatalf("Subject mismatch: got %s", claims.Subject)
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "mysecret"
	expiresIn := time.Hour

	tokenString, err := auth.MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	parsedID, err := auth.ValidateJWT(tokenString, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned error: %v", err)
	}

	if parsedID != userID {
		t.Fatalf("ValidateJWT returned wrong user ID: got %s, want %s", parsedID, userID)
	}

	// Wrong secret should fail
	_, err = auth.ValidateJWT(tokenString, "wrongsecret")
	if err == nil {
		t.Fatalf("ValidateJWT should fail with wrong secret")
	}
}
