package postgresql

import (
	"log"

	"primo/service/internal/product/repository/postgresql"

	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm")

	err := db.AutoMigrate(
		&postgresql.ProductEntity{},
	)
	if err != nil {
		log.Fatalf("Failed to run migration: %v", err)
	}

	db.Exec("CREATE INDEX IF NOT EXISTS idx_pg_gin_name ON tb_products USING gin (name gin_trgm_ops)")

	log.Println("Migration completed successfully")
}
