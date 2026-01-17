package handlers

import (
	"net/http"

	"github.com/aslam/backend/internal/database"
	"github.com/aslam/backend/internal/middleware"
	"github.com/aslam/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func Register(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Only admin can create users
		userRole, exists := c.Get("user_role")
		if !exists || userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "only admin can create users"})
			return
		}

		user, err := db.CreateUser(c.Request.Context(), req.Email, req.Password, req.Role)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "user created successfully",
			"user": gin.H{
				"id":    user.ID,
				"email": user.Email,
				"role":  user.Role,
			},
		})
	}
}

func Login(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := db.GetUserByEmail(c.Request.Context(), req.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		if !db.ValidatePassword(user.Password, req.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		token, err := middleware.GenerateToken(user.ID, user.Email, string(user.Role))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, models.LoginResponse{
			Token: token,
			User: struct {
				ID    string          `json:"id"`
				Email string          `json:"email"`
				Role  models.UserRole `json:"role"`
			}{
				ID:    user.ID,
				Email: user.Email,
				Role:  user.Role,
			},
		})
	}
}

func GetProfile(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		user, err := db.GetUserByID(c.Request.Context(), userID.(string))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		})
	}
}

func UpdateProfile(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		var req models.UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := db.UpdateUser(c.Request.Context(), userID.(string), &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		})
	}
}

func ListUsers(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only admin can list users
		userRole, exists := c.Get("user_role")
		if !exists || userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "only admin can list users"})
			return
		}

		users, err := db.ListUsers(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := make([]gin.H, len(users))
		for i, user := range users {
			response[i] = gin.H{
				"id":    user.ID,
				"email": user.Email,
				"role":  user.Role,
			}
		}

		c.JSON(http.StatusOK, gin.H{"users": response})
	}
}

func DeleteUser(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only admin can delete users
		userRole, exists := c.Get("user_role")
		if !exists || userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "only admin can delete users"})
			return
		}

		userID := c.Param("id")
		err := db.DeleteUser(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
	}
}
