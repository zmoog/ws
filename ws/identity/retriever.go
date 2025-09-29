package identity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

//nolint:lll
const (
	signInWithPasswordEndpoint = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key="
	tokenEndpoint              = "https://securetoken.googleapis.com/v1/token?key="
	tokenLifespanPercentage    = 90 // percentage of the token's lifespan to use the calculated expires at value.
)

// Retriever is an interface for retrieving a token.
type Retriever interface {
	GetToken() (Token, error)
	RefreshToken(refreshToken string) (Token, error)
}

// tokenRetriever is a concrete implementation of Retriever.
type tokenRetriever struct {
	httpClient *http.Client
	webApiKey  string
	username   string
	password   string
}

// GetToken retrieves a token from the token endpoint.
func (r *tokenRetriever) GetToken() (Token, error) {
	req := struct {
		Email             string `json:"email"`
		Password          string `json:"password"`
		ClientType        string `json:"clientType"`
		ReturnSecureToken bool   `json:"returnSecureToken"`
	}{
		Email:             r.username,
		Password:          r.password,
		ClientType:        "CLIENT_TYPE_WEB",
		ReturnSecureToken: true,
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return Token{}, err
	}

	request, err := http.NewRequest(
		"POST",
		signInWithPasswordEndpoint+r.webApiKey,
		bytes.NewReader(jsonReq),
	)
	if err != nil {
		return Token{}, err
	}

	request.Header.Add("Content-Type", "application/json")

	resp, err := r.httpClient.Do(request)
	if err != nil {
		return Token{}, err
	}
	defer resp.Body.Close() // nolint

	if resp.StatusCode != http.StatusOK {
		// return the response body as a string
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return Token{}, err
		}

		return Token{}, fmt.Errorf("unexpected status code: %d\n%s", resp.StatusCode, string(body))
	}

	var tokenResponse struct {
		IDToken      string `json:"idToken"`
		ExpiresIn    string `json:"expiresIn"`
		LocalID      string `json:"localId"`
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return Token{}, err
	}

	expiresIn, err := strconv.Atoi(tokenResponse.ExpiresIn)
	if err != nil {
		return Token{}, fmt.Errorf("failed to parse expiresIn: %v", err)
	}

	token := Token{
		ID:           tokenResponse.IDToken,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresIn:    expiresIn,
		ExpiresAt:    calculateExpiresAt(expiresIn),
		LocalID:      tokenResponse.LocalID,
	}

	return token, nil
}

// RefreshToken refreshes an expired token using a refresh token.
func (r *tokenRetriever) RefreshToken(refreshToken string) (Token, error) {
	req := struct {
		GrantType    string `json:"grant_type"`
		RefreshToken string `json:"refresh_token"`
	}{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return Token{}, err
	}

	request, err := http.NewRequest(
		"POST",
		tokenEndpoint+r.webApiKey,
		bytes.NewReader(jsonReq),
	)
	if err != nil {
		return Token{}, err
	}

	request.Header.Add("Content-Type", "application/json")

	resp, err := r.httpClient.Do(request)
	if err != nil {
		return Token{}, err
	}
	defer resp.Body.Close() // nolint

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return Token{}, err
		}
		return Token{}, fmt.Errorf("refresh token failed with status %d: %s", resp.StatusCode, string(body))
	}

	var refreshResponse struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    string `json:"expires_in"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
		IDToken      string `json:"id_token"`
		UserID       string `json:"user_id"`
		ProjectID    string `json:"project_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&refreshResponse); err != nil {
		return Token{}, err
	}

	expiresIn, err := strconv.Atoi(refreshResponse.ExpiresIn)
	if err != nil {
		return Token{}, fmt.Errorf("failed to parse expiresIn: %v", err)
	}

	token := Token{
		ID:           refreshResponse.IDToken,
		RefreshToken: refreshResponse.RefreshToken,
		ExpiresIn:    expiresIn,
		ExpiresAt:    calculateExpiresAt(expiresIn),
		LocalID:      refreshResponse.UserID,
	}

	return token, nil
}

// calculateExpiresAt calculates the expiration time for a token.
func calculateExpiresAt(expiresIn int) time.Time {
	return time.Now().
		// We calculate the expiration time once, so it's handy to have it
		// available later.
		//
		// We force a token usage duration of 90% of the token's
		// lifespan to avoid expiration while the request is in
		// flight.
		//
		// It happens if we request a token with a lifespan of 1 hour and
		// use the token every 60s: the token is valid when we check it,
		// but it will be expired when we use it.
		//
		// We refresh the token a little earlier, but it will be
		// always valid when we use it.
		Add(time.Duration(expiresIn*tokenLifespanPercentage/100) * time.Second)
}
