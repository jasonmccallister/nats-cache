package auth

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
