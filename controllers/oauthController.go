package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Replace these with your actual configuration details
var oauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("CLIENTID"),
	ClientSecret: os.Getenv("CLIENTKEY"),
	RedirectURL:  "http://localhost:3000",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

// jwtKey represents the secret key for JWT token signing
var jwtKey = []byte(os.Getenv("SECRET"))

// HandleGoogleAuth handles Google OAuth authentication
func HandleGoogleAuth(c *gin.Context) {
	code := c.Query("code")
	token, err := oauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Use the token to create a new client
	client := oauthConfig.Client(c, token)

	// Now you can use this client to make requests to protected resources
	userInfoEndpoint := "https://www.googleapis.com/oauth2/v3/userinfo"
	resp, err := client.Get(userInfoEndpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Decode the response body into your struct
	userInfo := struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create JWT token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": userInfo.Email,
		"name":  userInfo.Name,
	})

	// Generate encoded JWT token
	tokenString, err := jwtToken.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}

	// Send JWT token to client
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
