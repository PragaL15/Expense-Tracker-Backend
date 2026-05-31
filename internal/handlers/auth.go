package handlers

import (
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/PragaL15/Expense-Tracker/internal/models"
	"github.com/PragaL15/Expense-Tracker/internal/database"
)

var validate = validator.New()

// Request body for user registration
type registerReq struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// Request body for user login
type loginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Register handles new user signup
func Register(c *fiber.Ctx) error {
	var body registerReq
	// Parse JSON request body into struct
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Validate fields using go-playground/validator
	if err := validate.Struct(body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Create new user model
	u := models.User{
		Name:         body.Name,
		Email:        body.Email,
		PasswordHash: string(hash),
	}

	// Save user to database (fail if email already exists)
	if err := database.DB.Create(&u).Error; err != nil {
    return fiber.NewError(fiber.StatusBadRequest, err.Error())
}

	// Respond with created user (no password)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user_id": u.UserID,
		"name":    u.Name,
		"email":   u.Email,
	})
}

// Login authenticates user and returns a JWT token
func Login(c *fiber.Ctx) error {
	var body loginReq
	// Parse request body
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	// Validate input
	if err := validate.Struct(body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Find user by email
	var u models.User
	if err := database.DB.Where("email = ?", body.Email).First(&u).Error; err != nil {
		return fiber.ErrUnauthorized
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(body.Password)); err != nil {
		return fiber.ErrUnauthorized
	}

	// Create JWT token
	token, err := makeJWT(u.UserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"token": token})
}

// Profile returns authenticated user's profile
func Profile(c *fiber.Ctx) error {
	uid := c.Locals("user_id") // user_id is set by authentication middleware
	var u models.User
	if err := database.DB.First(&u, "user_id = ?", uid).Error; err != nil {
		return fiber.ErrNotFound
	}
	return c.JSON(fiber.Map{
		"user_id": u.UserID,
		"name":    u.Name,
		"email":   u.Email,
	})
}

func makeJWT(userID string) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	expHrs := time.Duration(envToInt("JWT_EXPIRES_HOURS", 720)) * time.Hour
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(expHrs).Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret)
}

func envToInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if dur, err := time.ParseDuration(v + "h"); err == nil {
			return int(dur.Hours())
		}
	}
	return def
}
