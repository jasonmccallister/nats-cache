package credentials

import (
	"fmt"
	"os"
)

const (
	templateFile = `-----BEGIN NATS USER JWT-----
%s
------END NATS USER JWT------

************************* IMPORTANT *************************
NKEY Seed printed below can be used to sign and prove identity.
NKEYs are sensitive and should be treated as secrets.

-----BEGIN USER NKEY SEED-----
%s
------END USER NKEY SEED------

*************************************************************`
)

// Generate generates a NATS credential file with the provided NKEY and JWT.
// If a file path is not provided a random file will be generated.
func Generate(nkey, jwt, file string) (string, error) {
	if nkey == "" {
		return "", fmt.Errorf("nkey is required")
	}

	if jwt == "" {
		return "", fmt.Errorf("jwt is required")
	}

	// file is optional, if not provided generate a random file name
	if file == "" {
		f, err := os.CreateTemp(os.TempDir(), "nats-creds-")
		if err != nil {
			return "", fmt.Errorf("failed to create temp file: %w", err)
		}
		defer f.Close()

		file = f.Name()
	}

	// write the nkey and jwt to the file
	if err := os.WriteFile(file, []byte(fmt.Sprintf(templateFile, jwt, nkey)), 0600); err != nil {
		return "", fmt.Errorf("failed to write credentials file: %w", err)
	}

	return file, nil
}
