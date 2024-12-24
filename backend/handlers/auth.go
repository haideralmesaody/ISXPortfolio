package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"isxportfolio-backend/config"

	"github.com/gin-gonic/gin"
)

func GoogleLogin(c *gin.Context) {
	if config.GoogleOAuthConfig == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OAuth config not initialized"})
		return
	}

	url := config.GoogleOAuthConfig.AuthCodeURL("state")
	log.Printf("Redirecting to OAuth URL: %s", url)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := config.GoogleOAuthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := config.GoogleOAuthConfig.Client(c, token)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user info"})
		return
	}
	defer userInfo.Body.Close()

	userData, err := ioutil.ReadAll(userInfo.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read user info"})
		return
	}

	var googleUser struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(userData, &googleUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse user info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email": googleUser.Email,
		"name":  googleUser.Name,
	})
}
