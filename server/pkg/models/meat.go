package models

import (
	"time"

	"github.com/uptrace/bun"
)

// Meat represents a piece of meat handled by the butcher.
type Meat struct {
	bun.BaseModel `bun:"table:meats,alias:meat"`

	Id         int64     `bun:"id,pk,autoincrement" json:"id"`
	Name       string    `bun:"name,notnull" json:"name"`
	Cut        string    `bun:"cut,notnull" json:"cut"`
	Species    string    `bun:"species,notnull" json:"species"`
	WeightG    float64   `bun:"weight_g,notnull" json:"weightG"`
	PricePerKg float64   `bun:"price_per_kg,notnull" json:"pricePerKg"`
	InStock    bool      `bun:"in_stock,notnull,default:true" json:"inStock"`
	ImagePath  string    `bun:"image_path" json:"imagePath"`
	CreatedAt  time.Time `bun:"created_at,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt  time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updatedAt"`
}
