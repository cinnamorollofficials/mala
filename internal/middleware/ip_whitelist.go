package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func IPWhitelist() fiber.Handler {
	return func(c *fiber.Ctx) error {
		allowedIPs := os.Getenv("ALLOWED_IPS")
		if allowedIPs == "" || allowedIPs == "*" {
			return c.Next()
		}

		clientIP := c.IP()
		allowedList := strings.Split(allowedIPs, ",")

		for _, ip := range allowedList {
			if strings.TrimSpace(ip) == clientIP {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "IP not allowed",
		})
	}
}
