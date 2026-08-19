package postgresql

import "github.com/google/uuid"

type ProductEntity struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description *string   `gorm:"type:text"`
	Price       float64   `gorm:"type:decimal(10,2);not null;index:idx_price,sort:asc"`
	SalePrice   *float64  `gorm:"type:decimal(10,2);index:idx_sale_price,sort:asc"`
}

func (ProductEntity) TableName() string {
	return "tb_products"
}

type UpdateProductEntity struct {
	Name        *string  `gorm:"column:name"`
	Description *string  `gorm:"column:description"`
	Price       *float64 `gorm:"column:price"`
	SalePrice   *float64 `gorm:"column:sale_price"`
}

func (UpdateProductEntity) TableName() string {
	return "tb_products"
}
