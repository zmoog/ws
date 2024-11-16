package identity

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Retriever interface {
	GetToken() (Token, error)
}

type tokenRetriever struct {
	client   *http.Client
	username string
	password string
}

func (r *tokenRetriever) GetToken() (Token, error) {
	// post a application/x-www-form-urlencoded request to the token endpoint
	// with the username and password
	form := url.Values{}
	form.Set("username", r.username)
	form.Set("password", r.password)
	form.Set("grant_type", "password")

	req, err := http.NewRequest("POST", "https://wavin-api.jablotron.cloud/v2.6/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return Token{}, err
	}

	req.Header.Add("Authorization", "Basic YXBwOnNlY3JldA==")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := r.client.Do(req)
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

	tokenResponse.ExpiresAt = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)

	return tokenResponse, nil
}
