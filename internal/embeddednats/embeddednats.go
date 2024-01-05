package embeddednats

import (
	"fmt"
	"github.com/jasonmccallister/nats-cache/natsremote"
	"github.com/nats-io/nats-server/v2/server"
)

// NewServer creates a new embedded NATS server and will return the server and the credentials file.
func NewServer(port, httpPort int) (*server.Server, string, error) {
	remote, creds, err := natsremote.RemoteLeafFromEnv()
	if err != nil {
		return nil, "", fmt.Errorf("failed to create remote leaf: %w", err)
	}

	opts := &server.Options{
		JetStream: true,
		LeafNode: server.LeafNodeOpts{
			Remotes: []*server.RemoteLeafOpts{remote},
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
