package auth_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/Joshua-Lucas/go-chirpy/internal/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestCreatingAHash(t *testing.T) {
	password := "hunter2"

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash == "" {
		t.Errorf("expected non-empty hash")
	}
	if hash == password {
		t.Errorf("hash should not equal the original password")
	}

	// Hashing the same password twice should produce different hashes
	hash2, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error on second hash: %v", err)
	}
	if hash == hash2 {
		t.Errorf("two hashes of the same password should not be identical")
	}
}

func TestCheckingHash(t *testing.T) {
	cases := map[string]struct {
		password      string
		passwordCheck string
		match         bool
	}{
		"correct password matches": {
			password:      "hunter2",
			passwordCheck: "hunter2",
			match:         true,
		},
		"wrong password does not match": {
			password:      "hunter2",
			passwordCheck: "wrongpassword",
			match:         false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			hash, err := auth.HashPassword(tc.password)
			if err != nil {
				t.Fatalf("could not hash password: %v", err)
			}

			match, err := auth.CheckPasswordHash(tc.passwordCheck, hash)
			if err != nil && tc.match {
				t.Fatalf("unexpected error: %v", err)
			}
			if match != tc.match {
				t.Errorf("expected match=%v, got match=%v", tc.match, match)
			}
		})
	}
}

func TestJWTParse(t *testing.T) {

	SECERT := "test-secret"
	userID := uuid.New()

	validToken, _ := auth.MakeJWT(userID, SECERT, time.Hour)
	expiredToken, _ := auth.MakeJWT(userID, SECERT, -time.Minute)

	invalidUUIDToken := func() string {
		claims := jwt.RegisteredClaims{Subject: "not-a-uuid"}
		tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, _ := tok.SignedString([]byte(SECERT))
		return signed
	}()

	tests := []struct {
		name        string
		token       string
		secret      string
		nowFunc     func() time.Time
		expectError bool
		expectedID  uuid.UUID
	}{
		{
			name:        "valid token",
			token:       validToken,
			secret:      SECERT,
			expectError: false,
			expectedID:  userID,
		},
		{
			name:        "expired token",
			token:       expiredToken,
			secret:      SECERT,
			expectError: true,
		},
		{
			name:        "wrong secret",
			token:       validToken,
			secret:      "wrong-secret",
			expectError: true,
		},
		{
			name:        "invalid token string",
			token:       "not-a-token",
			secret:      SECERT,
			expectError: true,
		},
		{
			name:        "invalid UUID in subject",
			token:       invalidUUIDToken,
			secret:      SECERT,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := auth.ValidateJWT(tt.token, tt.secret)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if id != tt.expectedID {
				t.Errorf("expected UUID %v, got %v", tt.expectedID, id)
			}
		})
	}

}

func TestGetBearToken(t *testing.T) {

	tests := []struct {
		name        string
		headers     http.Header
		expect      string
		expectError bool
	}{
		{
			name: "valid token",
			headers: http.Header{
				"Authorization": []string{"Bearer abc123"},
			},
			expect: "abc123",
		},
		{
			name:        "missing header",
			headers:     http.Header{},
			expectError: true,
		},
		{
			name: "empty bearer",
			headers: http.Header{
				"Authorization": []string{"Bearer "},
			},
			expectError: true,
		},
		{
			name: "wrong format",
			headers: http.Header{
				"Authorization": []string{"abc123"},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := auth.GetBearerToken(tt.headers)

			if (err != nil) != tt.expectError {
				t.Fatalf("expected error=%v, got err=%v", tt.expectError, err)
			}

			if got != tt.expect {
				t.Errorf("expected %q, got %q", tt.expect, got)
			}
		})
	}

}
