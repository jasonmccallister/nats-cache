package keygen

import (
	"fmt"

	"github.com/jasonmccallister/nats-cache/internal/auth"
)

// FromToken creates an internal key from a token using the subj, database and key provided
// and returns the internal key, the original key and an error if one occurred
func FromToken(t auth.Token, db uint32, key string) (string, string, error) {
	return fmt.Sprintf("%s.%d-%s", t.Subject, db, key), key, nil
}
