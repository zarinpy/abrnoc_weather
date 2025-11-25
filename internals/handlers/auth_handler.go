package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zarinpy/abrnoc_weather/internals/models"
	"github.com/zarinpy/abrnoc_weather/pkg/respond"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	db        *gorm.DB
	jwtSecret string
}

// NewAuthHandler constructs an AuthHandler.
func NewAuthHandler(db *gorm.DB, jwtSecret string) *AuthHandler {
	return &AuthHandler{db: db, jwtSecret: jwtSecret}
}

// LoginRequest represents the login request body.
type LoginRequest struct {
	Username string `json:"username" example:"testuser" binding:"required"`
	Password string `json:"password" example:"password123" binding:"required"`
}

// LoginResponse represents the login response.
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		respond.ValidationError(c, err.Error())
		return
	}

	var user models.User
	result := h.db.Where("username = ?", input.Username).First(&user)
	if result.Error != nil {
		respond.Error(c, http.StatusUnauthorized, "invalid_credentials", "Invalid username or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		respond.Error(c, http.StatusUnauthorized, "invalid_credentials", "Invalid username or password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"user_id":  user.ID.String(),
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		respond.Error(c, http.StatusInternalServerError, "token_generation_failed", "Failed to generate token")
		return
	}

	respond.JSON(c, http.StatusOK, LoginResponse{Token: tokenString})
}

// RegisterRequest represents the registration request body.
type RegisterRequest struct {
	Username string `json:"username" example:"testuser" binding:"required,min=3,max=50"`
	Password string `json:"password" example:"password123" binding:"required,min=6"`
}

// @Summary Register new user
// @Description Create a new user account with hashed password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 409 {object} map[string]string "Username already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		respond.ValidationError(c, err.Error())
		return
	}

	var existingUser models.User
	if err := h.db.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		respond.Error(c, http.StatusConflict, "username_exists", "Username already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		respond.Error(c, http.StatusInternalServerError, "password_hash_failed", "Failed to hash password")
		return
	}

	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
	}

	if err := h.db.Create(&user).Error; err != nil {
		respond.Error(c, http.StatusInternalServerError, "user_create_failed", "Failed to create user")
		return
	}

	user.Password = ""
	respond.JSON(c, http.StatusCreated, user)
}
