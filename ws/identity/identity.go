package identity

import "time"

// Token represents an access token.
type Token struct {
	ID           string    `json:"idToken"`
	DisplayName  string    `json:"displayName"`
	Kind         string    `json:"kind"`
	Email        string    `json:"email"`
	LocalID      string    `json:"localId"`
	Registered   bool      `json:"registered"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresIn    string    `json:"expiresIn"`
	ExpiresAt    time.Time `json:"expiresAt"`
}
