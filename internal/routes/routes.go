package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/hadigunawan/mala/internal/handlers"
	"github.com/hadigunawan/mala/internal/middleware"
)

func SetupRoutes(app *fiber.App) {
	// Health Check
	app.Get("/api/health", handlers.HealthCheck)

	// Data Plane (OpenAI-Compatible)
	v1 := app.Group("/v1")
	
	// Apply security chain in order
	v1.Use(middleware.IPWhitelist())
	v1.Use(middleware.VirtualKeyAuth)
	v1.Use(middleware.BudgetCheck)
	v1.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Locals("virtual_key_id").(string)
		},
	}))
	v1.Use(middleware.PIIScrubber)
	v1.Use(middleware.ProviderSetup)

	v1.Post("/chat/completions", handlers.ChatCompletions)
	v1.Post("/embeddings", handlers.Embeddings)
	v1.Get("/models", handlers.ListModels)

	// Control Plane (Admin)
	admin := app.Group("/admin")
	
	// Admin Keys
	admin.Post("/keys", handlers.CreateVirtualKey)
	admin.Get("/keys", handlers.ListVirtualKeys)
	admin.Patch("/keys/:id/topup", handlers.TopupVirtualKey)
	admin.Delete("/keys/:id", handlers.DeleteVirtualKey)

	// Admin Providers
	admin.Post("/providers", handlers.CreateProvider)
	admin.Get("/providers/health", handlers.ProvidersHealth)

	// Admin Models
	admin.Post("/models", handlers.CreateModel)
	admin.Get("/models", handlers.AdminListModels)
	admin.Put("/models/:id", handlers.UpdateModel)
	admin.Delete("/models/:id", handlers.DeleteModel)

	// Admin Usage & Analytics
	admin.Get("/usage/summary", handlers.GetUsageSummary)
	admin.Get("/usage/:key_id", handlers.GetKeyUsageHistory)
}
