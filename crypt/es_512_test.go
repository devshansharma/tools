package crypt_test

import (
	"testing"
	"time"

	"github.com/devshansharma/tools/crypt"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestJwtTokens(t *testing.T) {

	t.Run("token valid immediately", func(t *testing.T) {
		// Generate a private key
		privateKey, err := crypt.GenerateES512PrivateKey()
		if err != nil {
			t.Fatal(err)
		}

		// Create claims with iat and nbf set to now
		now := time.Now().Unix()
		claims := jwt.MapClaims{
			"iss": "example.com",
			"aud": "my-app",
			"exp": time.Now().Add(time.Hour * 24).Unix(),
			"sub": "test-user",
			"iat": now,
			"nbf": now,
		}

		accessToken, _, err := crypt.GenerateJWTTokens(privateKey, claims)
		if err != nil {
			t.Fatal(err)
		}

		publicKey := &privateKey.PublicKey
		_, err = crypt.ParseAndVerifyToken(accessToken, publicKey, "example.com", "my-app")
		assert.NoError(t, err, "Expected token to be valid, but got error:", err)
	})

	t.Run("token failed due to future nbf", func(t *testing.T) {
		// Generate a private key
		privateKey, err := crypt.GenerateES512PrivateKey()
		if err != nil {
			t.Fatal(err)
		}

		// Create claims with nbf in the future
		futureNBF := time.Now().Add(time.Hour).Unix()
		claims := jwt.MapClaims{
			"iss": "example.com",
			"aud": "my-app",
			"exp": time.Now().Add(time.Hour * 24).Unix(),
			"sub": "test-user",
			"iat": time.Now().Unix(),
			"nbf": futureNBF,
		}

		accessToken, _, err := crypt.GenerateJWTTokens(privateKey, claims)
		if err != nil {
			t.Fatal(err)
		}

		_, err = crypt.ParseAndVerifyToken(accessToken, &privateKey.PublicKey, "example.com", "my-app")
		assert.Equal(t, "failed to parse token: Token is not valid yet", err.Error())
	})

	t.Run("token failed due to expiration", func(t *testing.T) {
		// Generate a private key
		privateKey, err := crypt.GenerateES512PrivateKey()
		if err != nil {
			t.Fatal(err)
		}

		// Create some claims with expiration in the past
		claims := jwt.MapClaims{
			"iss": "example.com",
			"aud": "my-app",
			"exp": time.Now().Add(-time.Hour * 24).Unix(), // Expired one day ago
			"sub": "test-user",
			"iat": time.Now().Unix(),
			"nbf": time.Now().Unix(),
		}

		// Generate access token
		accessToken, _, err := crypt.GenerateJWTTokens(privateKey, claims)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the token - should fail due to expiration
		publicKey := &privateKey.PublicKey
		_, err = crypt.ParseAndVerifyToken(accessToken, publicKey, "example.com", "my-app")
		assert.Equal(t, "failed to parse token: Token is expired", err.Error())
	})

	t.Run("success case", func(t *testing.T) {
		// Generate a private key
		privateKey, err := crypt.GenerateES512PrivateKey()
		if err != nil {
			t.Fatal(err)
		}

		// Create some claims
		claims := jwt.MapClaims{
			"iss": "example.com",
			"aud": "my-app",
			"exp": time.Now().Add(time.Hour * 24).Unix(),
			"sub": "test-user",
			"iat": time.Now().Unix(),
			"nbf": time.Now().Unix(),
		}

		// Generate access and refresh tokens
		accessToken, refreshToken, err := crypt.GenerateJWTTokens(privateKey, claims)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the tokens
		publicKey := &privateKey.PublicKey
		_, err = crypt.ParseAndVerifyToken(accessToken, publicKey, "example.com", "my-app")
		if err != nil {
			t.Fatal(err)
		}

		_, err = crypt.ParseAndVerifyToken(refreshToken, publicKey, "example.com", "my-app")
		if err != nil {
			t.Fatal(err)
		}
	})
}
