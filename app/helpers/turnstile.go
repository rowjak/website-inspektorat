package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type TurnstileResponse struct {
	Success     bool     `json:"success"`
	ErrorCodes  []string `json:"error-codes"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
}

func VerifyTurnstile(token string, remoteIP string) (bool, error) {
	// secretKey := os.Getenv("TURNSTILE_SECRET_KEY")
	// if secretKey == "" {
	// 	return false, fmt.Errorf("TURNSTILE_SECRET_KEY not configured")
	// }

	// Prepare form data
	data := url.Values{}
	data.Set("secret", "0x4AAAAAAAfyGQqSmIQWaHYzeNy8zPUD6jE")
	data.Set("response", token)
	data.Set("remoteip", remoteIP)

	// Make request to Turnstile API
	resp, err := http.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify", data)
	if err != nil {
		return false, fmt.Errorf("failed to verify turnstile: %v", err)
	}
	defer resp.Body.Close()

	// Parse response
	var turnstileResp TurnstileResponse
	if err := json.NewDecoder(resp.Body).Decode(&turnstileResp); err != nil {
		return false, fmt.Errorf("failed to decode turnstile response: %v", err)
	}

	return turnstileResp.Success, nil
}

func GetTurnstileSiteKey() string {
	return "0x4AAAAAAAfyGSdfP3zSl_EN"
	// return os.Getenv("TURNSTILE_SITE_KEY")
}
