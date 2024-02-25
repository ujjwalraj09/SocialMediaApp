package controllers

import (
	"go/src/ujjwal/initializers"
	"go/src/ujjwal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentController struct {
	DB *gorm.DB
}

func NewCommentController() *CommentController {
	// Ensure database connection is established
	initializers.ConnectToDb()

	// Assign DB instance to controller
	return &CommentController{
		DB: initializers.DB,
	}
}

func (cc *CommentController) CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cc.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GetComment retrieves a specific comment by ID
func (cc *CommentController) GetComment(c *gin.Context) {
	var comment models.Comment
	commentID := c.Param("id")

	// Fetch comment from database based on comment_id
	if err := cc.DB.Where("comment_id = ?", commentID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// gets all the comments in a post

func (cc *CommentController) GetCommentOnPost(c *gin.Context) {
	var comments []models.Comment
	postID := c.Param("id")

	// Find comments by post_id
	if err := cc.DB.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comments not found"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// UpdateComment updates an existing comment
func (cc *CommentController) UpdateComment(c *gin.Context) {
	var comment models.Comment
	commentID := c.Param("id")

	// Fetch the comment from the database based on comment_id
	if err := cc.DB.Where("comment_id = ?", commentID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Bind the JSON payload to the comment struct
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the updated comment to the database
	if err := cc.DB.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// DeleteComment deletes a comment by its ID
func (cc *CommentController) DeleteComment(c *gin.Context) {
	var comment models.Comment
	commentID := c.Param("id")

	// Find the comment by its comment_id
	if err := cc.DB.Where("comment_id = ?", commentID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Delete the comment from the database
	if err := cc.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
