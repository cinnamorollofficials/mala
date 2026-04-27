package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hadigunawan/mala/internal/handlers"
	"github.com/hadigunawan/mala/internal/middleware"
)

func SetupRoutes(app *fiber.App) {
	// Health Check
	app.Get("/api/health", handlers.HealthCheck)

	// Data Plane (OpenAI-Compatible)
	v1 := app.Group("/v1", middleware.VirtualKeyMiddleware)
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

	// Admin Usage & Analytics
	admin.Get("/usage/summary", handlers.GetUsageSummary)
	admin.Get("/usage/:key_id", handlers.GetKeyUsageHistory)
}
