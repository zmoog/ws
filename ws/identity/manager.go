package identity

import (
	"fmt"
	"net/http"
	"os"
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

func (m *manager) GetToken() (Token, error) {
	token, exists, err := m.store.GetToken()
	if err != nil {
		return Token{}, err
	}

	if exists && token.ExpiresAt.After(time.Now()) {
		fmt.Fprintf(os.Stderr, "Using cached token (expires at %s)\n", token.ExpiresAt.Format(time.RFC3339))
		return token, nil
	}

	token, err = m.retriever.GetToken()
	if err != nil {
		return Token{}, err
	}

	fmt.Fprintf(os.Stderr, "Generated new token (expires at %s)\n", token.ExpiresAt.Format(time.RFC3339))

	err = m.store.StoreToken(token)
	if err != nil {
		return Token{}, err
	}

	return token, nil
}
