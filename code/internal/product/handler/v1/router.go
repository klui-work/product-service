package v1

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, h IProductHandler) {
	prodGroup := router.Group("/product")
	{
		prodGroup.Get("/health", h.HealthCheck)
		prodGroup.Get("", h.GetAllProducts)
		prodGroup.Post("", h.CreateProduct)
		prodGroup.Patch("/:id", h.PatchProduct)
	}
}
