package handlers

import (
	"net/http"

	"github.com/aslam/backend/internal/database"
	"github.com/aslam/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func CreateShelf(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only editor and admin can create shelves
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		if userRole != "editor" && userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		var req models.CreateShelfRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		shelf, err := db.CreateShelf(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, shelf)
	}
}

func GetShelf(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		shelf, err := db.GetShelfByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, shelf)
	}
}

func ListShelves(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		shelves, err := db.ListShelfs(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"shelves": shelves})
	}
}

func UpdateShelf(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only editor and admin can update shelves
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		if userRole != "editor" && userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		id := c.Param("id")
		var req models.UpdateShelfRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		shelf, err := db.UpdateShelf(c.Request.Context(), id, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, shelf)
	}
}

func DeleteShelf(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only admin can delete shelves
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		if userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "only admin can delete shelves"})
			return
		}

		id := c.Param("id")
		err := db.DeleteShelf(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "shelf deleted successfully"})
	}
}

func AddItemToShelf(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only editor and admin can add items
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		if userRole != "editor" && userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		shelfID := c.Param("id")
		var req models.AddItemToShelfRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		item, err := db.AddItemToShelf(c.Request.Context(), shelfID, req.SKU, req.Quantity)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, item)
	}
}

func RemoveItemFromShelf(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only editor and admin can remove items
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		if userRole != "editor" && userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		itemID := c.Param("itemId")
		err := db.RemoveItemFromShelf(c.Request.Context(), itemID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "item removed successfully"})
	}
}

func UpdateItemQuantity(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only editor and admin can update items
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		if userRole != "editor" && userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		itemID := c.Param("itemId")
		var req struct {
			Quantity int `json:"quantity" binding:"required,gt=0"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := db.UpdateItemQuantity(c.Request.Context(), itemID, req.Quantity)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "item quantity updated successfully"})
	}
}
