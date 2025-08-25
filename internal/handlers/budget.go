package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"

	"github.com/PragaL15/Expense-Tracker/internal/models"
	"github.com/PragaL15/Expense-Tracker/internal/database"
)

type budgetReq struct {
	Period      string  `json:"period" validate:"required"` 
	TotalBudget float64 `json:"total_budget" validate:"required,gte=0"`
}

func UpsertBudget(c *fiber.Ctx) error {
	uid := c.Locals("user_id").(string)
	var body budgetReq
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := validator.New().Struct(body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var b models.Budget
	if err := database.DB.Where("user_id = ? AND period = ?", uid, body.Period).First(&b).Error; err != nil {
		b = models.Budget{UserID: uid, Period: body.Period, TotalBudget: body.TotalBudget, TotalSpent: 0}
		if err := database.DB.Create(&b).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	} else {
		b.TotalBudget = body.TotalBudget
		if err := database.DB.Save(&b).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(b)
}

func GetBudget(c *fiber.Ctx) error {
	uid := c.Locals("user_id").(string)
	period := c.Query("period")
	if period == "" {
		return fiber.NewError(fiber.StatusBadRequest, "period is required (YYYY-MM)")
	}
	var b models.Budget
	if err := database.DB.Where("user_id = ? AND period = ?", uid, period).First(&b).Error; err != nil {
		return fiber.ErrNotFound
	}
	return c.JSON(b)
}
