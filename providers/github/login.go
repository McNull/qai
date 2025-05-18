package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"github.com/mcnull/qai/shared/throbber"
	"github.com/mcnull/qai/shared/utils"
	"strings"
	"time"
)

const DEVICE_CODE_URL = "https://github.com/login/device/code"
const COPILOT_API_KEY = "Iv1.b507a08c87ecfe98"
const OAUTH_TOKEN_URL = "https://github.com/login/oauth/access_token"
const API_TOKEN_URL = "https://api.github.com/copilot_internal/v2/token"

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type OAuthTokenResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	Scope            string `json:"scope"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorURI         string `json:"error_uri,omitempty"`
}

type ApiTokenResponse struct {
	ExpiresAt int    `json:"expires_at"`
	RefreshIn int    `json:"refresh_in"`
	Token     string `json:"token"`
}

func Login(debug bool) (string, error) {

	// 1. Request device and user codes
	dc, err := getDeviceCode()
	if err != nil {
		return "", fmt.Errorf("failed to request device code: %w", err)
	}

	if debug {
		utils.Dump(dc)
	}

	// 2. Prompt user to authorize
	fmt.Printf("Please visit: %s\n", dc.VerificationURI)
	fmt.Printf("And enter code: %s\n", dc.UserCode)

	// 3. Poll for the token
	interval := time.Duration(dc.Interval) * time.Second
	accessToken, err := pollForToken(dc.DeviceCode, interval)
	if err != nil {
		return "", fmt.Errorf("failed to poll for token: %w", err)
	}

	fmt.Println("Successfully authenticated!")

	return accessToken, nil
}

func requestOAuthToken(deviceCode string) (*OAuthTokenResponse, error) {

	clientID := COPILOT_API_KEY
	tokenURL := OAUTH_TOKEN_URL

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("device_code", deviceCode)
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send token request: %w", err)
	}
	defer resp.Body.Close()

	// GitHub returns 200 OK even for errors in the device flow, error is in JSON body
	var tokenResponse OAuthTokenResponse
	bodyBytes, err := readAll(resp.Body) // Read body once
	if err != nil {
		return nil, fmt.Errorf("failed to read token response body: %w", err)
	}

	if err := json.Unmarshal(bodyBytes, &tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w (body: %s)", err, string(bodyBytes))
	}

	if resp.StatusCode != http.StatusOK && tokenResponse.Error == "" {
		// If status is not OK and there's no error in JSON, it's an unexpected situation.
		return nil, fmt.Errorf("failed to request token: status %s, body: %s", resp.Status, string(bodyBytes))
	}

	return &tokenResponse, nil
}

func pollForToken(deviceCode string, interval time.Duration) (string, error) {
	throbber := throbber.NewThrobber().
		WithMessage("Waiting for authorization ... ").
		WithThrob(throbber.ThrobByName("dots"))

	throbber.Start()
	defer throbber.Stop()

	for {
		tokenResponse, err := requestOAuthToken(deviceCode)
		if err != nil {
			// This would be an error in making the request itself or decoding, not a GitHub error code.
			return "", fmt.Errorf("error requesting token: %w", err)
		}

		if tokenResponse.Error != "" {
			switch tokenResponse.Error {
			case "authorization_pending":
				// User has not yet entered the code. Wait, then poll again.
				time.Sleep(interval)
				continue
			case "slow_down":
				// App polled too fast. Wait for the interval plus 5 seconds, then poll again.
				time.Sleep(interval + 5*time.Second)
				continue
			case "expired_token":
				return "", fmt.Errorf("device code has expired. Please try the login process again. Description: %s", tokenResponse.ErrorDescription)
			case "access_denied":
				return "", fmt.Errorf("login cancelled by user or access denied. Description: %s", tokenResponse.ErrorDescription)
			default:
				return "", fmt.Errorf("received error from token endpoint: %s. Description: %s", tokenResponse.Error, tokenResponse.ErrorDescription)
			}
		}

		if tokenResponse.AccessToken != "" {
			return tokenResponse.AccessToken, nil
		}

		// If no access token and no error, something is unexpected. Wait and retry once more after a delay.
		// Or, this could be a point to return an error if this state is not expected.
		time.Sleep(interval)
	}
}

func getDeviceCode() (*DeviceCodeResponse, error) {

	// POST https://github.com/login/device/code
	// Accept: application/json
	// Content-Type: application/x-www-form-urlencoded
	//
	// client_id=Iv1.b507a08c87ecfe98

	data := url.Values{}
	data.Set("client_id", COPILOT_API_KEY)

	req, err := http.NewRequest("POST", DEVICE_CODE_URL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var deviceCode DeviceCodeResponse
	err = json.NewDecoder(resp.Body).Decode(&deviceCode)
	if err != nil {
		return nil, err
	}

	return &deviceCode, nil

}

func readAll(r io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func requestApiToken(oauth_token string) (string, error) {
	// GET https://api.github.com/copilot_internal/v2/token
	// Authorization: Bearer {{$dotenv GITHUB_AUTH_TOKEN}}
	// User-Agent: github.com/mcnull/qai
	// Accept: application/json

	req, err := http.NewRequest("GET", API_TOKEN_URL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+oauth_token)
	req.Header.Set("User-Agent", "github.com/mcnull/qai")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var t ApiTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
		return "", err
	}
	if t.Token == "" {
		return "", fmt.Errorf("empty token in response")
	}

	return t.Token, nil
}
