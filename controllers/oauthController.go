package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// OAuthConfig represents the Google OAuth configuration
var OAuthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

// HandleOAuthLogin initiates the Google OAuth flow
func HandleOAuthLogin(c *gin.Context) {
	url := OAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// HandleOAuthCallback handles the callback from Google OAuth service
func HandleOAuthCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code parameter"})
		return
	}

	token, err := OAuthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client := OAuthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		// Add other fields as needed
	}
	if err := json.Unmarshal(data, &userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userInfo})

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{

		"email": userInfo.Email,

		"name": userInfo.Name,
	})

	// Generate encoded JWT token
	var jwtKey = []byte(os.Getenv("SECRET"))

	tokenString, err := jwtToken.SignedString(jwtKey)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})

		return

	}

	// Send JWT token to client

	c.JSON(http.StatusOK, gin.H{"token": tokenString})

}
