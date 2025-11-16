package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "mySecurePassword123"

	hash, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash, "Hash should not equal plain password")
}

func TestHashPassword_DifferentHashesForSamePassword(t *testing.T) {
	password := "samePassword"

	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2, "Different salts should produce different hashes")
}

func TestCheckPassword_ValidPassword(t *testing.T) {
	password := "correctPassword"
	hash, _ := HashPassword(password)

	result := CheckPassword(password, hash)

	assert.True(t, result, "Valid password should return true")
}

func TestCheckPassword_InvalidPassword(t *testing.T) {
	password := "correctPassword"
	hash, _ := HashPassword(password)

	result := CheckPassword("wrongPassword", hash)

	assert.False(t, result, "Invalid password should return false")
}

func TestCheckPassword_EmptyPassword(t *testing.T) {
	hash, _ := HashPassword("somePassword")

	result := CheckPassword("", hash)

	assert.False(t, result, "Empty password should return false")
}

func TestHashPassword_EmptyString(t *testing.T) {
	hash, err := HashPassword("")

	assert.NoError(t, err)
	assert.NotEmpty(t, hash, "Should be able to hash empty string")
}
