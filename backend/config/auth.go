package config

import (
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	GoogleOAuthConfig *oauth2.Config
)

func InitGoogleOAuth(clientID, clientSecret string) {
	if clientID == "" {
		log.Fatal("GOOGLE_CLIENT_ID environment variable is not set")
	}
	if clientSecret == "" {
		log.Fatal("GOOGLE_CLIENT_SECRET environment variable is not set")
	}

	log.Printf("Initializing OAuth with Client ID: %s", clientID)

	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8000/auth/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
