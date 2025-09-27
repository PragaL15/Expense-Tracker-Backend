package routes

import (
	"github.com/PragaL15/Expense-Tracker/internal/handlers"
	"github.com/PragaL15/Expense-Tracker/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	api := app.Group("/api/v1")

	// --- Auth ---
	api.Post("/auth/register", handlers.Register)
	api.Post("/auth/login", handlers.Login)

	// --- Protected Routes ---
	p := api.Group("", middleware.JWTProtected())
	p.Get("/profile", handlers.Profile)

	// --- Categories ---
	p.Post("/categories", handlers.CreateCategory)
	p.Get("/categories", handlers.ListCategories)
	p.Put("/categories/:id", handlers.UpdateCategory)
	p.Delete("/categories/:id", handlers.DeleteCategory)

	// --- Transactions ---
	p.Post("/transactions/income", handlers.CreateIncome)
	p.Post("/transactions/expense", handlers.CreateExpense)
	p.Get("/transactions/income", handlers.ListIncome)
	p.Get("/transactions/expense", handlers.ListExpenses)
	p.Get("/transactionHistory", handlers.ListTransactions)
	p.Put("/transactions/:id", handlers.UpdateTransaction)
	p.Delete("/transactions/:id", handlers.DeleteTransaction)

	// --- Budgets ---
	p.Post("/budgets", handlers.UpsertBudget)
	p.Get("/budgets", handlers.GetBudget)
	p.Get("/budgets/summary", handlers.BudgetSummary)

	// --- Category Budgets ---
	p.Post("/category-budgets", handlers.UpsertCategoryBudget)
	p.Get("/category-budgets", handlers.GetCategoryBudget)

	// --- Investments ---
	p.Post("/investments", handlers.AddInvestment)
	p.Get("/investments", handlers.GetInvestments)
	p.Get("/investments/:id", handlers.GetInvestment)
	p.Put("/investments/:id", handlers.UpdateInvestment)
	p.Delete("/investments/:id", handlers.DeleteInvestment)
}
