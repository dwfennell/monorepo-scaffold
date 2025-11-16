package testutil

import (
	"context"
	"fmt"
	"os"

	"github.com/dwfennell/monorepo-scaffold/internal/database"
)

// GetTestDatabaseURL returns the test database URL from environment
// or fails if not configured
func GetTestDatabaseURL() (string, error) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		return "", fmt.Errorf("TEST_DATABASE_URL environment variable is not set - run 'make test-db-setup' first")
	}
	return dbURL, nil
}

// MustGetTestDatabaseURL returns the test database URL or panics
func MustGetTestDatabaseURL() string {
	url, err := GetTestDatabaseURL()
	if err != nil {
		panic(err)
	}
	return url
}

// NewTestDB creates a new database connection for testing
func NewTestDB(ctx context.Context) (*database.DB, error) {
	dbURL, err := GetTestDatabaseURL()
	if err != nil {
		return nil, err
	}

	return database.NewDB(ctx, dbURL)
}
