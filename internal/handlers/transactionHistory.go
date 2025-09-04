package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/PragaL15/Expense-Tracker/internal/database"

)

// TransactionWithCategory is the response structure
type TransactionWithCategory struct {
	TransactionID string    `json:"transaction_id"`
	UserID        string    `json:"user_id"`
	CategoryID    string    `json:"category_id"`
	CategoryName  string    `json:"category_name"`
	Amount        float64   `json:"amount"`
	TransactionType string  `json:"transaction_type"`
	Date          time.Time `json:"date"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
}

// ListTransactions returns all transactions with category name
func ListTransactions(c *fiber.Ctx) error {
	uid, ok := c.Locals("user_id").(string)
	if !ok || uid == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
	}

	var transactions []TransactionWithCategory

	q := database.DB.Table("transactions t").
		Select(`t.transaction_id,
		        t.user_id,
		        t.category_id,
		        COALESCE(c.name, '') AS category_name,
		        t.amount,
		        t.transaction_type,
		        t.date,
		        t.notes,
		        t.created_at`).
		Joins("LEFT JOIN categories c ON c.category_id = t.category_id").
		Where("t.user_id = ?", uid)


	if cid := c.Query("category_id"); cid != "" {
		q = q.Where("t.category_id = ?", cid)
	}
	if from := c.Query("from"); from != "" {
		if d, err := time.Parse("2006-01-02", from); err == nil {
			q = q.Where("t.date >= ?", d)
		}
	}
	if to := c.Query("to"); to != "" {
		if d, err := time.Parse("2006-01-02", to); err == nil {
			q = q.Where("t.date <= ?", d)
		}
	}

	sortOrder := c.Query("sort", "desc")
	orderStr := "t.date DESC, t.created_at DESC"
	if sortOrder == "asc" {
		orderStr = "t.date ASC, t.created_at ASC"
	}

	if err := q.Order(orderStr).Find(&transactions).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(transactions)
}
