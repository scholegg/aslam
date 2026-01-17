package models

import (
	"time"
)

type Shelf struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	RowIndex  int       `db:"row_index" json:"row_index"`
	ColIndex  int       `db:"col_index" json:"col_index"`
	MaxVolume float64   `db:"max_volume" json:"max_volume"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateShelfRequest struct {
	Name      string  `json:"name" binding:"required,min=3,max=100"`
	RowIndex  int     `json:"row_index" binding:"required,min=0"`
	ColIndex  int     `json:"col_index" binding:"required,min=0"`
	MaxVolume float64 `json:"max_volume" binding:"required,gt=0"`
}

type UpdateShelfRequest struct {
	Name      string  `json:"name" binding:"min=3,max=100"`
	MaxVolume float64 `json:"max_volume" binding:"gt=0"`
}

type ShelfResponse struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	RowIndex   int         `json:"row_index"`
	ColIndex   int         `json:"col_index"`
	MaxVolume  float64     `json:"max_volume"`
	UsedVolume float64     `json:"used_volume"`
	Items      []ShelfItem `json:"items"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type ShelfItem struct {
	ID          string    `json:"id"`
	ShelfID     string    `json:"shelf_id"`
	SKU         string    `json:"sku"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	Volume      float64   `json:"volume"`
	CreatedAt   time.Time `json:"created_at"`
}

type AddItemToShelfRequest struct {
	SKU      string `json:"sku" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
}

type RemoveItemRequest struct {
	ItemID string `json:"item_id" binding:"required"`
}
