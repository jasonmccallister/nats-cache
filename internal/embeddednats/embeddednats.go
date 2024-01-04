package embeddednats

import (
	"fmt"
	"github.com/jasonmccallister/nats-cache/credentials"
	"github.com/nats-io/nats-server/v2/server"
	"net/url"
)

// NewServer creates a new embedded NATS server and will return the server and the credentials file.
func NewServer(nkey, jwt string, port, httpPort int) (*server.Server, string, error) {
	creds, err := credentials.Generate(nkey, jwt, "")
	if err != nil {
		return nil, "", fmt.Errorf("unable to generate credentials: %w", err)
	}

	opts := &server.Options{
		JetStream: true,
		LeafNode: server.LeafNodeOpts{
			Remotes: []*server.RemoteLeafOpts{
				{
					URLs: []*url.URL{
						{
							Scheme: "tls",
							Host:   "connect.ngs.global",
						},
					},
					Credentials: creds,
				},
			},
		},
	}

	if port > 0 {
		opts.Port = port
	}

	if httpPort > 0 {
		opts.HTTPPort = httpPort
	}

	s, err := server.NewServer(opts)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create server: %w", err)
	}

	return s, creds, nil
}
