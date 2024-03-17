package controllers

import (
	"go/src/ujjwal/initializers"
	"go/src/ujjwal/models"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FeedController struct {
	DB *gorm.DB
}

// NewFeedController initializes a new FeedController with the provided database connection
func NewFeedController() *FeedController {
	// Ensure database connection is established
	initializers.ConnectToDb()

	// Assign DB instance to controller
	return &FeedController{
		DB: initializers.DB,
	}
}

func (fc *FeedController) CreatePost(c *gin.Context) {
	// Parse the request body to get the post content and tags
	var postInput struct {
		PostContent string   `json:"post_content"`
		PostTags    []string `json:"post_tags"`
	}
	if err := c.BindJSON(&postInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle file uploads for post images
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	files := form.File["post_images"]

	// Create the new post
	post := models.Post{
		PostContent: postInput.PostContent,
		PostTags:    []models.Tags{},
	}
	for _, tag := range postInput.PostTags {
		post.PostTags = append(post.PostTags, models.Tags{Tag: tag})
	}

	// Save the post to the database
	if err := fc.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Handle the uploaded post images
	for _, file := range files {
		if file != nil {
			imageData, err := file.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer imageData.Close()

			// Read the image data and store it in the PostImage model
			imgBytes, err := io.ReadAll(imageData)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			postImage := models.PostImage{
				PostID: post.ID,
				Image:  imgBytes,
			}

			if err := fc.DB.Create(&postImage).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusOK, post)
}

func (fc *FeedController) GetPost(c *gin.Context) {
	var post models.Post
	postID := c.Param("id")

	if err := fc.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (fc *FeedController) UpdatePost(c *gin.Context) {
	var post models.Post
	postID := c.Param("id")

	if err := fc.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := fc.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (fc *FeedController) DeletePost(c *gin.Context) {
	var post models.Post
	postID := c.Param("id")

	if err := fc.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if err := fc.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// GetLike retrieves a specific like by ID
func (fc *FeedController) GetLike(c *gin.Context) {
	var like models.Like
	likeID := c.Param("id")

	if err := fc.DB.First(&like, likeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Like not found"})
		return
	}

	c.JSON(http.StatusOK, like)
}

// DeleteLike deletes a like by its ID
func (fc *FeedController) DeleteLike(c *gin.Context) {
	var like models.Like
	likeID := c.Param("id")

	if err := fc.DB.First(&like, likeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Like not found"})
		return
	}
	var post models.Post

	if err := fc.DB.Delete(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete like"})
		return
	}
	post.LikeCount--
	c.JSON(http.StatusOK, gin.H{"message": "Like deleted successfully"})
}

func (fc *FeedController) CreateLike(c *gin.Context) {
	// Bind the like request to a Like struct
	var like models.Like
	if err := c.BindJSON(&like); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Increment the LikeCount for the post
	var post models.Post
	if err := fc.DB.First(&post, "id = ?", like.PostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	post.LikeCount++
	if err := fc.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like count"})
		return
	}

	// Create the like in the database
	if err := fc.DB.Create(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create like"})
		return
	}

	c.JSON(http.StatusCreated, like)
}

func (fc *FeedController) GetLikeCount(c *gin.Context) {
	postID := c.Param("id")

	var post models.Post
	if err := fc.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post_id": postID, "like_count": post.LikeCount})
}

// Tags Searching functionality
func (fc *FeedController) FindImagesByTags(c *gin.Context) {
	// Parse the tags from form data
	tags := c.PostFormArray("tags")

	// Query the database to find posts containing any of the specified tags
	var posts []models.Post
	if err := fc.DB.Where("tags IN ?", tags).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find posts by tags"})
		return
	}

	// Collect all images from the posts
	var images []string
	for _, post := range posts {
		for _, image := range post.PostImages {
			images = append(images, string(image.Image))
		}
	}

	// Return the collected images in the response
	c.JSON(http.StatusOK, images)
}
