package crypt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// Function to parse and verify the JWT
func ParseAndVerifyToken(tokenString string, publicKey *ecdsa.PublicKey, expectedIssuer, expectedAudience string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check issuer
		if claims["iss"] != expectedIssuer {
			return nil, fmt.Errorf("invalid issuer: %v", claims["iss"])
		}

		// Check audience
		if claims["aud"] != expectedAudience {
			return nil, fmt.Errorf("invalid audience: %v", claims["aud"])
		}

		// Check not before (nbf)
		if nbf, ok := claims["nbf"].(float64); ok {
			if time.Now().Unix() < int64(nbf) {
				return nil, fmt.Errorf("token not yet valid")
			}
		}

		// Check expiration
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, fmt.Errorf("token has expired")
			}
		} else {
			return nil, fmt.Errorf("invalid expiration claim")
		}

		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// Function to generate JWT access and refresh tokens
func GenerateJWTTokens(privateKey *ecdsa.PrivateKey, claims jwt.MapClaims) (string, string, error) {
	accessToken, err := CreateAccessToken(privateKey, claims)
	if err != nil {
		return "", "", fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := CreateRefreshToken(privateKey, claims)
	if err != nil {
		return "", "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// CreateAccessToken for creating access token
func CreateAccessToken(privateKey *ecdsa.PrivateKey, claims jwt.MapClaims) (string, error) {
	// Create the access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodES512, claims)
	signedAccessToken, err := accessToken.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return signedAccessToken, nil
}

// CreateRefreshToken for creating refresh token
func CreateRefreshToken(privateKey *ecdsa.PrivateKey, claims jwt.Claims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodES512, claims)
	signedRefreshToken, err := refreshToken.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return signedRefreshToken, nil
}

// Function to generate an ES512 private key
func GenerateES512PrivateKey() (*ecdsa.PrivateKey, error) {
	// Use the P-521 curve for ES512 (secp521r1)
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	return privateKey, nil
}

func ReadPrivateKeyFromPEM(fileName string) (*ecdsa.PrivateKey, error) {
	// Read the PEM file
	pemData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read PEM file: %w", err)
	}

	// Decode the PEM block
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing EC private key")
	}

	// Parse the EC private key
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse EC private key: %w", err)
	}

	return privateKey, nil
}

// Function to save the private key to a PEM file
func SavePrivateKeyToPEM(privateKey *ecdsa.PrivateKey, fileName string) error {
	// Marshal the private key to DER format
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}

	// Create a PEM block with the DER-encoded private key
	pemBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	// Write the PEM block to a file
	privateKeyFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer privateKeyFile.Close()

	err = pem.Encode(privateKeyFile, pemBlock)
	if err != nil {
		return fmt.Errorf("failed to encode private key to PEM: %w", err)
	}

	return nil
}
