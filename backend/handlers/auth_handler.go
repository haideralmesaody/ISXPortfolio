package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	oauthConfig *oauth2.Config
	jwtSecret   []byte
}

func NewAuthHandler(oauthConfig *oauth2.Config) *AuthHandler {
	return &AuthHandler{
		oauthConfig: oauthConfig,
		jwtSecret:   []byte(os.Getenv("JWT_SECRET")),
	}
}

// GoogleLogin initiates the Google OAuth2 flow
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	url := h.oauthConfig.AuthCodeURL("state")
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// GoogleCallback handles the OAuth2 callback from Google
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := h.oauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := h.oauthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	userBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read user info"})
		return
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(userBytes, &userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	// Create JWT token
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["email"] = userInfo["email"]
	claims["name"] = userInfo["name"]
	claims["picture"] = userInfo["picture"]

	tokenString, err := jwtToken.SignedString(h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	log.Printf("User logged in: %v", userInfo["email"])
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"email":   userInfo["email"],
			"name":    userInfo["name"],
			"picture": userInfo["picture"],
		},
	})
}

// GetCurrentUser returns the current user's information
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	tokenString := authHeader[7:] // Remove "Bearer " prefix
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return h.jwtSecret, nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.JSON(http.StatusOK, gin.H{
			"email":   claims["email"],
			"name":    claims["name"],
			"picture": claims["picture"],
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
	}
}
