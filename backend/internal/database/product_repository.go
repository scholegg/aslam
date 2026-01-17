package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aslam/backend/internal/models"
)

func (d *DB) CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error) {
	query := `
		INSERT INTO products (sku, name, volume, weight)
		VALUES ($1, $2, $3, $4)
		RETURNING sku, name, volume, weight, created_at, updated_at
	`

	product := &models.Product{}
	err := d.conn.QueryRowContext(ctx, query, req.SKU, req.Name, req.Volume, req.Weight).
		Scan(&product.SKU, &product.Name, &product.Volume, &product.Weight, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"products_pkey\"" {
			return nil, errors.New("product with this SKU already exists")
		}
		return nil, err
	}

	return product, nil
}

func (d *DB) GetProductBySKU(ctx context.Context, sku string) (*models.Product, error) {
	query := `
		SELECT sku, name, volume, weight, created_at, updated_at
		FROM products
		WHERE sku = $1
	`

	product := &models.Product{}
	err := d.conn.QueryRowContext(ctx, query, sku).
		Scan(&product.SKU, &product.Name, &product.Volume, &product.Weight, &product.CreatedAt, &product.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (d *DB) ListProducts(ctx context.Context) ([]models.Product, error) {
	query := `
		SELECT sku, name, volume, weight, created_at, updated_at
		FROM products
		ORDER BY name ASC
	`

	rows, err := d.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.SKU, &product.Name, &product.Volume, &product.Weight, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, rows.Err()
}

func (d *DB) UpdateProduct(ctx context.Context, sku string, req *models.UpdateProductRequest) (*models.Product, error) {
	query := `
		UPDATE products
		SET name = COALESCE(NULLIF($1, ''), name),
		    volume = CASE WHEN $2 > 0 THEN $2 ELSE volume END,
		    weight = CASE WHEN $3 > 0 THEN $3 ELSE weight END,
		    updated_at = CURRENT_TIMESTAMP
		WHERE sku = $4
		RETURNING sku, name, volume, weight, created_at, updated_at
	`

	product := &models.Product{}
	err := d.conn.QueryRowContext(ctx, query, req.Name, req.Volume, req.Weight, sku).
		Scan(&product.SKU, &product.Name, &product.Volume, &product.Weight, &product.CreatedAt, &product.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (d *DB) DeleteProduct(ctx context.Context, sku string) error {
	// Check if product is in use
	checkQuery := `SELECT COUNT(*) FROM shelf_items WHERE sku = $1`
	var count int
	err := d.conn.QueryRowContext(ctx, checkQuery, sku).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("cannot delete product that is in use on shelves")
	}

	query := `DELETE FROM products WHERE sku = $1`
	result, err := d.conn.ExecContext(ctx, query, sku)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product not found")
	}

	return nil
}
