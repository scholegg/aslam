package models

import (
	"time"
)

type Product struct {
	SKU       string    `db:"sku" json:"sku"`
	Name      string    `db:"name" json:"name"`
	Volume    float64   `db:"volume" json:"volume"`
	Weight    float64   `db:"weight" json:"weight"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateProductRequest struct {
	SKU    string  `json:"sku" binding:"required,min=3,max=50"`
	Name   string  `json:"name" binding:"required,min=3,max=255"`
	Volume float64 `json:"volume" binding:"required,gt=0"`
	Weight float64 `json:"weight" binding:"required,gt=0"`
}

type UpdateProductRequest struct {
	Name   string  `json:"name" binding:"min=3,max=255"`
	Volume float64 `json:"volume" binding:"gt=0"`
	Weight float64 `json:"weight" binding:"gt=0"`
}
