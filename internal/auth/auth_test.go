package auth_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/andrewbia818-sys/chirpy/internal/auth"
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

func TestGetBearerToken_Success(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Bearer abc123")

	token, err := auth.GetBearerToken(h)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token != "abc123" {
		t.Fatalf("expected token 'abc123', got %q", token)
	}
}

func TestGetBearerToken_MissingHeader(t *testing.T) {
	h := http.Header{} // no Authorization header

	_, err := auth.GetBearerToken(h)
	if err == nil {
		t.Fatalf("expected error for missing header, got nil")
	}
}

func TestGetBearerToken_InvalidFormat(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Token abc123") // wrong prefix

	_, err := auth.GetBearerToken(h)
	if err == nil {
		t.Fatalf("expected error for invalid format, got nil")
	}
}

func TestGetBearerToken_EmptyBearerValue(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Bearer ") // valid prefix, empty token

	token, err := auth.GetBearerToken(h)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token != "" {
		t.Fatalf("expected empty token, got %q", token)
	}
}

func TestGetAPIKey(t *testing.T) {
	tests := []struct {
		name    string
		header  http.Header
		wantKey string
		wantErr bool
	}{
		{
			name: "valid ApiKey header",
			header: http.Header{
				"Authorization": []string{"ApiKey abc123"},
			},
			wantKey: "abc123",
			wantErr: false,
		},
		{
			name: "valid with extra whitespace",
			header: http.Header{
				"Authorization": []string{"ApiKey     xyz789   "},
			},
			wantKey: "xyz789",
			wantErr: false,
		},
		{
			name:    "missing Authorization header",
			header:  http.Header{},
			wantErr: true,
		},
		{
			name: "wrong prefix",
			header: http.Header{
				"Authorization": []string{"Bearer abc123"},
			},
			wantErr: true,
		},
		{
			name: "empty key",
			header: http.Header{
				"Authorization": []string{"ApiKey   "},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := auth.GetAPIKey(tt.header)

			if tt.wantErr && err == nil {
				t.Fatalf("expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if key != tt.wantKey {
				t.Fatalf("expected key %q, got %q", tt.wantKey, key)
			}
		})
	}
}
