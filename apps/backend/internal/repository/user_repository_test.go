package repository

import (
	"context"
	"testing"

	"github.com/dwfennell/monorepo-scaffold/internal/database"
	"github.com/dwfennell/monorepo-scaffold/internal/models"
	"github.com/dwfennell/monorepo-scaffold/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserRepositoryTestSuite is an integration test suite that requires a running database
type UserRepositoryTestSuite struct {
	suite.Suite
	db   *database.DB
	repo *UserRepository
	ctx  context.Context
}

// SetupSuite runs once before all tests
func (suite *UserRepositoryTestSuite) SetupSuite() {
	var err error
	suite.ctx = context.Background()
	suite.db, err = testutil.NewTestDB(suite.ctx)
	suite.Require().NoError(err)

	suite.repo = NewUserRepository(suite.db)
}

// TearDownSuite runs once after all tests
func (suite *UserRepositoryTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

// SetupTest runs before each test - clean up test data
func (suite *UserRepositoryTestSuite) SetupTest() {
	// Clean up users table before each test
	_, err := suite.db.Pool.Exec(suite.ctx, "DELETE FROM users")
	suite.Require().NoError(err, "Failed to clean up test data")
}

func (suite *UserRepositoryTestSuite) TestCreate_Success() {
	user := &models.User{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
	}

	err := suite.repo.Create(suite.ctx, user)

	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), user.ID, "User ID should be set after creation")
	assert.NotZero(suite.T(), user.CreatedAt, "CreatedAt should be set")
	assert.NotZero(suite.T(), user.UpdatedAt, "UpdatedAt should be set")
}

func (suite *UserRepositoryTestSuite) TestCreate_DuplicateEmail() {
	user1 := &models.User{
		Email:        "duplicate@example.com",
		PasswordHash: "hash1",
		Name:         "User 1",
	}

	user2 := &models.User{
		Email:        "duplicate@example.com", // Same email
		PasswordHash: "hash2",
		Name:         "User 2",
	}

	err1 := suite.repo.Create(suite.ctx, user1)
	assert.NoError(suite.T(), err1)

	err2 := suite.repo.Create(suite.ctx, user2)
	assert.Error(suite.T(), err2, "Should fail on duplicate email")
}

func (suite *UserRepositoryTestSuite) TestGetByEmail_Found() {
	// Create test user
	original := &models.User{
		Email:        "find@example.com",
		PasswordHash: "hashedpass",
		Name:         "Find Me",
	}
	suite.repo.Create(suite.ctx, original)

	// Find by email
	found, err := suite.repo.GetByEmail(suite.ctx, "find@example.com")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), found)
	assert.Equal(suite.T(), original.Email, found.Email)
	assert.Equal(suite.T(), original.Name, found.Name)
	assert.Equal(suite.T(), original.PasswordHash, found.PasswordHash)
}

func (suite *UserRepositoryTestSuite) TestGetByEmail_NotFound() {
	found, err := suite.repo.GetByEmail(suite.ctx, "nonexistent@example.com")

	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), found, "Should return nil for non-existent email")
}

func (suite *UserRepositoryTestSuite) TestGetByID_Found() {
	// Create test user
	original := &models.User{
		Email:        "findbyid@example.com",
		PasswordHash: "hashedpass",
		Name:         "Find By ID",
	}
	suite.repo.Create(suite.ctx, original)

	// Find by ID
	found, err := suite.repo.GetByID(suite.ctx, original.ID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), found)
	assert.Equal(suite.T(), original.ID, found.ID)
	assert.Equal(suite.T(), original.Email, found.Email)
}

func (suite *UserRepositoryTestSuite) TestGetByID_NotFound() {
	found, err := suite.repo.GetByID(suite.ctx, 99999)

	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), found, "Should return nil for non-existent ID")
}

// Run the test suite
func TestUserRepositoryTestSuite(t *testing.T) {
	// Skip integration tests if SHORT flag is set
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	suite.Run(t, new(UserRepositoryTestSuite))
}
