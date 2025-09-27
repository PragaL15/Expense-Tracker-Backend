package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/PragaL15/Expense-Tracker/internal/database"
	"github.com/PragaL15/Expense-Tracker/internal/models"
)

var validate5= validator.New()

// ================= Add Investment =================
func AddInvestment(c *fiber.Ctx) error {
	uidStr := c.Locals("user_id").(string)
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	var body models.Investment
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Assign server-side fields
	body.UserID = uid
	body.InvestmentID = uuid.New() // Correct type
	body.CreatedAt = time.Now()

	// Validate required fields
	if err := validate.Struct(body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := database.DB.Create(&body).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(body)
}

// ================= Get All Investments =================
func GetInvestments(c *fiber.Ctx) error {
	uidStr := c.Locals("user_id").(string)
	uid, _ := uuid.Parse(uidStr)

	var investments []models.Investment
	if err := database.DB.Where("user_id = ?", uid).Order("created_at desc").Find(&investments).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(investments)
}

// ================= Get Single Investment =================
func GetInvestment(c *fiber.Ctx) error {
	uidStr := c.Locals("user_id").(string)
	uid, _ := uuid.Parse(uidStr)
	idStr := c.Params("id")
	invID, err := uuid.Parse(idStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid investment ID")
	}

	var inv models.Investment
	if err := database.DB.Where("investment_id = ? AND user_id = ?", invID, uid).First(&inv).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Investment not found")
	}

	return c.JSON(inv)
}

// ================= Update Investment =================
func UpdateInvestment(c *fiber.Ctx) error {
	uidStr := c.Locals("user_id").(string)
	uid, _ := uuid.Parse(uidStr)
	idStr := c.Params("id")
	invID, err := uuid.Parse(idStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid investment ID")
	}

	var body models.Investment
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	var inv models.Investment
	if err := database.DB.Where("investment_id = ? AND user_id = ?", invID, uid).First(&inv).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Investment not found")
	}

	// Only update allowed fields
	updates := map[string]interface{}{
		"type":            body.Type,
		"amount_invested": body.AmountInvested,
		"current_value":   body.CurrentValue,
		"date_invested":   body.DateInvested,
		"reminder_date":   body.ReminderDate,
		"notes":           body.Notes,
	}

	if err := database.DB.Model(&inv).Updates(updates).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(inv)
}

// ================= Delete Investment =================
func DeleteInvestment(c *fiber.Ctx) error {
	uidStr := c.Locals("user_id").(string)
	uid, _ := uuid.Parse(uidStr)
	idStr := c.Params("id")
	invID, err := uuid.Parse(idStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid investment ID")
	}

	if err := database.DB.Where("investment_id = ? AND user_id = ?", invID, uid).Delete(&models.Investment{}).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
