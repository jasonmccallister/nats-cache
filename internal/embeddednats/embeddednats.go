package embeddednats

import (
	"fmt"
	"github.com/jasonmccallister/nats-cache/natsremote"
	"github.com/nats-io/nats-server/v2/server"
	"os"
	"strconv"
)

// NewServer creates a new embedded NATS server and will return the server and the credentials file.
func NewServer(port, httpPort int, remotes []*server.RemoteLeafOpts, creds string) (*server.Server, string, error) {
	opts := &server.Options{
		JetStream: true,
		LeafNode: server.LeafNodeOpts{
			Remotes: remotes,
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

// NewServerFromEnvironment will create a new embedded NATS server from the environment variables.
func NewServerFromEnvironment() (*server.Server, string, error) {
	portV, _ := os.LookupEnv("NATS_PORT")
	port, _ := strconv.Atoi(portV)
	httpPortV, _ := os.LookupEnv("NATS_HTTP_PORT")
	httpPort, _ := strconv.Atoi(httpPortV)

	remote, creds, err := natsremote.RemoteLeafFromEnv()
	if err != nil {
		return nil, "", fmt.Errorf("failed to create remote leaf: %w", err)
	}

	return NewServer(port, httpPort, []*server.RemoteLeafOpts{remote}, creds)
}
