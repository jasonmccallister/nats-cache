package natsremote

import (
	"fmt"
	"github.com/jasonmccallister/nats-cache/credentials"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"net/url"
	"os"
)

// RemoteLeafFromEnv will create a remote leaf from the environment variables.
func RemoteLeafFromEnv() (*server.RemoteLeafOpts, string, error) {
	u, err := url.Parse("connect.ngs.global")
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse url: %w", err)
	}

	// we need the jwt and seed to be set

	// if nats seed is set, and jwt is set, use it
	var creds string
	if s, ok := os.LookupEnv("NATS_SEED"); ok {
		if j, ok := os.LookupEnv("NATS_JWT"); ok {
			creds, err = credentials.Generate(s, j, "")
			if err != nil {
				return nil, "", fmt.Errorf("failed to generate credentials: %w", err)
			}
		} else {
			return nil, "", fmt.Errorf("NATS_JWT is not set")
		}
	} else {
		return nil, "", fmt.Errorf("NATS_SEED is not set")
	}

	return &server.RemoteLeafOpts{
		URLs:        []*url.URL{u},
		Credentials: creds,
	}, creds, nil
}

// FromEnv will create a nats connection from the environment variables. Currently not used.
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
