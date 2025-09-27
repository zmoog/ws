package identity

import (
	"errors"
	"testing"
	"time"
)

// Mock implementations for testing
type mockRetriever struct {
	token Token
	err   error
}

func (m *mockRetriever) GetToken() (Token, error) {
	return m.token, m.err
}

func (m *mockRetriever) RefreshToken(refreshToken string) (Token, error) {
	return m.token, m.err
}

type mockStorer struct {
	token      Token
	exists     bool
	getError   error
	storeError error
}

func (m *mockStorer) GetToken() (Token, bool, error) {
	return m.token, m.exists, m.getError
}

func (m *mockStorer) StoreToken(token Token) error {
	m.token = token
	m.exists = true
	return m.storeError
}

func TestManager_GetToken_ValidCachedToken(t *testing.T) {
	// Arrange
	validToken := Token{
		ID:        "valid-token",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	mockStore := &mockStorer{
		token:  validToken,
		exists: true,
	}
	mockRetriever := &mockRetriever{}

	manager := &manager{
		retriever: mockRetriever,
		storer:    mockStore,
	}

	// Act
	token, err := manager.GetToken()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if token.ID != validToken.ID {
		t.Errorf("Expected token ID %s, got %s", validToken.ID, token.ID)
	}
}

func TestManager_GetToken_ExpiredCachedToken(t *testing.T) {
	// Arrange
	expiredToken := Token{
		ID:        "expired-token",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	newToken := Token{
		ID:        "new-token",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	mockStore := &mockStorer{
		token:  expiredToken,
		exists: true,
	}
	mockRetriever := &mockRetriever{
		token: newToken,
	}

	manager := &manager{
		retriever: mockRetriever,
		storer:    mockStore,
	}

	// Act
	token, err := manager.GetToken()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if token.ID != newToken.ID {
		t.Errorf("Expected token ID %s, got %s", newToken.ID, token.ID)
	}
	// Verify new token was stored
	if mockStore.token.ID != newToken.ID {
		t.Errorf("Expected stored token ID %s, got %s", newToken.ID, mockStore.token.ID)
	}
}

func TestManager_GetToken_NoTokenExists(t *testing.T) {
	// Arrange
	newToken := Token{
		ID:        "new-token",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	mockStore := &mockStorer{
		exists: false,
	}
	mockRetriever := &mockRetriever{
		token: newToken,
	}

	manager := &manager{
		retriever: mockRetriever,
		storer:    mockStore,
	}

	// Act
	token, err := manager.GetToken()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if token.ID != newToken.ID {
		t.Errorf("Expected token ID %s, got %s", newToken.ID, token.ID)
	}
	// Verify new token was stored
	if mockStore.token.ID != newToken.ID {
		t.Errorf("Expected stored token ID %s, got %s", newToken.ID, mockStore.token.ID)
	}
}

func TestManager_GetToken_StoreGetError(t *testing.T) {
	// Arrange
	mockStore := &mockStorer{
		getError: errors.New("store error"),
	}
	mockRetriever := &mockRetriever{}

	manager := &manager{
		retriever: mockRetriever,
		storer:    mockStore,
	}

	// Act
	_, err := manager.GetToken()

	// Assert
	if err == nil {
		t.Error("Expected error from store, got nil")
	}
	if err.Error() != "store error" {
		t.Errorf("Expected 'store error', got %v", err)
	}
}

func TestManager_GetToken_RetrieverError(t *testing.T) {
	// Arrange
	mockStore := &mockStorer{
		exists: false,
	}
	mockRetriever := &mockRetriever{
		err: errors.New("retriever error"),
	}

	manager := &manager{
		retriever: mockRetriever,
		storer:    mockStore,
	}

	// Act
	_, err := manager.GetToken()

	// Assert
	if err == nil {
		t.Error("Expected error from retriever, got nil")
	}
	if err.Error() != "retriever error" {
		t.Errorf("Expected 'retriever error', got %v", err)
	}
}

func TestManager_GetToken_StoreTokenError(t *testing.T) {
	// Arrange
	newToken := Token{
		ID:        "new-token",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	mockStore := &mockStorer{
		exists:     false,
		storeError: errors.New("store token error"),
	}
	mockRetriever := &mockRetriever{
		token: newToken,
	}

	manager := &manager{
		retriever: mockRetriever,
		storer:    mockStore,
	}

	// Act
	_, err := manager.GetToken()

	// Assert
	if err == nil {
		t.Error("Expected error from store token, got nil")
	}
	if err.Error() != "store token error" {
		t.Errorf("Expected 'store token error', got %v", err)
	}
}

func TestNewManager(t *testing.T) {
	// Act
	manager := NewManager("user", "pass", "key")

	// Assert
	if manager == nil {
		t.Error("Expected manager to be created, got nil")
	}

	// Verify manager implements Manager interface
	var _ Manager = manager //nolint
}
