package auth

import (
	"fmt"
	"strings"

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
	fields := strings.Fields(token)
	if len(fields) != 2 {
		return nil, fmt.Errorf("invalid token, must be in the format of Bearer <token>")
	}

	publicKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(p.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create public key: %w", err)
	}

	parser := paseto.NewParserWithoutExpiryCheck()
	t, err := parser.ParseV4Public(publicKey, fields[1], nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// aud, err := t.GetAudience()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get audience: %w", err)
	// }

	// exp, err := t.GetExpiration()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get expiration: %w", err)
	// }

	// iat, err := t.GetIssuedAt()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get issued at: %w", err)
	// }

	sub, err := t.GetSubject()
	if err != nil {
		return nil, fmt.Errorf("failed to get subject: %w", err)
	}

	return &Token{
		//Audience: aud,
		//Expires:  exp.Unix(),
		//IssuedAt: iat.Unix(),
		Subject: sub,
	}, nil
}
