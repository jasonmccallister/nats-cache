package auth

import (
	"fmt"

	"aidanwoods.dev/go-paseto"
)

type pasetoPublicKey struct {
	PublicKey string
}

// NewPasetoPublicKey creates a new Paseto public key authorizer
func NewPasetoPublicKey(publicKey string) Authorizer {
	return &pasetoPublicKey{
		PublicKey: publicKey,
	}
}

func (p *pasetoPublicKey) Authorize(token string) (*Token, error) {
	publicKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(p.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create public key: %w", err)
	}

	parser := paseto.NewParserWithoutExpiryCheck()
	t, err := parser.ParseV4Public(publicKey, token, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	aud, err := t.GetAudience()
	if err != nil {
		return nil, fmt.Errorf("failed to get audience: %w", err)
	}

	exp, err := t.GetExpiration()
	if err != nil {
		return nil, fmt.Errorf("failed to get expiration: %w", err)
	}

	iat, err := t.GetIssuedAt()
	if err != nil {
		return nil, fmt.Errorf("failed to get issued at: %w", err)
	}

	sub, err := t.GetSubject()
	if err != nil {
		return nil, fmt.Errorf("failed to get subject: %w", err)
	}

	return &Token{
		Audience: aud,
		Expires:  exp.Unix(),
		IssuedAt: iat.Unix(),
		Subject:  sub,
	}, nil
}
