package identity

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Storer interface {
	StoreToken(token Token) error
	GetToken() (Token, bool, error)
}

type tokenStorer struct{}

func (s *tokenStorer) StoreToken(token Token) error {
	settingsPath, err := s.getSettingsPath()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(settingsPath), 0700)
	if err != nil {
		return err
	}

	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	return os.WriteFile(settingsPath, data, 0600)
}

func (s *tokenStorer) GetToken() (Token, bool, error) {

	settingsPath, err := s.getSettingsPath()
	if err != nil {
		return Token{}, false, err
	}

	if _, err := os.Stat(settingsPath); errors.Is(err, os.ErrNotExist) {
		return Token{}, false, nil
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return Token{}, false, err
	}

	token := Token{}
	err = json.Unmarshal(data, &token)
	if err != nil {
		return Token{}, false, err
	}

	return token, true, nil
}

func (s *tokenStorer) getSettingsPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".ws", "identity"), nil
}
