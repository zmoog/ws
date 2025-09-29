package identity

import (
	"net/http"
)

type Manager interface {
	GetToken() (Token, error)
}

type manager struct {
	retriever Retriever
	storer    Storer
}

func NewManager(config Config) Manager {
	retriever := tokenRetriever{
		httpClient: http.DefaultClient,
		username:   config.Username,
		password:   config.Password,
		webApiKey:  config.WebApiKey,
	}
	storer := tokenStorer{}

	return &manager{
		retriever: &retriever,
		storer:    &storer,
	}
}

// NewInMemoryManager creates a manager that stores the token in memory.
func NewInMemoryManager(config Config) Manager {
	retriever := tokenRetriever{
		httpClient: http.DefaultClient,
		username:   config.Username,
		password:   config.Password,
		webApiKey:  config.WebApiKey,
	}
	storer := NewInMemoryStorer()

	return &manager{
		retriever: &retriever,
		storer:    storer,
	}
}

// GetToken returns a token from the store if it exists and is not expired,
// otherwise it retrieves a new token.
func (m *manager) GetToken() (Token, error) {
	token, exists, err := m.storer.GetToken()
	if err != nil {
		return Token{}, err
	}

	if exists && !token.IsExpired() {
		// Token is still valid
		return token, nil
	}

	// Token is expired or does not exist
	// Try to refresh first if we have a refresh token
	if exists && token.RefreshToken != "" {
		refreshedToken, refreshErr := m.retriever.RefreshToken(token.RefreshToken)
		if refreshErr == nil {
			// Successfully refreshed token
			err = m.storer.StoreToken(refreshedToken)
			if err != nil {
				return Token{}, err
			}
			return refreshedToken, nil
		}
		// Refresh failed, fall back to full authentication
	}

	// Get new token with credentials
	token, err = m.retriever.GetToken()
	if err != nil {
		return Token{}, err
	}

	// Store the new token
	err = m.storer.StoreToken(token)
	if err != nil {
		return Token{}, err
	}

	return token, nil
}
