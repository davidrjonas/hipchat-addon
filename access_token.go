package addon

import (
	"encoding/json"
	"io"
	"time"
)

const graceSeconds = 120

type AccessToken struct {
	Value     string `json:"access_token"`
	TokenType string `json:"token_type"`
	Scope     string `json:"scope"`
	ExpiresIn int32  `json:"expires_in"`
	ExpiresAt int64  `json:"-"`
}

func (t *AccessToken) IsExpired() bool {
	return t.ExpiresAt < time.Now().Unix()
}

func (t *AccessToken) Valid() bool {
	return t.Value != "" && !t.IsExpired()
}

func (t *AccessToken) String() string {
	return t.Value
}

func NewAccessTokenFromJson(r io.Reader) (*AccessToken, error) {

	token := &AccessToken{}

	if err := json.NewDecoder(r).Decode(token); err != nil {
		return nil, err
	}

	token.ExpiresAt = nextUnixExpiry(token.ExpiresIn)

	return token, nil
}

func nextUnixExpiry(secondsFromNow int32) int64 {
	return time.Now().Unix() + int64(secondsFromNow) - graceSeconds
}
