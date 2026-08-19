package v1

import (
	"primo/service/internal/product/application"
	"primo/service/internal/product/handler/mapper"
	"primo/service/internal/product/handler/response"
	"primo/service/internal/product/handler/validation"
	portin "primo/service/internal/product/port/in"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productService portin.ProductService
}

func NewProductHandler(productService portin.ProductService) IProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// HealthCheck godoc
// @Summary      Health check
// @Description  Check if the Product API is healthy
// @Tags         product
// @Produce      json
// @Success      200  {object}  response.HealthCheckResponse
// @Router       /product/health [get]
func (h *ProductHandler) HealthCheck(c *fiber.Ctx) error {
	return response.Success(c, fiber.Map{
		"msg": "Product API is healthy",
	})
}

// GetAllProducts godoc
// @Summary      Get all products
// @Description  Retrieve a list of all products
// @Tags         product
// @Produce      json
// @Success      200  {object}  response.ProductListResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /product [get]
func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	ctx := c.UserContext()

	products, err := h.productService.GetAllProducts(ctx)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, mapper.MapProductsToDtos(products))
}

// CreateProduct godoc
// @Summary      Create a product
// @Description  Create a new product
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        body  body      application.CreateProductDto  true  "Create product request"
// @Success      200   {object}  response.CreateProductResponse
// @Failure      400   {object}  response.ErrorResponse
// @Failure      500   {object}  response.ErrorResponse
// @Router       /product [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var dto application.CreateProductDto
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, response.NewValidationError(response.ErrBadRequest, err.Error()))
	}

	if err := validation.ValidateCreateProductDto(dto); err != nil {
		return response.Error(c, err)
	}

	if err := h.productService.CreateProduct(ctx, dto); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c)
}

// PatchProduct godoc
// @Summary      Patch a product
// @Description  Partially update a product by ID. Uses seed data ID as example.
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id    path      string                        true  "Product ID (UUID)"  example(a1b2c3d4-e5f6-7890-abcd-ef1234567801)
// @Param        body  body      application.PatchProductDto   true  "Patch product request"
// @Success      200   {object}  response.PatchResponse
// @Failure      400   {object}  response.ErrorResponse
// @Failure      404   {object}  response.ErrorResponse
// @Failure      500   {object}  response.ErrorResponse
// @Router       /product/{id} [patch]
func (h *ProductHandler) PatchProduct(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id := c.Params("id")

	if err := validation.ValidateUUID(id); err != nil {
		return response.Error(c, err)
	}

	productID, err := uuid.Parse(id)
	if err != nil {
		return response.Error(c, response.NewValidationError(response.ErrBadRequest, "invalid product id"))
	}

	var dto application.PatchProductDto
	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, response.NewValidationError(response.ErrBadRequest, err.Error()))
	}

	if err := validation.ValidatePatchProductDto(dto); err != nil {
		return response.Error(c, err)
	}

	if err := h.productService.PatchProduct(ctx, productID, dto); err != nil {
		return response.Error(c, err)
	}

	return response.SuccessPatch(c)
}
