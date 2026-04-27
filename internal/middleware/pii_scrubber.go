package middleware

import (
	"encoding/json"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

var (
	emailRegex = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	phoneRegex = regexp.MustCompile(`(\+62|08)[0-9]{9,12}`)
	nikRegex   = regexp.MustCompile(`\b\d{16}\b`)
)

func PIIScrubber(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPost {
		return c.Next()
	}

	body := c.Body()
	if len(body) == 0 {
		return c.Next()
	}

	// Simple string replacement on the raw body for speed
	// Better way: Parse JSON and scrub specific fields (like "content")
	bodyStr := string(body)
	
	scrubbed := emailRegex.ReplaceAllString(bodyStr, "[EMAIL_REDACTED]")
	scrubbed = phoneRegex.ReplaceAllString(scrubbed, "[PHONE_REDACTED]")
	scrubbed = nikRegex.ReplaceAllString(scrubbed, "[NIK_REDACTED]")

	if scrubbed != bodyStr {
		// Update body for the next handler
		c.Request().SetBody([]byte(scrubbed))
		
		// Also update Parsed body if it was already parsed
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(scrubbed), &parsed); err == nil {
			// This is a bit hacky in Fiber, usually better to re-parse in handler
			// or use a custom context key
		}
	}

	return c.Next()
}
