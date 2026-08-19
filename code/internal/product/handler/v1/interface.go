package v1

import "github.com/gofiber/fiber/v2"

type IProductHandler interface {
	HealthCheck(c *fiber.Ctx) error
	GetAllProducts(c *fiber.Ctx) error
	CreateProduct(c *fiber.Ctx) error
	PatchProduct(c *fiber.Ctx) error
}
