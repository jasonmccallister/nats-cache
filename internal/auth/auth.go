package auth

import (
	"fmt"
	"os"
)

type Token struct {
	Audience string `json:"aud,omitempty"`
	Expires  int64  `json:"exp,omitempty"`
	IssuedAt int64  `json:"iat,omitempty"`
	Subject  string `json:"sub,omitempty"`
}

// Authorizer is an interface for authorizing a token
type Authorizer interface {
	Authorize(token string) (*Token, error)
}

// NewFromEnvironment creates an Authorizer from the environment variables
func NewFromEnvironment() (Authorizer, error) {
	authType := "paseto"
	if v, ok := os.LookupEnv("AUTH_TYPE"); ok {
		authType = v
	}

	switch authType {
	case "paseto":
		v, ok := os.LookupEnv("AUTH_PUBLIC_KEY")
		if !ok {
			return nil, fmt.Errorf("AUTH_PUBLIC_KEY not set")
		}

		return NewPasetoPublicKey(v), nil
	}

	return nil, fmt.Errorf("unknown AUTH_TYPE: %s", authType)
}
