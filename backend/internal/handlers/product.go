package handlers

import (
	"net/http"

	"github.com/aslam/backend/internal/database"
	"github.com/aslam/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func CreateProduct(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only editor and admin can create products
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		if userRole != "editor" && userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		var req models.CreateProductRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product, err := db.CreateProduct(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, product)
	}
}

func GetProduct(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sku := c.Param("sku")
		product, err := db.GetProductBySKU(c.Request.Context(), sku)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func ListProducts(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := db.ListProducts(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"products": products})
	}
}

func UpdateProduct(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only editor and admin can update products
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		if userRole != "editor" && userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		sku := c.Param("sku")
		var req models.UpdateProductRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product, err := db.UpdateProduct(c.Request.Context(), sku, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func DeleteProduct(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only admin can delete products
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		if userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "only admin can delete products"})
			return
		}

		sku := c.Param("sku")
		err := db.DeleteProduct(c.Request.Context(), sku)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
	}
}
