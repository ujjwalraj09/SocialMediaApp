package controllers

import (
	"go/src/ujjwal/initializers"
	"go/src/ujjwal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FollowController struct {
	DB *gorm.DB
}

func NewFollowController(db *gorm.DB) *FollowController {
	return &FollowController{
		DB: initializers.DB,
	}
}

// FollowUser follows another user
func (fc *FollowController) FollowUser(c *gin.Context) {
	var follow models.Follow
	if err := c.BindJSON(&follow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the follow relationship in the database
	if err := fc.DB.Create(&follow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		return
	}

	c.JSON(http.StatusCreated, follow)
}

// UnfollowUser unfollows a user
func (fc *FollowController) UnfollowUser(c *gin.Context) {
	// Get follower ID and followed ID from request params
	followerID := c.Param("follower_id")
	followedID := c.Param("followed_id")

	// Delete the follow relationship from the database
	if err := fc.DB.Where("follower_id = ? AND followed_id = ?", followerID, followedID).Delete(&models.Follow{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unfollowed successfully"})
}

// GetFollowers returns the list of followers for a user
func (fc *FollowController) GetFollowers(c *gin.Context) {
	// Get user ID from request params
	userID := c.Param("user_id")

	// Fetch followers from database
	var followers []models.Follow
	if err := fc.DB.Where("followed_id = ?", userID).Find(&followers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch followers"})
		return
	}

	c.JSON(http.StatusOK, followers)
}

// GetFollowing returns the list of users being followed by a user
func (fc *FollowController) GetFollowing(c *gin.Context) {
	// Get user ID from request params
	userID := c.Param("user_id")

	// Fetch followed users from database
	var following []models.Follow
	if err := fc.DB.Where("follower_id = ?", userID).Find(&following).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch followed users"})
		return
	}

	c.JSON(http.StatusOK, following)
}
