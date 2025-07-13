package identity

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Retriever is an interface for retrieving a token.
type Retriever interface {
	GetToken() (Token, error)
}

// tokenRetriever is a concrete implementation of Retriever.
type tokenRetriever struct {
	client    *http.Client
	webApiKey string
	username  string
	password  string
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
		"https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key="+r.webApiKey,
		strings.NewReader(string(jsonReq)),
	)
	if err != nil {
		return Token{}, err
	}

	request.Header.Add("Content-Type", "application/json")

	resp, err := r.client.Do(request)
	if err != nil {
		return Token{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// return the response body as a string
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return Token{}, err
		}

		return Token{}, fmt.Errorf("unexpected status code: %d\n%s", resp.StatusCode, string(body))
	}

	var tokenResponse Token
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return Token{}, err
	}

	expiresIn, err := strconv.Atoi(tokenResponse.ExpiresIn)
	if err != nil {
		return Token{}, fmt.Errorf("failed to parse expiresIn: %v", err)
	}

	// We calculate the expiration time once, so it's handy
	// to have it available later.
	tokenResponse.ExpiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)

	return tokenResponse, nil
}
