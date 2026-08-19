package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	productv1 "primo/service/internal/product/handler/v1"
	productrepo "primo/service/internal/product/repository/postgresql"
	productusecase "primo/service/internal/product/usecase"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func strPtr(v string) *string    { return &v }
func f64Ptr(v float64) *float64 { return &v }

type responseBody struct {
	Successful bool   `json:"successful"`
	ErrorCode  string `json:"error_code"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

type E2ETestSuite struct {
	suite.Suite
	app       *fiber.App
	db        *gorm.DB
	container testcontainers.Container
	ctx       context.Context
}

func (s *E2ETestSuite) SetupSuite() {
	s.ctx = context.Background()

	pgContainer, err := tcpostgres.Run(s.ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("testuser"),
		tcpostgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	s.Require().NoError(err)
	s.container = pgContainer

	host, err := pgContainer.Host(s.ctx)
	s.Require().NoError(err)

	port, err := pgContainer.MappedPort(s.ctx, "5432")
	s.Require().NoError(err)

	dsn := fmt.Sprintf("host=%s user=testuser password=testpass dbname=testdb port=%s sslmode=disable", host, port.Port())

	s.db, err = gorm.Open(pg.Open(dsn), &gorm.Config{})
	s.Require().NoError(err)

	err = s.db.AutoMigrate(&productrepo.ProductEntity{})
	s.Require().NoError(err)

	repo := productrepo.NewProductRepository(s.db)
	svc := productusecase.NewService(repo)
	handler := productv1.NewProductHandler(svc)

	s.app = fiber.New()
	productv1.RegisterRoutes(s.app, handler)
}

func (s *E2ETestSuite) TearDownSuite() {
	if s.container != nil {
		_ = s.container.Terminate(s.ctx)
	}
}

func (s *E2ETestSuite) SetupTest() {
	s.db.Exec("DELETE FROM tb_products")
}

func (s *E2ETestSuite) doRequest(method, path string, body any) (*http.Response, responseBody) {
	var reqBody io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, "http://localhost"+path, reqBody)
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req, -1)
	s.Require().NoError(err)

	var result responseBody
	json.NewDecoder(resp.Body).Decode(&result)
	return resp, result
}

func (s *E2ETestSuite) doRequestRaw(method, path string, body any) (*http.Response, map[string]any) {
	var reqBody io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, "http://localhost"+path, reqBody)
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req, -1)
	s.Require().NoError(err)

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	return resp, result
}

func (s *E2ETestSuite) TestCreateProduct_OmitNullableFields() {
	body := map[string]any{
		"name":  "Minimal Product",
		"price": 100,
	}

	resp, result := s.doRequest("POST", "/product", body)

	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
	assert.True(s.T(), result.Successful)

	var found productrepo.ProductEntity
	s.db.Where("name = ?", "Minimal Product").First(&found)
	assert.Nil(s.T(), found.Description)
	assert.Nil(s.T(), found.SalePrice)
}

func (s *E2ETestSuite) TestCreateProduct_NullDescription() {
	body := map[string]any{
		"name":        "Null Desc Product",
		"description": nil,
		"price":       200,
	}

	resp, result := s.doRequest("POST", "/product", body)

	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
	assert.True(s.T(), result.Successful)

	var found productrepo.ProductEntity
	s.db.Where("name = ?", "Null Desc Product").First(&found)
	assert.Nil(s.T(), found.Description)
}

func (s *E2ETestSuite) TestPatchProduct_ResponseFormat() {
	id := uuid.New()
	s.db.Create(&productrepo.ProductEntity{
		ID:    id,
		Name:  "Product",
		Price: 100,
	})

	resp, result := s.doRequestRaw("PATCH", "/product/"+id.String(), map[string]any{
		"name": "Updated",
	})

	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
	assert.Equal(s.T(), 2, len(result))
	assert.Equal(s.T(), true, result["successful"])
	assert.Equal(s.T(), "", result["error_code"])
	_, hasMessage := result["message"]
	_, hasData := result["data"]
	assert.False(s.T(), hasMessage)
	assert.False(s.T(), hasData)
}

func (s *E2ETestSuite) TestHealthCheck() {
	resp, result := s.doRequest("GET", "/product/health", nil)

	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
	assert.True(s.T(), result.Successful)
}

func (s *E2ETestSuite) TestCreateProduct_Success() {
	body := map[string]any{
		"name":        "iPhone 16",
		"description": "Apple iPhone 16",
		"price":       48900,
		"sale_price":  45900,
	}

	resp, result := s.doRequest("POST", "/product", body)

	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
	assert.True(s.T(), result.Successful)
	assert.NotNil(s.T(), result.Data)
	_, isObject := result.Data.(map[string]any)
	assert.True(s.T(), isObject, "create success data should be {} not null")

	var count int64
	s.db.Model(&productrepo.ProductEntity{}).Where("name = ?", "iPhone 16").Count(&count)
	assert.Equal(s.T(), int64(1), count)
}

func (s *E2ETestSuite) TestCreateProduct_MissingName() {
	body := map[string]any{
		"price": 48900,
	}

	resp, result := s.doRequest("POST", "/product", body)

	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	assert.False(s.T(), result.Successful)
	assert.Equal(s.T(), "E_VALIDATION_REQUIRED", result.ErrorCode)
	assert.NotNil(s.T(), result.Data)
	_, isObject := result.Data.(map[string]any)
	assert.True(s.T(), isObject, "error data should be {} not null")
}

func (s *E2ETestSuite) TestCreateProduct_ZeroPrice() {
	body := map[string]any{
		"name":  "iPhone",
		"price": 0,
	}

	resp, result := s.doRequest("POST", "/product", body)

	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	assert.False(s.T(), result.Successful)
}

func (s *E2ETestSuite) TestCreateProduct_InvalidBody() {
	req, _ := http.NewRequest("POST", "http://localhost/product", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req, -1)
	s.Require().NoError(err)

	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
}

func (s *E2ETestSuite) TestPatchProduct_Success() {
	id := uuid.New()
	s.db.Create(&productrepo.ProductEntity{
		ID:          id,
		Name:        "Original",
		Description: strPtr("Original Desc"),
		Price:       100,
		SalePrice:   f64Ptr(80),
	})

	body := map[string]any{
		"name": "Updated Name",
	}

	resp, result := s.doRequest("PATCH", "/product/"+id.String(), body)

	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
	assert.True(s.T(), result.Successful)

	var found productrepo.ProductEntity
	s.db.First(&found, "id = ?", id)
	assert.Equal(s.T(), "Updated Name", found.Name)
	assert.Equal(s.T(), strPtr("Original Desc"), found.Description)
	assert.Equal(s.T(), 100.0, found.Price)
}

func (s *E2ETestSuite) TestPatchProduct_PartialUpdate() {
	id := uuid.New()
	s.db.Create(&productrepo.ProductEntity{
		ID:        id,
		Name:      "Product",
		Price:     100,
		SalePrice: f64Ptr(80),
	})

	body := map[string]any{
		"price":      150,
		"sale_price": 120,
	}

	resp, result := s.doRequest("PATCH", "/product/"+id.String(), body)

	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
	assert.True(s.T(), result.Successful)

	var found productrepo.ProductEntity
	s.db.First(&found, "id = ?", id)
	assert.Equal(s.T(), "Product", found.Name)
	assert.Equal(s.T(), 150.0, found.Price)
	assert.Equal(s.T(), f64Ptr(120.0), found.SalePrice)
}

func (s *E2ETestSuite) TestPatchProduct_NotFound() {
	body := map[string]any{
		"name": "Updated",
	}

	resp, result := s.doRequest("PATCH", "/product/"+uuid.New().String(), body)

	assert.Equal(s.T(), http.StatusNotFound, resp.StatusCode)
	assert.False(s.T(), result.Successful)
	assert.Equal(s.T(), "E_NOT_FOUND", result.ErrorCode)
}

func (s *E2ETestSuite) TestPatchProduct_InvalidUUID() {
	body := map[string]any{
		"name": "Updated",
	}

	resp, result := s.doRequest("PATCH", "/product/invalid-uuid", body)

	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	assert.False(s.T(), result.Successful)
}

func (s *E2ETestSuite) TestPatchProduct_EmptyName() {
	id := uuid.New()
	s.db.Create(&productrepo.ProductEntity{
		ID:    id,
		Name:  "Product",
		Price: 100,
	})

	body := map[string]any{
		"name": "",
	}

	resp, result := s.doRequest("PATCH", "/product/"+id.String(), body)

	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	assert.False(s.T(), result.Successful)
}

func TestE2ESuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
