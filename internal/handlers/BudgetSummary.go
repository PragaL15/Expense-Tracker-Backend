package handlers

import (
	"github.com/PragaL15/Expense-Tracker/internal/database"
	"github.com/gofiber/fiber/v2"
)

type CategoryBudgetRow struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	BudgetLimit  float64 `json:"budget_limit"`
	Spent        float64 `json:"spent"`
}

type BudgetSummaryResponse struct {
	UserID          string              `json:"user_id"`
	CategoryBudgets []CategoryBudgetRow `json:"category_budgets"`
	TotalBudget     float64             `json:"total_budget"`
	TotalSpent      float64             `json:"total_spent"`
	TotalIncome     float64             `json:"total_income"`
}

func BudgetSummary(c *fiber.Ctx) error {
	uid, ok := c.Locals("user_id").(string)
	if !ok || uid == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
	}

	var rows []CategoryBudgetRow
	raw := `
WITH budgets AS (
  SELECT category_id, COALESCE(SUM(budget_limit),0) AS budget_limit
  FROM category_budgets
  WHERE user_id = ?
  GROUP BY category_id
),
spent AS (
  SELECT category_id, COALESCE(SUM(amount),0) AS spent
  FROM transactions
  WHERE user_id = ?
    AND transaction_type IN ('Expense', 'Investment')   -- ✅ count investments like expenses
  GROUP BY category_id
)
SELECT c.category_id,
       c.name AS category_name,
       COALESCE(budgets.budget_limit, 0) AS budget_limit,
       COALESCE(spent.spent, 0)          AS spent
FROM categories c
LEFT JOIN budgets ON budgets.category_id = c.category_id
LEFT JOIN spent   ON spent.category_id   = c.category_id
WHERE budgets.category_id IS NOT NULL OR spent.category_id IS NOT NULL
ORDER BY c.name;
`
	if err := database.DB.Raw(raw, uid, uid).Scan(&rows).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var totalBudget float64
	if err := database.DB.Raw(
		`SELECT COALESCE(SUM(budget_limit),0) FROM category_budgets WHERE user_id = ?`,
		uid,
	).Scan(&totalBudget).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var totalSpent float64
	if err := database.DB.Raw(
		`SELECT COALESCE(SUM(amount),0) 
         FROM transactions 
         WHERE user_id = ? 
           AND transaction_type IN ('Expense', 'Investment')`,  // ✅ include investments
		uid,
	).Scan(&totalSpent).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Total income (still only Income)
	var totalIncome float64
	if err := database.DB.Raw(
		`SELECT COALESCE(SUM(amount),0) 
         FROM transactions 
         WHERE user_id = ? 
           AND transaction_type = 'Income'`,
		uid,
	).Scan(&totalIncome).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	resp := BudgetSummaryResponse{
		UserID:          uid,
		CategoryBudgets: rows,
		TotalBudget:     totalBudget,
		TotalSpent:      totalSpent,
		TotalIncome:     totalIncome,
	}

	return c.JSON(resp)
}
