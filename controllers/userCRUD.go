package controllers

import (
	"go/src/ujjwal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// basic user functinalities
type UserController struct {
	DB *gorm.DB
}

// UpdateUser updates an existing user profile
func (uc *UserController) UpdateUser(c *gin.Context) {
	// Parse form data or JSON payload
	userID := c.Param("id")
	// Fetch user from database
	var user models.User
	if err := uc.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update user fields
	user.Name = c.PostForm("name")
	user.Email = c.PostForm("email")
	// Update other fields...

	// Save updated user record
	if err := uc.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
func (uc *UserController) DeleteUser(c *gin.Context) {
	// Get user ID from the URL parameter
	userID := c.Param("id")

	// Fetch user from database
	var user models.User
	if err := uc.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete the user from the database
	if err := uc.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// UploadProfileImage uploads a profile image for a user

func (uc *UserController) UploadProfileImage(c *gin.Context) {
	// Get user ID from the URL parameter
	userID := c.Param("id")

	// Fetch user from database
	var user models.User
	if err := uc.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Parse JSON payload to get profile image data
	var requestData struct {
		ProfileImage string `json:"profile_image"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// Update user's profile image
	user.ProfileImage = requestData.ProfileImage

	// Save updated user record with profile image
	if err := uc.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile image uploaded successfully"})
}
