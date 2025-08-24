package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/PragaL15/Expense-Tracker/internal/database"
	"github.com/PragaL15/Expense-Tracker/internal/models"
)

// ListTransactions returns all transactions (income + expenses) for the user
// Optional query params:
//   ?category_id=<uuid>
//   ?from=YYYY-MM-DD
//   ?to=YYYY-MM-DD
//   ?sort=asc|desc  (default: desc)
func ListTransactions(c *fiber.Ctx) error {
	uid, ok := c.Locals("user_id").(string)
	if !ok || uid == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
	}

	var transactions []models.Transaction
	q := database.DB.Where("user_id = ?", uid)

	// Optional filters
	if cid := c.Query("category_id"); cid != "" {
		q = q.Where("category_id = ?", cid)
	}
	if from := c.Query("from"); from != "" {
		if d, err := time.Parse("2006-01-02", from); err == nil {
			q = q.Where("date >= ?", d)
		}
	}
	if to := c.Query("to"); to != "" {
		if d, err := time.Parse("2006-01-02", to); err == nil {
			q = q.Where("date <= ?", d)
		}
	}

	sortOrder := c.Query("sort", "desc") // default is descending
	orderStr := "date DESC, created_at DESC"
	if sortOrder == "asc" {
		orderStr = "date ASC, created_at ASC"
	}

	// Fetch transactions
	if err := q.Order(orderStr).Find(&transactions).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return as JSON
	return c.JSON(transactions)
}
