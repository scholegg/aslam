package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func New(dsn string) (*DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return &DB{conn: db}, nil
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func (d *DB) GetConnection() *sql.DB {
	return d.conn
}

func (d *DB) RunMigrations() error {
	migrations := []string{
		createUsersTable,
		createProductsTable,
		createShelfsTable,
		createShelfItemsTable,
	}

	for _, migration := range migrations {
		if _, err := d.conn.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	log.Println("✓ Database migrations completed successfully")
	return nil
}

const (
	createUsersTable = `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'viewer',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`

	createProductsTable = `
		CREATE TABLE IF NOT EXISTS products (
			sku VARCHAR(50) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			volume DECIMAL(10, 2) NOT NULL,
			weight DECIMAL(10, 2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
	`

	createShelfsTable = `
		CREATE TABLE IF NOT EXISTS shelfs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(100) NOT NULL,
			row_index INTEGER NOT NULL,
			col_index INTEGER NOT NULL,
			max_volume DECIMAL(10, 2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_shelfs_position ON shelfs(row_index, col_index);
	`

	createShelfItemsTable = `
		CREATE TABLE IF NOT EXISTS shelf_items (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			shelf_id UUID NOT NULL REFERENCES shelfs(id) ON DELETE CASCADE,
			sku VARCHAR(50) NOT NULL REFERENCES products(sku) ON DELETE RESTRICT,
			quantity INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT shelf_item_quantity_positive CHECK (quantity > 0)
		);

		CREATE INDEX IF NOT EXISTS idx_shelf_items_shelf_id ON shelf_items(shelf_id);
		CREATE INDEX IF NOT EXISTS idx_shelf_items_sku ON shelf_items(sku);
	`
)
