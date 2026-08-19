package repo

import (
	"context"
	"fmt"
	"primo/service/internal/product/domain"
	productrepo "primo/service/internal/product/repository/postgresql"
	portout "primo/service/internal/product/port/out"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func strPtr(v string) *string    { return &v }
func f64Ptr(v float64) *float64 { return &v }

type ProductRepoTestSuite struct {
	suite.Suite
	db        *gorm.DB
	repo      portout.ProductRepository
	container testcontainers.Container
	ctx       context.Context
}

func (s *ProductRepoTestSuite) SetupSuite() {
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

	s.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	s.Require().NoError(err)

	err = s.db.AutoMigrate(&productrepo.ProductEntity{})
	s.Require().NoError(err)

	s.repo = productrepo.NewProductRepository(s.db)
}

func (s *ProductRepoTestSuite) TearDownSuite() {
	if s.container != nil {
		_ = s.container.Terminate(s.ctx)
	}
}

func (s *ProductRepoTestSuite) SetupTest() {
	s.db.Exec("DELETE FROM tb_products")
}

func (s *ProductRepoTestSuite) TestCreateProduct_Success() {
	product := domain.Product{
		ID:          uuid.New(),
		Name:        "Test Product",
		Description: strPtr("Test Description"),
		Price:       100.00,
		SalePrice:   f64Ptr(80.00),
	}

	err := s.repo.CreateProduct(s.ctx, product)
	assert.NoError(s.T(), err)

	var found productrepo.ProductEntity
	s.db.First(&found, "id = ?", product.ID)
	assert.Equal(s.T(), product.Name, found.Name)
	assert.Equal(s.T(), product.Price, found.Price)
	assert.Equal(s.T(), product.Description, found.Description)
}

func (s *ProductRepoTestSuite) TestCreateProduct_NullableFields() {
	product := domain.Product{
		ID:    uuid.New(),
		Name:  "Minimal Product",
		Price: 100,
	}

	err := s.repo.CreateProduct(s.ctx, product)
	assert.NoError(s.T(), err)

	var found productrepo.ProductEntity
	s.db.First(&found, "id = ?", product.ID)
	assert.Nil(s.T(), found.Description)
	assert.Nil(s.T(), found.SalePrice)
}

func (s *ProductRepoTestSuite) TestCreateProduct_DuplicateID() {
	id := uuid.New()
	product := domain.Product{
		ID:    id,
		Name:  "Product 1",
		Price: 100,
	}

	err := s.repo.CreateProduct(s.ctx, product)
	assert.NoError(s.T(), err)

	product2 := domain.Product{
		ID:    id,
		Name:  "Product 2",
		Price: 200,
	}

	err = s.repo.CreateProduct(s.ctx, product2)
	assert.Error(s.T(), err)
}

func (s *ProductRepoTestSuite) TestGetAllProducts_Success() {
	id1 := uuid.New()
	id2 := uuid.New()
	s.db.Create(&productrepo.ProductEntity{ID: id1, Name: "Product 1", Price: 100})
	s.db.Create(&productrepo.ProductEntity{ID: id2, Name: "Product 2", Price: 200})

	products, err := s.repo.GetAllProducts(s.ctx)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), products, 2)
}

func (s *ProductRepoTestSuite) TestPatchProduct_Success() {
	id := uuid.New()
	entity := productrepo.ProductEntity{
		ID:          id,
		Name:        "Original",
		Description: strPtr("Original Desc"),
		Price:       100.00,
		SalePrice:   f64Ptr(80.00),
	}
	s.db.Create(&entity)

	newName := "Updated Name"
	update := domain.UpdateProduct{
		Name: &newName,
	}

	err := s.repo.PatchProduct(s.ctx, id, update)
	assert.NoError(s.T(), err)

	var found productrepo.ProductEntity
	s.db.First(&found, "id = ?", id)
	assert.Equal(s.T(), "Updated Name", found.Name)
	assert.Equal(s.T(), strPtr("Original Desc"), found.Description)
	assert.Equal(s.T(), 100.00, found.Price)
}

func (s *ProductRepoTestSuite) TestPatchProduct_NotFound() {
	newName := "Updated"
	update := domain.UpdateProduct{
		Name: &newName,
	}

	err := s.repo.PatchProduct(s.ctx, uuid.New(), update)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "product not found")
}

func (s *ProductRepoTestSuite) TestPatchProduct_PartialUpdate() {
	id := uuid.New()
	entity := productrepo.ProductEntity{
		ID:          id,
		Name:        "Original",
		Description: strPtr("Original Desc"),
		Price:       100.00,
		SalePrice:   f64Ptr(80.00),
	}
	s.db.Create(&entity)

	newPrice := 150.00
	newSalePrice := 120.00
	update := domain.UpdateProduct{
		Price:     &newPrice,
		SalePrice: &newSalePrice,
	}

	err := s.repo.PatchProduct(s.ctx, id, update)
	assert.NoError(s.T(), err)

	var found productrepo.ProductEntity
	s.db.First(&found, "id = ?", id)
	assert.Equal(s.T(), "Original", found.Name)
	assert.Equal(s.T(), 150.00, found.Price)
	assert.Equal(s.T(), f64Ptr(120.00), found.SalePrice)
}

func TestProductRepoSuite(t *testing.T) {
	suite.Run(t, new(ProductRepoTestSuite))
}
