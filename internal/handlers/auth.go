package handlers

import (
	"log"
	"os"
	"time"

	"github.com/PragaL15/Expense-Tracker/internal/database"
	"github.com/PragaL15/Expense-Tracker/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

type registerReq struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type loginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func Register(c *fiber.Ctx) error {
	log.Println("REGISTER HIT")

	var body registerReq

	if err := c.BodyParser(&body); err != nil {
		log.Println("BodyParser Error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := validate.Struct(body); err != nil {
		log.Println("Validation Error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var existing models.User
	if err := database.DB.Where("email = ?", body.Email).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "email already exists",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Hash Error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to hash password",
		})
	}

	user := models.User{
		Name:         body.Name,
		Email:        body.Email,
		PasswordHash: string(hash),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Println("DB Create Error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user_id": user.UserID,
		"name":    user.Name,
		"email":   user.Email,
	})
}

func Login(c *fiber.Ctx) error {
	log.Println("LOGIN HIT")

	var body loginReq

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var user models.User

	if err := database.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(body.Password),
	); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid password",
		})
	}

	token, err := makeJWT(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

func makeJWT(userID string) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		log.Println("JWT_SECRET missing")
	}

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(30 * 24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
