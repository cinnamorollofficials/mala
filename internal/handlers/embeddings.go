package handlers

import "github.com/gofiber/fiber/v2"

func Embeddings(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error": "Embeddings endpoint not yet implemented",
	})
}
