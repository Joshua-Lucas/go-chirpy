package auth_test

import (
	"testing"

	"github.com/Joshua-Lucas/go-chirpy/internal/auth"
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
