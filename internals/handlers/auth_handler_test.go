package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/zarinpy/abrnoc_weather/internals/models"
	"github.com/zarinpy/abrnoc_weather/internals/test"
	"golang.org/x/crypto/bcrypt"
)

func setupAuthTestRouter(t *testing.T) (*gin.Engine, *AuthHandler) {
	gin.SetMode(gin.TestMode)
	testDB := test.SetupTestDB()
	handler := NewAuthHandler(testDB, "test_secret_key")

	router := gin.New()
	router.POST("/auth/register", handler.Register)
	router.POST("/auth/login", handler.Login)

	return router, handler
}

func TestRegister_Success(t *testing.T) {
	router, handler := setupAuthTestRouter(t)
	defer test.CleanupTestDB(handler.db)
	test.ResetDB(handler.db)

	payload := RegisterRequest{
		Username: "testuser",
		Password: "password123",
	}
	jsonValue, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.User
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "testuser", response.Username)
	assert.Empty(t, response.Password) // Password should not be returned
	assert.NotEmpty(t, response.ID)
}

func TestRegister_DuplicateUsername(t *testing.T) {
	router, handler := setupAuthTestRouter(t)
	defer test.CleanupTestDB(handler.db)
	test.ResetDB(handler.db)

	// Create first user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 12)
	user := models.User{
		Username: "testuser",
		Password: string(hashedPassword),
	}
	handler.db.Create(&user)

	// Try to register with same username
	payload := RegisterRequest{
		Username: "testuser",
		Password: "password123",
	}
	jsonValue, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRegister_InvalidInput(t *testing.T) {
	router, handler := setupAuthTestRouter(t)
	defer test.CleanupTestDB(handler.db)
	test.ResetDB(handler.db)

	// Test with short password
	payload := RegisterRequest{
		Username: "testuser",
		Password: "short", // Less than 6 characters
	}
	jsonValue, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_MissingFields(t *testing.T) {
	router, handler := setupAuthTestRouter(t)
	defer test.CleanupTestDB(handler.db)
	test.ResetDB(handler.db)

	payload := map[string]string{
		"username": "testuser",
		// Missing password
	}
	jsonValue, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_Success(t *testing.T) {
	router, handler := setupAuthTestRouter(t)
	defer test.CleanupTestDB(handler.db)
	test.ResetDB(handler.db)

	// Create user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 12)
	user := models.User{
		Username: "testuser",
		Password: string(hashedPassword),
	}
	handler.db.Create(&user)

	payload := LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	jsonValue, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response LoginResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response.Token)
}

func TestLogin_InvalidUsername(t *testing.T) {
	router, handler := setupAuthTestRouter(t)
	defer test.CleanupTestDB(handler.db)
	test.ResetDB(handler.db)

	payload := LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}
	jsonValue, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_InvalidPassword(t *testing.T) {
	router, handler := setupAuthTestRouter(t)
	defer test.CleanupTestDB(handler.db)
	test.ResetDB(handler.db)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 12)
	user := models.User{
		Username: "testuser",
		Password: string(hashedPassword),
	}
	handler.db.Create(&user)

	payload := LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}
	jsonValue, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_MissingFields(t *testing.T) {
	router, handler := setupAuthTestRouter(t)
	defer test.CleanupTestDB(handler.db)
	test.ResetDB(handler.db)

	payload := map[string]string{
		"username": "testuser",
		// Missing password
	}
	jsonValue, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
