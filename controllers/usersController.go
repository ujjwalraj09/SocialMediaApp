package controllers

import (
	"go/src/ujjwal/initializers"
	"go/src/ujjwal/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	// Get user input from the request body
	var body struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the email is already taken
	var existingUser models.User
	if err := initializers.DB.Where("email = ?", body.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already taken"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create a new user record
	newUser := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hashedPassword),
	}
	if err := initializers.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func Login(c *gin.Context) {
	//Get the email and pass of req body
	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	//lookup requested user

	var user models.User
	initializers.DB.First(&user, "email= ?", body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email or password",
		})
		return
	}
	//compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email or password",
		})
		return
	}
	//generate jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	//sign and get the complete encoded token as a string using the secret

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create token",
		})
		return
	}

	//send it
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}
func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged Out",
	})
}

func Validate(c *gin.Context) {
	//user, _ := c.Get("user")
	//user.(models.User)
	c.JSON(http.StatusOK, gin.H{
		"message": "i'm logged in ",
	})
}

// implementing home function
func Home(c *gin.Context) {
	// Check if user is authenticated
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// User is authenticated, serve content
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the home page",
		"user":    user,
	})
}
