package handlers

import (
	"encoding/json"
	"fmt"
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
	log.Printf("Generated OAuth URL: %s", url)

	c.JSON(http.StatusOK, gin.H{
		"redirect_url": url,
	})
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	log.Printf("Received callback with code: %s", code)

	oauthToken, err := config.GoogleOAuthConfig.Exchange(c, code)
	if err != nil {
		log.Printf("Token exchange error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
		return
	}
	log.Printf("Successfully exchanged token")

	client := config.GoogleOAuthConfig.Client(c, oauthToken)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user info"})
		return
	}
	defer userInfo.Body.Close()

	userData, err := ioutil.ReadAll(userInfo.Body)
	if err != nil {
		log.Printf("Failed to read user info: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read user info"})
		return
	}
	log.Printf("Received user data: %s", string(userData))

	var googleUser struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(userData, &googleUser); err != nil {
		log.Printf("Failed to parse user info: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse user info"})
		return
	}

	// Insert or update user in database
	stmt, err := config.DB.Prepare(`
		INSERT INTO users (email, name) 
		VALUES (?, ?)
		ON CONFLICT(email) 
		DO UPDATE SET name = excluded.name
	`)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(googleUser.Email, googleUser.Name)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	log.Printf("Successfully saved user to database")

	// Generate JWT token
	jwtToken := config.GenerateJWTToken(googleUser.Email)

	// Return response with HTML that calls parent window
	data := gin.H{
		"email": googleUser.Email,
		"name":  googleUser.Name,
		"token": jwtToken,
	}

	jsonData, _ := json.Marshal(data)
	html := `
	<!DOCTYPE html>
	<html>
		<script>
			try {
				window.opener.postMessage({
					type: 'oauth-completion',
					data: %s
				}, 'http://localhost:3000');
				window.close();
			} catch (e) {
				console.error('Failed to communicate with main window:', e);
				document.body.innerHTML = 'Please return to the main application window.';
			}
		</script>
	</html>`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, fmt.Sprintf(html, string(jsonData)))
}

func GetCurrentUser(c *gin.Context) {
	// TODO: Implement proper session management
	// For now, just return error as unauthorized
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "Not authenticated",
	})
}
