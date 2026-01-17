package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aslam/backend/internal/models"
	"github.com/google/uuid"
)

func (d *DB) CreateShelf(ctx context.Context, req *models.CreateShelfRequest) (*models.Shelf, error) {
	id := uuid.New().String()

	query := `
		INSERT INTO shelfs (id, name, row_index, col_index, max_volume)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, row_index, col_index, max_volume, created_at, updated_at
	`

	shelf := &models.Shelf{}
	err := d.conn.QueryRowContext(ctx, query, id, req.Name, req.RowIndex, req.ColIndex, req.MaxVolume).
		Scan(&shelf.ID, &shelf.Name, &shelf.RowIndex, &shelf.ColIndex, &shelf.MaxVolume, &shelf.CreatedAt, &shelf.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return shelf, nil
}

func (d *DB) GetShelfByID(ctx context.Context, id string) (*models.ShelfResponse, error) {
	shelfQuery := `
		SELECT id, name, row_index, col_index, max_volume, created_at, updated_at
		FROM shelfs
		WHERE id = $1
	`

	shelf := &models.Shelf{}
	err := d.conn.QueryRowContext(ctx, shelfQuery, id).
		Scan(&shelf.ID, &shelf.Name, &shelf.RowIndex, &shelf.ColIndex, &shelf.MaxVolume, &shelf.CreatedAt, &shelf.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("shelf not found")
	}
	if err != nil {
		return nil, err
	}

	// Get items
	items, err := d.getShelfItems(ctx, id)
	if err != nil {
		return nil, err
	}

	// Calculate used volume
	usedVolume := 0.0
	for _, item := range items {
		usedVolume += item.Volume
	}

	response := &models.ShelfResponse{
		ID:         shelf.ID,
		Name:       shelf.Name,
		RowIndex:   shelf.RowIndex,
		ColIndex:   shelf.ColIndex,
		MaxVolume:  shelf.MaxVolume,
		UsedVolume: usedVolume,
		Items:      items,
		CreatedAt:  shelf.CreatedAt,
		UpdatedAt:  shelf.UpdatedAt,
	}

	return response, nil
}

func (d *DB) getShelfItems(ctx context.Context, shelfID string) ([]models.ShelfItem, error) {
	query := `
		SELECT si.id, si.shelf_id, si.sku, p.name, si.quantity, (p.volume * si.quantity) as volume, si.created_at
		FROM shelf_items si
		JOIN products p ON si.sku = p.sku
		WHERE si.shelf_id = $1
		ORDER BY si.created_at ASC
	`

	rows, err := d.conn.QueryContext(ctx, query, shelfID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.ShelfItem
	for rows.Next() {
		var item models.ShelfItem
		err := rows.Scan(&item.ID, &item.ShelfID, &item.SKU, &item.ProductName, &item.Quantity, &item.Volume, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (d *DB) ListShelfs(ctx context.Context) ([]models.ShelfResponse, error) {
	query := `
		SELECT id, name, row_index, col_index, max_volume, created_at, updated_at
		FROM shelfs
		ORDER BY row_index ASC, col_index ASC
	`

	rows, err := d.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shelfs []models.ShelfResponse
	for rows.Next() {
		var shelf models.Shelf
		err := rows.Scan(&shelf.ID, &shelf.Name, &shelf.RowIndex, &shelf.ColIndex, &shelf.MaxVolume, &shelf.CreatedAt, &shelf.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Get items and used volume
		items, err := d.getShelfItems(ctx, shelf.ID)
		if err != nil {
			return nil, err
		}

		usedVolume := 0.0
		for _, item := range items {
			usedVolume += item.Volume
		}

		shelfs = append(shelfs, models.ShelfResponse{
			ID:         shelf.ID,
			Name:       shelf.Name,
			RowIndex:   shelf.RowIndex,
			ColIndex:   shelf.ColIndex,
			MaxVolume:  shelf.MaxVolume,
			UsedVolume: usedVolume,
			Items:      items,
			CreatedAt:  shelf.CreatedAt,
			UpdatedAt:  shelf.UpdatedAt,
		})
	}

	return shelfs, rows.Err()
}

func (d *DB) UpdateShelf(ctx context.Context, id string, req *models.UpdateShelfRequest) (*models.Shelf, error) {
	query := `
		UPDATE shelfs
		SET name = COALESCE(NULLIF($1, ''), name),
		    max_volume = CASE WHEN $2 > 0 THEN $2 ELSE max_volume END,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, name, row_index, col_index, max_volume, created_at, updated_at
	`

	shelf := &models.Shelf{}
	err := d.conn.QueryRowContext(ctx, query, req.Name, req.MaxVolume, id).
		Scan(&shelf.ID, &shelf.Name, &shelf.RowIndex, &shelf.ColIndex, &shelf.MaxVolume, &shelf.CreatedAt, &shelf.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("shelf not found")
	}
	if err != nil {
		return nil, err
	}

	return shelf, nil
}

func (d *DB) DeleteShelf(ctx context.Context, id string) error {
	query := `DELETE FROM shelfs WHERE id = $1`
	result, err := d.conn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("shelf not found")
	}

	return nil
}

func (d *DB) AddItemToShelf(ctx context.Context, shelfID, sku string, quantity int) (*models.ShelfItem, error) {
	// Get product to verify existence and get volume
	product, err := d.GetProductBySKU(ctx, sku)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Get shelf to check volume
	shelf, err := d.GetShelfByID(ctx, shelfID)
	if err != nil {
		return nil, err
	}

	// Check if item already exists
	checkQuery := `SELECT id, quantity FROM shelf_items WHERE shelf_id = $1 AND sku = $2`
	var existingID string
	var existingQuantity int
	err = d.conn.QueryRowContext(ctx, checkQuery, shelfID, sku).Scan(&existingID, &existingQuantity)

	if err == nil {
		// Item exists, update quantity
		newQuantity := existingQuantity + quantity
		volumeAdded := product.Volume * float64(quantity)

		if shelf.UsedVolume+volumeAdded > shelf.MaxVolume {
			return nil, errors.New("insufficient shelf volume")
		}

		updateQuery := `
			UPDATE shelf_items
			SET quantity = $1
			WHERE id = $2
			RETURNING id, shelf_id, sku, quantity, created_at
		`

		item := &models.ShelfItem{}
		err := d.conn.QueryRowContext(ctx, updateQuery, newQuantity, existingID).
			Scan(&item.ID, &item.ShelfID, &item.SKU, &item.Quantity, &item.CreatedAt)

		if err != nil {
			return nil, err
		}

		item.ProductName = product.Name
		item.Volume = product.Volume * float64(item.Quantity)
		return item, nil
	}

	// New item
	volumeNeeded := product.Volume * float64(quantity)
	if shelf.UsedVolume+volumeNeeded > shelf.MaxVolume {
		return nil, errors.New("insufficient shelf volume")
	}

	id := uuid.New().String()
	insertQuery := `
		INSERT INTO shelf_items (id, shelf_id, sku, quantity)
		VALUES ($1, $2, $3, $4)
		RETURNING id, shelf_id, sku, quantity, created_at
	`

	item := &models.ShelfItem{}
	err = d.conn.QueryRowContext(ctx, insertQuery, id, shelfID, sku, quantity).
		Scan(&item.ID, &item.ShelfID, &item.SKU, &item.Quantity, &item.CreatedAt)

	if err != nil {
		return nil, err
	}

	item.ProductName = product.Name
	item.Volume = product.Volume * float64(item.Quantity)
	return item, nil
}

func (d *DB) RemoveItemFromShelf(ctx context.Context, itemID string) error {
	query := `DELETE FROM shelf_items WHERE id = $1`
	result, err := d.conn.ExecContext(ctx, query, itemID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("item not found")
	}

	return nil
}

func (d *DB) UpdateItemQuantity(ctx context.Context, itemID string, quantity int) error {
	if quantity <= 0 {
		return d.RemoveItemFromShelf(ctx, itemID)
	}

	query := `UPDATE shelf_items SET quantity = $1 WHERE id = $2`
	result, err := d.conn.ExecContext(ctx, query, quantity, itemID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("item not found")
	}

	return nil
}
