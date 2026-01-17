package tests

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aslam/backend/internal/database"
	"github.com/aslam/backend/internal/models"
)

func TestUserOperations(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Test create user
	user, err := db.CreateUser(ctx, "test@example.com", "password123", models.RoleEditor)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", user.Email)
	}

	// Test get user
	retrieved, err := db.GetUserByEmail(ctx, "test@example.com")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if retrieved.ID != user.ID {
		t.Errorf("User ID mismatch")
	}

	// Test validate password
	if !db.ValidatePassword(user.Password, "password123") {
		t.Error("Password validation failed")
	}

	// Test delete user
	err = db.DeleteUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}
}

func TestProductOperations(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Test create product
	req := &models.CreateProductRequest{
		SKU:    "SKU001",
		Name:   "Test Product",
		Volume: 10.5,
		Weight: 2.5,
	}

	product, err := db.CreateProduct(ctx, req)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	if product.SKU != "SKU001" {
		t.Errorf("Expected SKU SKU001, got %s", product.SKU)
	}

	// Test get product
	retrieved, err := db.GetProductBySKU(ctx, "SKU001")
	if err != nil {
		t.Fatalf("Failed to get product: %v", err)
	}

	if retrieved.Volume != 10.5 {
		t.Errorf("Expected volume 10.5, got %f", retrieved.Volume)
	}

	// Test list products
	products, err := db.ListProducts(ctx)
	if err != nil {
		t.Fatalf("Failed to list products: %v", err)
	}

	if len(products) == 0 {
		t.Error("Expected at least one product")
	}

	// Test delete product
	err = db.DeleteProduct(ctx, "SKU001")
	if err != nil {
		t.Fatalf("Failed to delete product: %v", err)
	}
}

func TestShelfOperations(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create a product first
	productReq := &models.CreateProductRequest{
		SKU:    "SKU002",
		Name:   "Shelf Product",
		Volume: 5.0,
		Weight: 1.0,
	}
	_, err := db.CreateProduct(ctx, productReq)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Test create shelf
	shelfReq := &models.CreateShelfRequest{
		Name:      "Shelf A",
		RowIndex:  0,
		ColIndex:  0,
		MaxVolume: 100.0,
	}

	shelf, err := db.CreateShelf(ctx, shelfReq)
	if err != nil {
		t.Fatalf("Failed to create shelf: %v", err)
	}

	// Test add item to shelf
	item, err := db.AddItemToShelf(ctx, shelf.ID, "SKU002", 5)
	if err != nil {
		t.Fatalf("Failed to add item to shelf: %v", err)
	}

	if item.SKU != "SKU002" {
		t.Errorf("Expected SKU SKU002, got %s", item.SKU)
	}

	if item.Quantity != 5 {
		t.Errorf("Expected quantity 5, got %d", item.Quantity)
	}

	// Test get shelf
	retrieved, err := db.GetShelfByID(ctx, shelf.ID)
	if err != nil {
		t.Fatalf("Failed to get shelf: %v", err)
	}

	if retrieved.UsedVolume != 25.0 {
		t.Errorf("Expected used volume 25.0, got %f", retrieved.UsedVolume)
	}

	if len(retrieved.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(retrieved.Items))
	}

	// Test volume validation
	_, err = db.AddItemToShelf(ctx, shelf.ID, "SKU002", 50)
	if err == nil {
		t.Error("Expected volume exceeded error")
	}

	// Test remove item
	err = db.RemoveItemFromShelf(ctx, item.ID)
	if err != nil {
		t.Fatalf("Failed to remove item: %v", err)
	}

	// Test delete shelf
	err = db.DeleteShelf(ctx, shelf.ID)
	if err != nil {
		t.Fatalf("Failed to delete shelf: %v", err)
	}
}

func setupTestDB(t *testing.T) *database.DB {
	dbHost := os.Getenv("TEST_DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("TEST_DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := os.Getenv("TEST_DB_USER")
	if dbUser == "" {
		dbUser = "aslam_user"
	}

	dbPassword := os.Getenv("TEST_DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "aslam_password"
	}

	dbName := os.Getenv("TEST_DB_NAME")
	if dbName == "" {
		dbName = "aslam_test"
	}

	// AGORA sim, fmt.Sprintf faz sentido
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	db, err := database.New(dsn)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.RunMigrations(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}
