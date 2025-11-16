package auth

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	userID := 123
	email := "test@example.com"

	token, err := GenerateToken(userID, email)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateToken_MissingSecret(t *testing.T) {
	os.Unsetenv("JWT_SECRET")

	token, err := GenerateToken(1, "test@example.com")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "JWT_SECRET not set", err.Error())
}

func TestValidateToken_ValidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	userID := 456
	email := "validate@example.com"

	token, _ := GenerateToken(userID, email)
	claims, err := ValidateToken(token)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	invalidToken := "invalid.token.string"

	claims, err := ValidateToken(invalidToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_TamperedToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	token, _ := GenerateToken(1, "test@example.com")

	// Tamper with token
	tamperedToken := token + "tampered"

	claims, err := ValidateToken(tamperedToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_WrongSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "original-secret")
	token, _ := GenerateToken(1, "test@example.com")

	// Change secret
	os.Setenv("JWT_SECRET", "different-secret")
	defer os.Unsetenv("JWT_SECRET")

	claims, err := ValidateToken(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// Create an expired token manually
	claims := Claims{
		UserID: 1,
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret-key"))

	validatedClaims, err := ValidateToken(tokenString)

	assert.Error(t, err)
	assert.Nil(t, validatedClaims)
}

func TestValidateToken_MissingSecret(t *testing.T) {
	os.Unsetenv("JWT_SECRET")

	claims, err := ValidateToken("some.token.string")

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, "JWT_SECRET not set", err.Error())
}
