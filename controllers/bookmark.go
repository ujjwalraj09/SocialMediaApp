package controllers

import (
	"go/src/ujjwal/initializers"
	"go/src/ujjwal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookMarkController struct {
	DB *gorm.DB
}


func NewBookMarkController() *BookMarkController {
	// Ensure database connection is established
	initializers.ConnectToDb()

	// Assign DB instance to controller
	return &BookMarkController{
		DB: initializers.DB,
	}
}

func (bkmk *BookMarkController) CreateBookmark(c *gin.Context) {
	var bookmark models.BookMark

	if err := c.BindJSON(&bookmark); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if len(bookmark.UserPosts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bookmark required"})
		return
	}

	// Create the post in the database
	if err := bkmk.DB.Create(&bookmark).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bookmark"})
		return

	}

	c.JSON(http.StatusCreated, bookmark)

}
func (bkmk *BookMarkController) GetBookMark(c *gin.Context) {
	// Get user ID from the URL parameter
	userID := c.Param("id")

	// Fetch bookmark from database
	var bookmark models.BookMark
	if err := bkmk.DB.Preload("UserPosts").First(&bookmark, "user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bookmark not found"})
		return
	}

	c.JSON(http.StatusOK, bookmark)
}

func (bkmk *BookMarkController) UpdateBookMark(c *gin.Context) {
	// Get user ID from the URL parameter
	userID := c.Param("id")

	// Fetch bookmark from database
	var bookmark models.BookMark
	if err := bkmk.DB.Preload("UserPosts").First(&bookmark, "user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bookmark not found"})
		return
	}

	// Update bookmark fields
	if err := c.BindJSON(&bookmark); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save updated bookmark record
	if err := bkmk.DB.Save(&bookmark).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bookmark"})
		return
	}

	c.JSON(http.StatusOK, bookmark)
}

func (bkmk *BookMarkController) DeleteBookMark(c *gin.Context) {
	// Get user ID from the URL parameter
	userID := c.Param("id")

	// Fetch bookmark from database
	var bookmark models.BookMark
	if err := bkmk.DB.First(&bookmark, "user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bookmark not found"})
		return
	}

	// Delete the bookmark from the database
	if err := bkmk.DB.Delete(&bookmark).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete bookmark"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bookmark deleted successfully"})
}
