package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/google/uuid"
	"time"
	"github.com/PragaL15/Expense-Tracker/internal/models"
)

type InvestmentHandler struct {
	DB *sqlx.DB
}

// Add Investment
func (h *InvestmentHandler) AddInvestment(c *fiber.Ctx) error {
	var inv models.Investment
	if err := c.BodyParser(&inv); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	inv.InvestmentID = uuid.New()
	inv.CreatedAt = time.Now()

	_, err := h.DB.NamedExec(`
		INSERT INTO investments 
		(investment_id, user_id, type, amount_invested, current_value, date_invested, reminder_date, notes, created_at)
		VALUES (:investment_id, :user_id, :type, :amount_invested, :current_value, :date_invested, :reminder_date, :notes, :created_at)
	`, inv)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB insert failed", "details": err.Error()})
	}

	return c.Status(201).JSON(inv)
}

// Get All Investments
func (h *InvestmentHandler) GetInvestments(c *fiber.Ctx) error {
	userID := c.Query("user_id") 
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "user_id required"})
	}

	var investments []models.Investment
	err := h.DB.Select(&investments, "SELECT * FROM investments WHERE user_id=$1 ORDER BY created_at DESC", userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch", "details": err.Error()})
	}

	return c.JSON(investments)
}

// Get Single Investment
func (h *InvestmentHandler) GetInvestment(c *fiber.Ctx) error {
	id := c.Params("id")
	var inv models.Investment
	err := h.DB.Get(&inv, "SELECT * FROM investments WHERE investment_id=$1", id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}
	return c.JSON(inv)
}

// Update Investment
func (h *InvestmentHandler) UpdateInvestment(c *fiber.Ctx) error {
	id := c.Params("id")
	var inv models.Investment
	if err := c.BodyParser(&inv); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	inv.InvestmentID = uuid.MustParse(id)

	_, err := h.DB.NamedExec(`
		UPDATE investments 
		SET type=:type, amount_invested=:amount_invested, current_value=:current_value, 
		    date_invested=:date_invested, reminder_date=:reminder_date, notes=:notes
		WHERE investment_id=:investment_id
	`, inv)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB update failed", "details": err.Error()})
	}

	return c.JSON(inv)
}

// Delete Investment
func (h *InvestmentHandler) DeleteInvestment(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := h.DB.Exec("DELETE FROM investments WHERE investment_id=$1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB delete failed"})
	}
	return c.SendStatus(204)
}
