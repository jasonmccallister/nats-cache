package auth

type Token struct {
	Audience string `json:"aud"`
	Expires  int64  `json:"exp"`
	IssuedAt int64  `json:"iat"`
	Subject  string `json:"sub"`
}

// Authorizer is an interface for authorizing a token
type Authorizer interface {
	Authorize(token string) (*Token, error)
}
