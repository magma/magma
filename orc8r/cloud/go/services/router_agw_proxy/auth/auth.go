package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
	"io"
	"fmt"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

var (
	currentToken string
	expiryTime   time.Time
	mu           sync.Mutex
)

func GetToken(tokenURL, username, password, clientID, clientSecret string) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	if currentToken != "" && time.Now().Before(expiryTime) {
		return currentToken, nil
	}

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", username)
	data.Set("password", password)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		log.Fatal("Error fetching token:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Token request failed, status=%d, body=%s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		log.Fatal("Error decoding token response:", err)
	}

	fmt.Println("Got access token:", tokenResp.AccessToken)

	currentToken = tokenResp.AccessToken
	expiryTime = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second) // refresh 1 min before expiry

	log.Println("Fetched new access token")
	return currentToken, nil
}
