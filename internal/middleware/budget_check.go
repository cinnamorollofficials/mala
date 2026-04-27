package middleware

import "github.com/gofiber/fiber/v2"

func BudgetCheck(c *fiber.Ctx) error {
	budgetRemaining := c.Locals("budget_remaining").(float64)
	totalBudget := c.Locals("total_budget").(float64)

	// If totalBudget is 0, we assume unlimited for now or handled elsewhere
	if totalBudget > 0 && budgetRemaining <= 0 {
		return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
			"error": "Budget exceeded for this Virtual Key. Please top up.",
		})
	}

	return c.Next()
}
