package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const tokenFile = "token.json"

// GetTokenSource возвращает источник токенов OAuth2.
func GetTokenSource(ctx context.Context, serviceAccountFile string, scope string) (oauth2.TokenSource, error) {
	jsonKey, err := os.ReadFile(serviceAccountFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read service account file: %w", err)
	}

	config, err := google.ConfigFromJSON(jsonKey, scope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service account JSON: %w", err)
	}

	token, err := readTokenFromFile(tokenFile)
	if err != nil {
		log.Println("No existing token found, requesting new token...")
		token = getNewToken(ctx, config)
		saveTokenToFile(tokenFile, token)
	}

	return config.TokenSource(ctx, token), nil
}

// getNewToken выполняет ручную авторизацию пользователя для получения токена.
func getNewToken(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Visit the URL for the auth dialog: %v", authURL)

	var authCode string
	log.Print("Enter the authorization code: ")
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	token, err := config.Exchange(ctx, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	return token
}

// readTokenFromFile читает токен из файла.
func readTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

// saveTokenToFile сохраняет токен в файл.
func saveTokenToFile(file string, token *oauth2.Token) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache OAuth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
