package natsconn

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"net/url"
	"os"
)

// FromEnv will create a nats connection from the environment variables. Environment variables are
func FromEnv() (*nats.Conn, error) {
	u, err := url.Parse("connect.ngs.global")
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	// is the url set?
	if v, ok := os.LookupEnv("NATS_URL"); ok {
		u, err = url.Parse(v)
		if err != nil {
			return nil, fmt.Errorf("failed to parse url: %w", err)
		}
	}

	var opts []nats.Option

	// check for the name
	if v, ok := os.LookupEnv("NATS_NAME"); ok {
		opts = append(opts, nats.Name(v))
	}

	// check for the echo
	if _, ok := os.LookupEnv("NATS_NO_ECHO"); ok {
		opts = append(opts, nats.NoEcho())
	}

	// if we have a nkey seed, use it only if we have a jwt
	if s, ok := os.LookupEnv("NATS_SEED"); ok {
		// check for jwt as well
		if j, ok := os.LookupEnv("NATS_JWT"); ok {
			opts = append(opts, nats.UserJWTAndSeed(j, s))
		}
	}

	// is we have the nats user/pass, use it only if we have a pass
	if u, ok := os.LookupEnv("NATS_USER"); ok {
		if p, ok := os.LookupEnv("NATS_PASS"); ok {
			opts = append(opts, nats.UserInfo(u, p))
		}
	}

	// is we have the nats token, use it
	if t, ok := os.LookupEnv("NATS_AUTH"); ok {
		opts = append(opts, nats.Token(t))
	}

	return nats.Connect(u.String(), opts...)
}
