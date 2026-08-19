package postgresql

import (
	"log"

	productrepo "primo/service/internal/product/repository/postgresql"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func strPtr(v string) *string    { return &v }
func f64Ptr(v float64) *float64 { return &v }

func RunSeed(db *gorm.DB) {
	products := []productrepo.ProductEntity{
		{
			ID:          uuid.MustParse("a1b2c3d4-e5f6-7890-abcd-ef1234567801"),
			Name:        "iPhone 16 Pro",
			Description: strPtr("Apple iPhone 16 Pro 256GB"),
			Price:       48900,
			SalePrice:   f64Ptr(45900),
		},
		{
			ID:          uuid.MustParse("a1b2c3d4-e5f6-7890-abcd-ef1234567802"),
			Name:        "MacBook Air M3",
			Description: strPtr("Apple MacBook Air M3 15-inch 256GB"),
			Price:       44900,
			SalePrice:   f64Ptr(42900),
		},
		{
			ID:          uuid.MustParse("a1b2c3d4-e5f6-7890-abcd-ef1234567803"),
			Name:        "iPad Pro M4",
			Description: strPtr("Apple iPad Pro M4 11-inch 256GB"),
			Price:       39900,
			SalePrice:   f64Ptr(37900),
		},
	}

	for _, p := range products {
		result := db.FirstOrCreate(&p, productrepo.ProductEntity{ID: p.ID})
		if result.Error != nil {
			log.Printf("Failed to seed product %s: %v", p.Name, result.Error)
		}
	}

	log.Println("Seed completed successfully")
}
