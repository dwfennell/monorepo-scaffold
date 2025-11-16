package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dwfennell/monorepo-scaffold/internal/database"
	"github.com/dwfennell/monorepo-scaffold/internal/models"
	"github.com/dwfennell/monorepo-scaffold/internal/repository"
	"github.com/dwfennell/monorepo-scaffold/internal/testutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	db      *database.DB
	handler *AuthHandler
	router  *gin.Engine
	ctx     context.Context
}

func (suite *AuthHandlerTestSuite) SetupSuite() {
	// Set JWT secret for tests
	os.Setenv("JWT_SECRET", "test-secret-key")

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Connect to test database
	var err error
	suite.ctx = context.Background()
	suite.db, err = testutil.NewTestDB(suite.ctx)
	suite.Require().NoError(err)

	// Setup handler and router
	userRepo := repository.NewUserRepository(suite.db)
	suite.handler = NewAuthHandler(userRepo)

	suite.router = gin.New()
	suite.router.POST("/register", suite.handler.Register)
	suite.router.POST("/login", suite.handler.Login)
	suite.router.GET("/me", AuthMiddleware(), suite.handler.GetCurrentUser)
}

func (suite *AuthHandlerTestSuite) TearDownSuite() {
	os.Unsetenv("JWT_SECRET")
	if suite.db != nil {
		suite.db.Close()
	}
}

func (suite *AuthHandlerTestSuite) SetupTest() {
	// Clean up users table before each test
	_, err := suite.db.Pool.Exec(suite.ctx, "DELETE FROM users")
	suite.Require().NoError(err, "Failed to clean up test data")
}

func (suite *AuthHandlerTestSuite) TestRegister_Success() {
	reqBody := models.RegisterRequest{
		Email:    "newuser@example.com",
		Password: "password123",
		Name:     "New User",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response models.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), response.Token)
	assert.Equal(suite.T(), reqBody.Email, response.User.Email)
	assert.Equal(suite.T(), reqBody.Name, response.User.Name)
}

func (suite *AuthHandlerTestSuite) TestRegister_DuplicateEmail() {
	// Create first user
	reqBody1 := models.RegisterRequest{
		Email:    "duplicate@example.com",
		Password: "password123",
		Name:     "First User",
	}
	body1, _ := json.Marshal(reqBody1)
	req1 := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	suite.router.ServeHTTP(w1, req1)
	assert.Equal(suite.T(), http.StatusCreated, w1.Code)

	// Try to create second user with same email
	reqBody2 := models.RegisterRequest{
		Email:    "duplicate@example.com",
		Password: "different",
		Name:     "Second User",
	}
	body2, _ := json.Marshal(reqBody2)
	req2 := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	suite.router.ServeHTTP(w2, req2)

	assert.Equal(suite.T(), http.StatusConflict, w2.Code)
}

func (suite *AuthHandlerTestSuite) TestRegister_InvalidJSON() {
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *AuthHandlerTestSuite) TestLogin_Success() {
	// First register a user
	registerBody := models.RegisterRequest{
		Email:    "logintest@example.com",
		Password: "password123",
		Name:     "Login Test",
	}
	regJSON, _ := json.Marshal(registerBody)
	regReq := httptest.NewRequest("POST", "/register", bytes.NewBuffer(regJSON))
	regReq.Header.Set("Content-Type", "application/json")
	regW := httptest.NewRecorder()
	suite.router.ServeHTTP(regW, regReq)

	// Now login with same credentials
	loginBody := models.LoginRequest{
		Email:    "logintest@example.com",
		Password: "password123",
	}
	loginJSON, _ := json.Marshal(loginBody)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), response.Token)
	assert.Equal(suite.T(), loginBody.Email, response.User.Email)
}

func (suite *AuthHandlerTestSuite) TestLogin_InvalidPassword() {
	// Register user
	registerBody := models.RegisterRequest{
		Email:    "wrongpass@example.com",
		Password: "correctpassword",
		Name:     "Wrong Pass Test",
	}
	regJSON, _ := json.Marshal(registerBody)
	regReq := httptest.NewRequest("POST", "/register", bytes.NewBuffer(regJSON))
	regReq.Header.Set("Content-Type", "application/json")
	regW := httptest.NewRecorder()
	suite.router.ServeHTTP(regW, regReq)

	// Try to login with wrong password
	loginBody := models.LoginRequest{
		Email:    "wrongpass@example.com",
		Password: "wrongpassword",
	}
	loginJSON, _ := json.Marshal(loginBody)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *AuthHandlerTestSuite) TestLogin_NonExistentUser() {
	loginBody := models.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}
	loginJSON, _ := json.Marshal(loginBody)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *AuthHandlerTestSuite) TestGetCurrentUser_ValidToken() {
	// Register and get token
	registerBody := models.RegisterRequest{
		Email:    "tokentest@example.com",
		Password: "password123",
		Name:     "Token Test",
	}
	regJSON, _ := json.Marshal(registerBody)
	regReq := httptest.NewRequest("POST", "/register", bytes.NewBuffer(regJSON))
	regReq.Header.Set("Content-Type", "application/json")
	regW := httptest.NewRecorder()
	suite.router.ServeHTTP(regW, regReq)

	var authResponse models.AuthResponse
	json.Unmarshal(regW.Body.Bytes(), &authResponse)
	token := authResponse.Token

	// Use token to get current user
	req := httptest.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var user models.User
	err := json.Unmarshal(w.Body.Bytes(), &user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), registerBody.Email, user.Email)
}

func (suite *AuthHandlerTestSuite) TestGetCurrentUser_MissingToken() {
	req := httptest.NewRequest("GET", "/me", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *AuthHandlerTestSuite) TestGetCurrentUser_InvalidToken() {
	req := httptest.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func TestAuthHandlerTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	suite.Run(t, new(AuthHandlerTestSuite))
}
