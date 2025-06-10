package identity

import (
	"net/http"
	"time"
)

type Manager interface {
	Retriever
}

type manager struct {
	retriever Retriever
	store     Storer
}

func NewManager(username, password string) Manager {
	retriever := tokenRetriever{
		client:   http.DefaultClient,
		username: username,
		password: password,
	}
	store := tokenStorer{}

	return &manager{
		retriever: &retriever,
		store:     &store,
	}
}

// GetToken returns a token from the store if it exists and is not expired,
// otherwise it retrieves a new token.
func (m *manager) GetToken() (Token, error) {
	token, exists, err := m.store.GetToken()
	if err != nil {
		return Token{}, err
	}

	if exists && token.ExpiresAt.After(time.Now()) {
		return token, nil
	}

	token, err = m.retriever.GetToken()
	if err != nil {
		return Token{}, err
	}

	err = m.store.StoreToken(token)
	if err != nil {
		return Token{}, err
	}

	return token, nil
}
